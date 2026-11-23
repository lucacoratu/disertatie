# Adaptive Honeypot System
Master's degree final project - Adaptive Honeypot System

# Architecture
The honeypot system consists of 4 main components: 
- an agent which will be deployed on the target machine
- an API which will handle all the data processing and database interaction
- a dashboard which will display the data in a human friendly manner 
- an Apache Cassandra cluster which will store all the data provided by the agents. 
- an ElasticSearch cluster which will store the logs from the agents for easier indexing.

# Dependencies
The honeypot system is built using Golang and NextJS

Installing NextJS:
```bash
#Install nodejs
curl -fsSL https://deb.nodesource.com/setup_21.x | sudo -E bash - && sudo apt-get install -y nodejs

#Check the node 21 was installed
node -v

#Check that npm was installed
npm -v
```

Installing modules for dashboard
```bash
cd ./dashboard/
npm install
```

Checking if the dashboard has all dependencies installed
```bash
#Start the dashboard in development mode
npm run dev
```

# Configure database
In order for the API to run, a Apache Cassandra cluster is needed. The easiest way is to start a single node cluster is to run docker-compose in the root of the repository.
```bash
docker-compose up -d
```
The docker-compose command will create the Cassandra container and the necessary keyspace for the API. It will also create 2 volumes, one for initializing the Cassandra keyspace and one which will assure persistence of the Cassandra database.

If the `init-cassandra` service fails to start, you should check the permissions of the `scripts/cassandra/init.sh` script (it should be executable - `chmod +x scripts/cassandra/init.sh`)

# Configuring Agent and API
Both the agent and the API have a configuration file inside their directory (`agent.conf` and `api.conf`). The format of the data in the configuration file is JSON.

## 1. API configuration
The API needs the following fields in the configuration file in order to start:
- address - This is the IP address the API will listen on
- port - The port to listen on
- cassandraNodes - The nodes of the Apache Cassandra cluster (list of strings)
- cassandraKeyspace - The keyspace of the Cassandra cluster
- exploitTemplatePath - The path to the exploit generation template.

Example configuration:
```json
{
    "address": "0.0.0.0",
    "port":"8081",
    "cassandraNodes": ["127.0.0.1"],
    "cassandraKeyspace": "api",
    "exploitTemplatePath": "templates/exploit.tmpl"
}
```

## 2. Agent configuration
The Agent needs the following fields in the configuration file in order to run:
- protocol - The protocol to listen on (http or https)
- address - The IP address the agent will listen on
- port - The port the agent will listen on
- forbiddenPagePath - Path to the forbidden page
- blacklistUserAgentPath - Path to the list of blacklisted user agents
- forwardServerProtocol - The protocol of the forward server (protected server)
- forwardServerAddress - The address of the forward server
- forwardServerPort - The port of the forward server
- apiProtocol - The protocol of the API
- apiIpAddress - The IP address of the API
- apiPort - The port of the API
- uuid - This will be completed at first run (pulled from API) - leave empty
- rulesDirectory - The rules directory
- operationMode - The mode of operation of the agent
- classifierProcesses - (optional) The number of python classifier processes started when `useAIClassifier` is true. `0` or missing means one process for each logical CPU. Each process uses about 130 MB of memory
- pythonInterpreter - (optional) The python interpreter used to run the classifier (default `/usr/bin/python`), for example the interpreter of a virtual environment

Example Agent configuration:
```json
{
    "protocol": "http",
    "address": "0.0.0.0",
    "port": "8080",
    "forbiddenPagePath": "./static/forbidden.html",
    "blacklistUserAgentPath": "./lists/user_agent.list",
    "forwardServerProtocol": "http",
    "forwardServerAddress": "127.0.0.1",
    "forwardServerPort": "8000",
    "apiProtocol": "http",
    "apiIpAddress": "127.0.0.1",
    "apiPort": "8081",
    "uuid": "237870e8-f2f5-49dd-b63c-8d67c4003513",
    "rulesDirectory": "./rules",
    "operationMode": "testing"
}
```


# Docker images
The agent, the API and the dashboard have a Dockerfile in their directory:
```bash
docker build -t honeypot-agent agent        # the configuration is mounted in /config/agent.conf
docker build -t honeypot-api api            # the configuration is mounted in /config/api.conf
docker build -t honeypot-dashboard \
    --build-arg NEXT_PUBLIC_API_URL=http://<api>:8081/api/v1 \
    --build-arg NEXT_PUBLIC_API_WS_URL=ws://<public address of the api>:8081/api/v1/ws \
    dashboard
```
`NEXT_PUBLIC_API_URL` is used by the server side of the dashboard and `NEXT_PUBLIC_API_WS_URL` by the browser, so the second one must be reachable by the users of the dashboard. Both default to `127.0.0.1:8081`.

# Evaluation deployment
The `evaluation` directory contains the Terraform configuration which creates the virtual machines on OpenStack (one for the API, its databases and the dashboard, and two agents protecting DVWA and OWASP Juice Shop) and the Ansible playbooks which deploy every component in Docker. See [evaluation/README.md](evaluation/README.md).

# Throughput benchmark
The `scripts/benchmark` directory contains the scripts used to measure the throughput of the agent (the number of requests per second it can handle) with different modules enabled.

| File | Description |
|---|---|
| `run_benchmark.sh` | Main script. Prepares the environment, runs all the measurements, writes a summary and cleans up |
| `bench_run.sh` | Runs and validates a single measurement (called by `run_benchmark.sh`) |
| `summarize.py` | Computes the mean and the standard deviation of every configuration and writes Markdown tables |
| `backend/` | Minimal Go web server used as the protected application. It returns the same static page for every request, so the measurement reflects the cost of the agent |

