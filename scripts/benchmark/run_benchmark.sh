#!/usr/bin/env bash
# Measures the throughput of the agent for every combination of modules (logging, rule engine,
# ML classifier) and for different sizes of the classifier pool.
#
# The script builds the agent, the API and a minimal target web server, starts temporary
# Cassandra and Elasticsearch containers, runs the measurements with ffuf, writes a summary
# and removes everything it started.
#
# Usage: scripts/benchmark/run_benchmark.sh
# The options are environment variables, see the README (section "Throughput benchmark").
set -u

BENCH_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
export REPO_DIR=$(cd "$BENCH_DIR/../.." && pwd)
export WORK_DIR=${WORK_DIR:-$BENCH_DIR/work}
RESULTS_DIR=${RESULTS_DIR:-$BENCH_DIR/results/$(date +%Y%m%d-%H%M%S)}

REPETITIONS=${REPETITIONS:-3}
REQUESTS=${REQUESTS:-12000}
WARMUP_REQUESTS=${WARMUP_REQUESTS:-1000}
export THREADS=${THREADS:-40}
POOL_SIZES=${POOL_SIZES:-"1 2 4 8 16"}
EXPERIMENTS=${EXPERIMENTS:-"modules pool"}
SOURCE_WORDLIST=${SOURCE_WORDLIST:-/usr/share/seclists/Discovery/Web-Content/common.txt}
KEEP_CONTAINERS=${KEEP_CONTAINERS:-0}

export AGENT_PORT=${AGENT_PORT:-8000}
export BACKEND_PORT=${BACKEND_PORT:-8080}
export API_PORT=${API_PORT:-8081}
export ELASTIC_URL=http://127.0.0.1:9200

CASSANDRA_CONTAINER=bench-cassandra
ELASTIC_CONTAINER=bench-elastic

log() { echo "[$(date +%H:%M:%S)] $*"; }
die() {
	log "ERROR: $*"
	exit 1
}

# ---------------------------------------------------------------- dependencies
for tool in go ffuf docker curl python3 ss; do
	command -v "$tool" >/dev/null || die "$tool is not installed"
done
[ -f "$SOURCE_WORDLIST" ] || die "wordlist $SOURCE_WORDLIST not found (install seclists or set SOURCE_WORDLIST)"

# Use sudo for docker if the user cannot access the docker socket
DOCKER="docker"
docker ps >/dev/null 2>&1 || DOCKER="sudo docker"
$DOCKER ps >/dev/null 2>&1 || die "cannot run docker"

for port in "$AGENT_PORT" "$BACKEND_PORT" "$API_PORT" 9042 9200; do
	ss -ltn | grep -q ":$port " && die "port $port is already in use (stop the service using it before running the benchmark)"
done

mkdir -p "$WORK_DIR/logs" "$RESULTS_DIR"

# ---------------------------------------------------------------- cleanup
PIDS=()
cleanup() {
	log "Cleaning up"
	for pid in "${PIDS[@]}"; do
		kill "$pid" 2>/dev/null
	done
	pkill -f "$WORK_DIR/agent-bin" 2>/dev/null
	if [ "$KEEP_CONTAINERS" != 1 ]; then
		$DOCKER rm -f -v "$CASSANDRA_CONTAINER" "$ELASTIC_CONTAINER" >/dev/null 2>&1
	fi
}
trap cleanup EXIT
trap 'exit 130' INT TERM

# ---------------------------------------------------------------- python environment
export PYTHON_BIN="$WORK_DIR/venv/bin/python"
if [ ! -x "$PYTHON_BIN" ]; then
	log "Creating the python virtual environment for the classifier"
	python3 -m venv "$WORK_DIR/venv" || die "could not create the virtual environment"
	# The models were serialized with scikit-learn 1.6.1
	"$WORK_DIR/venv/bin/pip" install -q scikit-learn==1.6.1 pandas numpy || die "could not install the python packages"
