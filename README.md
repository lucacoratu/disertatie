# disertatie
Master's degree final project - Adaptive Honeypot System

# Architecture
The honeypot system consists of 4 main components: an agent which will be deployed on the target machine, an API which will handle all the data processing and database interaction, a dashboard which will display the data in a human friendly manner and a Apache Cassandra cluster which will store all the data provided by the agents. 

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

Checking if the dashboard has all dependencies install
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

# Starting IDE
You can start the Visual Studio Code IDE using one of the scripts `start_ide.ps1` or `start_ide.sh` depending on the operating system you have installed on your machine (.ps1 for Windows, .sh for Linux) 