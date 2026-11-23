#!/usr/bin/env bash
# Runs one throughput measurement of the agent and checks that it is valid.
#
# Usage: bench_run.sh <name> <logging 0|1> <rules 0|1> <ml 0|1>
#
# Environment (set by run_benchmark.sh):
#   REPO_DIR, WORK_DIR, PYTHON_BIN    - repository, work directory and python interpreter of the venv
#   AGENT_PORT, BACKEND_PORT, API_PORT - ports of the agent, the target web server and the API
#   ELASTIC_URL                        - the elasticsearch URL (used to check that the logs were stored)
#   WORDLIST, WARMUP_WORDLIST          - the measured requests and the warm-up requests
#   THREADS                            - number of concurrent ffuf threads
#   CLASSIFIER_PROCESSES               - size of the classifier pool (0 = one per CPU)
#
# Prints one result line:
#   RESULT <name> logging=<0|1> rules=<0|1> ml=<0|1> processes=<n> requests=<n> time=<s> rps=<n> status=OK|INVALID (<reason>)
set -u

NAME=$1
LOGGING=$2
RULES=$3
ML=$4

CONFIG="$WORK_DIR/conf-$NAME.json"
AGENT_LOG="$WORK_DIR/logs/agent-$NAME.log"

# Without logging, point the agent at a closed port so it does not connect to the API
AGENT_API_PORT=$API_PORT
[ "$LOGGING" = 0 ] && AGENT_API_PORT=1
RULES_DIRECTORY=""
[ "$RULES" = 1 ] && RULES_DIRECTORY="./rules"
USE_AI=false
[ "$ML" = 1 ] && USE_AI=true

fail() {
	echo "RESULT $NAME logging=$LOGGING rules=$RULES ml=$ML processes=$CLASSIFIER_PROCESSES requests=0 time=0 rps=0 status=INVALID ($1)"
	[ -n "${AGENT_PID:-}" ] && kill "$AGENT_PID" 2>/dev/null
	exit 1
}

# Returns the number of documents in the logs index
es_count() {
	curl -s -XPOST "$ELASTIC_URL/logs/_refresh" >/dev/null
	curl -s "$ELASTIC_URL/logs/_count" | python3 -c "import json,sys; print(json.load(sys.stdin).get('count', 0))" 2>/dev/null || echo 0
}

# The port must be free, otherwise the requests would reach another process
timeout 30 bash -c "while ss -ltn | grep -q '127.0.0.1:$AGENT_PORT '; do sleep 1; done" || fail "port $AGENT_PORT is busy"

cat >"$CONFIG" <<JSON
{
    "protocol": "http",
    "address": "127.0.0.1",
    "port": "$AGENT_PORT",
    "forbiddenPagePath": "./static/forbidden.html",
    "blacklistUserAgentPath": "./lists/user_agent.list",
    "forwardServerProtocol": "http",
    "forwardServerAddress": "127.0.0.1",
    "forwardServerPort": "$BACKEND_PORT",
    "apiProtocol": "http",
    "apiIpAddress": "127.0.0.1",
    "apiPort": "$AGENT_API_PORT",
    "uuid": "",
    "rulesDirectory": "$RULES_DIRECTORY",
    "operationMode": "waf",
    "ignoreRulesDirectories": ["CVEs"],
    "useAIClassifier": $USE_AI,
    "classifier": "svc",
    "classifierProcesses": $CLASSIFIER_PROCESSES,
    "pythonInterpreter": "$PYTHON_BIN",
    "llmAPIURL": "",
    "createDataset": false,
    "datasetPath": "./datasets/benchmark.csv",
    "adaptiveTemplates": []
}
JSON

# The agent loads the rules, models and scripts relative to its directory
cd "$REPO_DIR/agent" || fail "agent directory not found"
"$WORK_DIR/agent-bin" -config "$CONFIG" >"$AGENT_LOG" 2>&1 &
AGENT_PID=$!

timeout 120 bash -c "until curl -s -o /dev/null 127.0.0.1:$AGENT_PORT/healthprobe; do sleep 0.5; done" || fail "agent did not start"
sleep 0.5
if grep -q 'address already in use' "$AGENT_LOG" || ! kill -0 "$AGENT_PID" 2>/dev/null; then
	fail "agent did not bind port $AGENT_PORT"
fi

# Warm-up
ffuf -s -u "http://127.0.0.1:$AGENT_PORT/search?q=FUZZ" -w "$WARMUP_WORDLIST" -mc all -t "$THREADS" -sa >/dev/null 2>&1

LOG_START=$(wc -l <"$AGENT_LOG")
ES_BEFORE=0
[ "$LOGGING" = 1 ] && ES_BEFORE=$(es_count)

START=$(date +%s.%N)
ffuf -s -u "http://127.0.0.1:$AGENT_PORT/search?q=FUZZ" -w "$WORDLIST" -mc all -t "$THREADS" >/dev/null 2>&1
END=$(date +%s.%N)

kill -INT "$AGENT_PID"
wait "$AGENT_PID" 2>/dev/null
AGENT_PID=""

# Checks: the agent received every request, no log or classifier errors, every log stored in elasticsearch
EXPECTED=$(wc -l <"$WORDLIST")
RECEIVED=$(tail -n +"$LOG_START" "$AGENT_LOG" | grep -c "Received GET request")
ERRORS=$(tail -n +"$LOG_START" "$AGENT_LOG" | grep -c "error on the server\|could not send the logs\|restarting it\|Failed to get prediction\|Classifier stderr")
STORED=0
[ "$LOGGING" = 1 ] && sleep 1 && STORED=$(($(es_count) - ES_BEFORE))

STATUS=OK
if [ "$RECEIVED" -lt "$EXPECTED" ]; then
	STATUS="INVALID (agent received $RECEIVED of $EXPECTED requests)"
elif [ "$ERRORS" -gt 0 ]; then
	STATUS="INVALID ($ERRORS log or classifier errors, see $AGENT_LOG)"
elif [ "$LOGGING" = 1 ] && [ "$STORED" -lt "$RECEIVED" ]; then
	STATUS="INVALID (only $STORED of $RECEIVED logs stored in elasticsearch)"
fi

# The throughput is the number of requests sent divided by the time ffuf needed to receive all the responses
python3 - "$EXPECTED" "$START" "$END" <<PY
import sys
requests, start, end = int(sys.argv[1]), float(sys.argv[2]), float(sys.argv[3])
print(f"RESULT $NAME logging=$LOGGING rules=$RULES ml=$ML processes=$CLASSIFIER_PROCESSES "
      f"requests={requests} time={end - start:.2f} rps={requests / (end - start):.1f} status=$STATUS")
PY