fi

# ---------------------------------------------------------------- build
log "Building the agent, the API and the target web server"
(cd "$REPO_DIR/agent" && go build -o "$WORK_DIR/agent-bin" .) || die "could not build the agent"
(cd "$REPO_DIR/api" && go build -o "$WORK_DIR/api-bin" .) || die "could not build the API"
(cd "$BENCH_DIR/backend" && go build -o "$WORK_DIR/backend-bin" .) || die "could not build the target web server"

# ---------------------------------------------------------------- wordlists
# Lines containing spaces are removed because the HTTP server rejects them before they reach the agent.
# The line number is appended so that every request is unique.
export WORDLIST="$WORK_DIR/wordlist.txt"
export WARMUP_WORDLIST="$WORK_DIR/warmup.txt"
grep -v '[[:space:]]' "$SOURCE_WORDLIST" | grep -v '^#' | grep -v '^$' >"$WORK_DIR/words.txt"
python3 - "$WORK_DIR/words.txt" "$REQUESTS" "$WARMUP_REQUESTS" "$WORDLIST" "$WARMUP_WORDLIST" <<'PY'
import sys
words = open(sys.argv[1]).read().split("\n")
words = [w for w in words if w]
requests, warmup = int(sys.argv[2]), int(sys.argv[3])
with open(sys.argv[4], "w") as f:
    for i in range(requests):
        f.write(f"{words[i % len(words)]}&n={i}\n")
with open(sys.argv[5], "w") as f:
    for i in range(warmup):
        f.write(f"{words[i % len(words)]}&w={i}\n")
PY

# ---------------------------------------------------------------- databases
log "Starting Cassandra and Elasticsearch"
$DOCKER rm -f -v "$CASSANDRA_CONTAINER" "$ELASTIC_CONTAINER" >/dev/null 2>&1
# The host network avoids the docker proxy, which drops connections under the load of the benchmark
$DOCKER run -d --name "$CASSANDRA_CONTAINER" --network host \
	-e CASSANDRA_CLUSTER_NAME=honeypot -e MAX_HEAP_SIZE=2G -e HEAP_NEWSIZE=400M \
	cassandra:latest >/dev/null || die "could not start cassandra"
# The disk watermark is disabled, otherwise elasticsearch rejects the logs when the disk is almost full
$DOCKER run -d --name "$ELASTIC_CONTAINER" --network host \
	-e discovery.type=single-node -e xpack.security.enabled=false \
	-e cluster.routing.allocation.disk.threshold_enabled=false -e ES_JAVA_OPTS="-Xms2g -Xmx2g" \
	docker.elastic.co/elasticsearch/elasticsearch:8.7.1 >/dev/null || die "could not start elasticsearch"

timeout 300 bash -c "until curl -s $ELASTIC_URL >/dev/null; do sleep 3; done" || die "elasticsearch did not start"
timeout 300 bash -c "until $DOCKER exec $CASSANDRA_CONTAINER cqlsh -e 'describe keyspaces' >/dev/null 2>&1; do sleep 5; done" ||
	die "cassandra did not start"
$DOCKER exec "$CASSANDRA_CONTAINER" cqlsh -e \
	"CREATE KEYSPACE IF NOT EXISTS api WITH replication = {'class': 'SimpleStrategy', 'replication_factor': '1'};" ||
	die "could not create the cassandra keyspace"

