# Evaluation setup

This directory deploys the honeypot system on OpenStack for the real-world evaluation:

| Virtual machine | Containers |
|---|---|
| `central` | Cassandra, Elasticsearch, the API and the dashboard |
| `agent-dvwa` | An agent protecting [DVWA](https://github.com/digininja/DVWA) |
| `agent-juiceshop` | An agent protecting [OWASP Juice Shop](https://github.com/juice-shop/juice-shop) |

```
Internet ──:80──► agent-dvwa       [agent] ──► [DVWA]        ─┐
Internet ──:80──► agent-juiceshop  [agent] ──► [Juice Shop]  ─┤ logs, alerts
                                                              │ (private network, :8081)
Administrators ──:3000, :8081──► central  [dashboard] ──► [API] ◄─┘
                                                           ├──► [Cassandra]
                                                           └──► [Elasticsearch]
```

- [`terraform/`](terraform) creates the network, the security groups, the virtual machines and the floating IPs, and writes the Ansible inventory.
- [`ansible/`](ansible) installs Docker on the virtual machines, copies the source code, builds the images and starts the containers with Docker Compose.

Every component runs in Docker. The images are built on the virtual machines from the Dockerfiles in `agent/`, `api/` and `dashboard/`.

## Requirements
- An OpenStack project with an external network for the floating IPs and an Ubuntu 22.04 or 24.04 image
- [Terraform](https://developer.hashicorp.com/terraform/install) 1.5 or newer (or OpenTofu)
- Ansible 2.15 or newer (`ansible-core`) with the collections in `ansible/requirements.yml`, and `rsync`
- An SSH key pair (the public key is uploaded to OpenStack)
- The OpenStack credentials, either in `~/.config/openstack/clouds.yaml` or as an OpenStack RC file (`source project-openrc.sh`)

Recommended flavors: 4 vCPUs and 8 GB RAM for `central` (Cassandra and Elasticsearch use 2 GB of heap each), 2 vCPUs and 4 GB RAM for the agents.

## 1. Create the virtual machines
```bash
cd evaluation/terraform
cp terraform.tfvars.example terraform.tfvars   # set the cloud, the external network, the image, the flavors and admin_cidrs
terraform init
terraform apply
```

The main variables (see `variables.tf` for all of them):

| Variable | Default | Description |
|---|---|---|
| `cloud` | `""` | The cloud in `clouds.yaml`; empty to use the `OS_*` environment variables |
| `external_network_name` | (required) | The external network used by the router and the floating IPs |
| `image_name` | `Ubuntu 24.04` | The image of the virtual machines |
| `ssh_user` | `ubuntu` | The user of the image |
| `central_flavor`, `agent_flavor` | `m1.large`, `m1.medium` | The flavors of the virtual machines |
| `boot_from_volume` | `true` | Boot from volumes of `central_volume_size` (60 GB) and `agent_volume_size` (30 GB); set to `false` on clouds without block storage |
| `admin_cidrs` | `["0.0.0.0/0"]` | The networks allowed to reach SSH, the dashboard (3000) and the API (8081). Restrict it to your address |
| `honeypot_ports` | `[80]` | The ports of the agents open to the Internet |
| `agents` | `dvwa` and `juiceshop` | The agent virtual machines and the application protected by each one |

Terraform creates a private network (`10.20.0.0/24`) with a router, two security groups, a key pair, the three virtual machines with floating IPs, and the file `ansible/inventory.ini`. The agents reach the API through the private network; only the honeypot ports of the agents are open to everyone.

`terraform output` shows the addresses of the dashboard, the API and the agents.

## 2. Configure the virtual machines
```bash
cd evaluation/ansible
ansible-galaxy collection install -r requirements.yml
# change admin_password (and the other settings) in group_vars/all.yml
ansible-playbook site.yml
```

The playbook (`site.yml`) runs three roles:

| Role | Hosts | What it does |
|---|---|---|
| `common` | all | Installs Docker and the Compose plugin from the Docker repository |
| `central` | `central` | Copies `api/` and `dashboard/`, writes the API configuration and the Compose file in `/opt/honeypot/central`, builds the images and starts Cassandra, a container which creates the keyspace, Elasticsearch, the API and the dashboard, then waits until the API and the dashboard answer |
| `agent` | `agents` | Copies `agent/`, writes the agent configuration and the Compose file in `/opt/honeypot/agent`, waits until the API is reachable through the private network, builds the agent image and starts the agent and the vulnerable application |

The central virtual machine is configured first, because an agent registers to the API only when it starts. The UUID received by an agent is kept in `/opt/honeypot/agent/config/agent.conf`, and the playbook keeps it when it runs again, so an agent is not registered twice.

Cassandra and Elasticsearch are only reachable from the containers of the central virtual machine (their data is kept in Docker volumes). The vulnerable applications are only reachable through their agent.

The settings in `group_vars/all.yml`:

| Variable | Default | Description |
|---|---|---|
| `admin_username`, `admin_password` | `admin`, `change-me` | The account of the dashboard. Change the password, for example with `ansible-vault encrypt_string` |
| `api_port`, `dashboard_port` | `8081`, `3000` | The ports of the API and the dashboard |
| `api_public_address` | floating IP of `central` | The address of the API used by the browser for the WebSocket connection of the dashboard |
| `cassandra_heap_size`, `elasticsearch_heap_size` | `2G`, `2g` | The memory of the databases |
| `agent_port` | `80` | The port of the agents open to the Internet |
| `agent_operation_mode` | `waf` | `testing`, `waf` or `adaptive` (`adaptive` needs `llm_api_url`) |
| `agent_use_ai_classifier`, `agent_classifier` | `true`, `svc` | The ML classifier (`svc`, `knn` or `random-forest`) |
| `agent_classifier_processes` | `0` | The size of the classifier pool (`0` = one process for each CPU) |
| `llm_api_url` | `""` | The URL of the LLM API (`llm/app.py`), which is not deployed by this playbook because it needs a GPU server with Ollama |
| `applications` | DVWA and Juice Shop | The image and the port of every vulnerable application |

The playbook can be run again after changing the code or the settings: it copies the code, rebuilds the images and recreates the containers. The vulnerable applications are recreated as well, so their state is reset.

## 3. Use the setup
- Dashboard: `http://<central floating IP>:3000` (log in with `admin_username` and `admin_password`)
- Agents: `http://<agent floating IP>` (DVWA: the default login is `admin` / `password`, then press "Create / Reset Database")

Useful commands on the virtual machines:
```bash
cd /opt/honeypot/central && sudo docker compose ps          # state of the central services
cd /opt/honeypot/central && sudo docker compose logs -f api # logs of the API
cd /opt/honeypot/agent && sudo docker compose logs -f agent # logs of the agent
```

## 4. Remove the setup
```bash
cd evaluation/terraform
terraform destroy
```
This deletes the virtual machines with their volumes, so the collected logs are lost. Export them from the dashboard first.

## Notes
- The agents use HTTP. For HTTPS, set `protocol` to `https` in `roles/agent/templates/agent.conf.j2` and mount a certificate and a key in the agent container (the `agent/ssl` directory is not included in the image).
- Elasticsearch stops accepting logs when its disk is more than 90% full, so keep enough free space on `central` for the duration of the deployment.
- The API endpoints used by the agents are not authenticated, which is why the security groups allow port 8081 only from the agents and `admin_cidrs`.