## Requirements
- Go, Python 3, `curl`, `ss` and Docker (the script uses `sudo docker` if the current user cannot access the Docker socket)
- [ffuf](https://github.com/ffuf/ffuf), used to send the requests
- The SecLists `common.txt` wordlist (`/usr/share/seclists/Discovery/Web-Content/common.txt` on Kali, or set `SOURCE_WORDLIST`)
- Free ports: 8000 (agent), 8080 (target web server), 8081 (API), 9042 (Cassandra) and 9200 (Elasticsearch). The script stops if any of them is in use, so stop the containers from `docker-compose.yml` first
- About 6 GB of free memory (Cassandra, Elasticsearch and up to one classifier process per CPU)

## Running
```bash
./scripts/benchmark/run_benchmark.sh
```
With the default options the benchmark takes about 15 minutes. The script:
1. creates a Python virtual environment in `scripts/benchmark/work/venv` with the version of scikit-learn used to serialize the models (only the first time)
2. builds the agent, the API and the target web server
3. starts temporary Cassandra and Elasticsearch containers (`bench-cassandra` and `bench-elastic`) on the host network
4. starts the API and the target web server
5. runs the two experiments described below, starting a new agent for every measurement
6. writes the results to `scripts/benchmark/results/<date>/` and removes the containers and the processes it started

The agent runs in `waf` mode with the rules from `agent/rules` (without the `CVEs` directory) and the SVC classifier. Every measurement sends a warm-up of 1,000 requests, then 12,000 requests of the form `/search?q=<word>&n=<i>` from 40 concurrent ffuf threads. The throughput is the number of requests divided by the time ffuf needs to receive all the responses.

## Experiments
- `modules` - all the combinations of logging, rule engine and ML classifier in Table 8 of the paper, with the default size of the classifier pool:

| Name | Logging | Rule engine | ML classifier |
|---|---|---|---|
| `proxy` | | | |
| `logging` | ✓ | | |
| `rules` | | ✓ | |
| `logging-rules` | ✓ | ✓ | |
| `ml` | | | ✓ |
| `rules-ml` | | ✓ | ✓ |
| `all` | ✓ | ✓ | ✓ |

- `pool` - the ML classifier alone and all the modules together, for every size of the classifier pool in `POOL_SIZES`

"Logging" means that the agent sends every request to the API, which stores it in Cassandra and Elasticsearch. Without logging, the agent is configured with a closed API port, so it does not connect to the API.

## Validation
A measurement is marked `INVALID` and excluded from the tables if:
- the agent did not start or could not bind its port (for example because another process uses it)
- the agent did not log every request it should have received
- the agent log contains errors from sending logs to the API or from the classifier
- with logging enabled, Elasticsearch did not store one document for every request

The script exits with an error if any measurement was invalid.

## Options
The options are set as environment variables:

| Variable | Default | Description |
|---|---|---|
| `EXPERIMENTS` | `modules pool` | The experiments to run |
| `REPETITIONS` | `3` | The number of times every configuration is measured |
| `REQUESTS` | `12000` | The number of requests in every measurement |
| `WARMUP_REQUESTS` | `1000` | The number of warm-up requests sent before every measurement |
| `THREADS` | `40` | The number of concurrent ffuf threads |
| `POOL_SIZES` | `1 2 4 8 16` | The sizes of the classifier pool used in the `pool` experiment |
| `SOURCE_WORDLIST` | SecLists `common.txt` | The wordlist from which the requests are generated (lines containing spaces are removed) |
| `AGENT_PORT`, `BACKEND_PORT`, `API_PORT` | `8000`, `8080`, `8081` | The ports of the agent, the target web server and the API |
| `WORK_DIR` | `scripts/benchmark/work` | The directory with the binaries, the virtual environment, the generated configurations and the logs |
| `RESULTS_DIR` | `scripts/benchmark/results/<date>` | The directory where the results are written |
| `KEEP_CONTAINERS` | `0` | Set to `1` to keep the Cassandra and Elasticsearch containers after the benchmark |

For example, a quick run of only the first experiment:
```bash
EXPERIMENTS=modules REPETITIONS=1 ./scripts/benchmark/run_benchmark.sh
```
Very short runs (a few thousand requests) finish in less than a second for the fastest configurations, so their results vary more.

## Results
The results directory contains:
- `modules.txt` and `pool.txt` - one `RESULT` line for every measurement
- `summary.md` - the hardware, the commit and the tables with the mean and the standard deviation of every configuration

The logs of the agent (one file for every measurement), of the API and of the target web server are in `scripts/benchmark/work/logs`.

# Starting IDE
You can start the Visual Studio Code IDE using one of the scripts `start_ide.ps1` or `start_ide.sh` depending on the operating system you have installed on your machine (.ps1 for Windows, .sh for Linux)

# How to setup elasticsearch with volume
In order to setup the volume for the elasticsearch cluster you need to change the permission of the folder to be owned by the user with uid 1000.
```bash
sudo chown 1000:1000 ./data/elastic
```

# How to deploy the agent as a systemd service
To add the agent as a systemd service you have to create a file in /etc/systemd/system folder (ex: /etc/systemd/system/agent.service).

Add the following options in the file

```bash
[Unit]
Description=Disertatie Agent
ConditionPathExists=<PATH_TO_AGENT_CODE>
After=network.target
[Service]
Type=simple
User=<USER>
Group=<GROUP>
WorkingDirectory=<PATH_TO_AGENT_CODE>
ExecStart=/usr/local/go/bin/go run main.go -config ./agent.conf
Restart=always
RestartSec=10
StandardOutput=syslog
StandardError=syslog
SyslogIdentifier=agent
[Install]
WantedBy=multi-user.target
```