# ---------------------------------------------------------------- API and target web server
cat >"$WORK_DIR/api.conf" <<JSON
{
    "address": "127.0.0.1",
    "port": "$API_PORT",
    "cassandraNodes": ["127.0.0.1"],
    "cassandraKeyspace": "api",
    "elasticUrl": "$ELASTIC_URL",
    "elasticIndex": "logs",
    "exploitTemplatePath": "templates/exploit.tmpl",
    "adminUsername": "admin",
    "adminPassword": "admin"
}
JSON
log "Starting the API and the target web server"
(cd "$REPO_DIR/api" && exec "$WORK_DIR/api-bin" -config "$WORK_DIR/api.conf" >"$WORK_DIR/logs/api.log" 2>&1) &
PIDS+=($!)
"$WORK_DIR/backend-bin" -address "127.0.0.1:$BACKEND_PORT" >"$WORK_DIR/logs/backend.log" 2>&1 &
PIDS+=($!)
timeout 120 bash -c "until curl -s 127.0.0.1:$API_PORT/api/v1/healthcheck >/dev/null; do sleep 1; done" || die "the API did not start"
timeout 30 bash -c "until curl -s 127.0.0.1:$BACKEND_PORT/ >/dev/null; do sleep 1; done" || die "the target web server did not start"

# ---------------------------------------------------------------- measurements
# name logging rules ml
CONFIGURATIONS=(
	"proxy 0 0 0"
	"logging 1 0 0"
	"rules 0 1 0"
	"logging-rules 1 1 0"
	"ml 0 0 1"
	"rules-ml 0 1 1"
	"all 1 1 1"
)

run() {
	local output=$1
	shift
	"$BENCH_DIR/bench_run.sh" "$@" | tee -a "$output"
}

MODULES_RESULTS="$RESULTS_DIR/modules.txt"
POOL_RESULTS="$RESULTS_DIR/pool.txt"
: >"$MODULES_RESULTS"
: >"$POOL_RESULTS"

for experiment in $EXPERIMENTS; do
	case $experiment in
	modules)
		log "Experiment 1: throughput by module configuration ($REPETITIONS repetitions, default pool size)"
		for rep in $(seq 1 "$REPETITIONS"); do
			for configuration in "${CONFIGURATIONS[@]}"; do
				set -- $configuration
				CLASSIFIER_PROCESSES=0 run "$MODULES_RESULTS" "$1-r$rep" "$2" "$3" "$4"
			done
		done
		;;
	pool)
		log "Experiment 2: throughput by classifier pool size ($POOL_SIZES)"
		for rep in $(seq 1 "$REPETITIONS"); do
			for size in $POOL_SIZES; do
				CLASSIFIER_PROCESSES=$size run "$POOL_RESULTS" "ml-p$size-r$rep" 0 0 1
				CLASSIFIER_PROCESSES=$size run "$POOL_RESULTS" "all-p$size-r$rep" 1 1 1
			done
		done
		;;
	*) die "unknown experiment $experiment (use modules and/or pool)" ;;
	esac
done

# ---------------------------------------------------------------- summary
{
	echo "# Agent throughput benchmark"
	echo
	echo "- Date: $(date '+%Y-%m-%d %H:%M')"
	echo "- CPU: $(lscpu | sed -n 's/^Model name: *//p') ($(nproc) logical CPUs)"
	echo "- Memory: $(free -g | awk '/^Mem:/ {print $2}') GB"
	echo "- OS: $(. /etc/os-release && echo "$PRETTY_NAME"), kernel $(uname -r)"
	echo "- Commit: $(git -C "$REPO_DIR" rev-parse --short HEAD)$(git -C "$REPO_DIR" diff --quiet || echo ' (with uncommitted changes)')"
	echo "- Requests per run: $REQUESTS (warm-up $WARMUP_REQUESTS), $THREADS threads, $REPETITIONS repetitions"
	echo
	SUMMARY_ARGS=("$MODULES_RESULTS")
	[ -s "$POOL_RESULTS" ] && SUMMARY_ARGS+=("$POOL_RESULTS")
	python3 "$BENCH_DIR/summarize.py" "${SUMMARY_ARGS[@]}"
} >"$RESULTS_DIR/summary.md"
SUMMARY_STATUS=$?

cat "$RESULTS_DIR/summary.md"
log "Results written to $RESULTS_DIR"
[ "$SUMMARY_STATUS" = 0 ] || die "some runs were invalid, see the summary"
