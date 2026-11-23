variable "cloud" {
  description = "The name of the cloud in clouds.yaml. Leave empty to use the OS_* environment variables (OpenStack RC file)"
  type        = string
  default     = ""
}

variable "name_prefix" {
  description = "Prefix added to the names of all the created resources"
  type        = string
  default     = "honeypot"
}

variable "image_name" {
  description = "The image of the virtual machines (the Ansible playbooks support Ubuntu 22.04 and 24.04)"
  type        = string
  default     = "Ubuntu 24.04"
}

variable "ssh_user" {
  description = "The user used by Ansible to connect to the virtual machines (depends on the image)"
  type        = string
  default     = "ubuntu"
}

variable "central_flavor" {
  description = "The flavor of the central virtual machine (API, Cassandra, Elasticsearch and dashboard). At least 4 vCPUs and 8 GB RAM are recommended"
  type        = string
  default     = "m1.large"
}

variable "agent_flavor" {
  description = "The flavor of the agent virtual machines (agent and vulnerable application). At least 2 vCPUs and 4 GB RAM are recommended"
  type        = string
  default     = "m1.medium"
}

variable "central_volume_size" {
  description = "The size in GB of the boot volume of the central virtual machine (it stores the logs)"
  type        = number
  default     = 60
}

variable "agent_volume_size" {
  description = "The size in GB of the boot volumes of the agent virtual machines"
  type        = number
  default     = 30
}

variable "external_network_name" {
  description = "The name of the external (public) network, used for the router and the floating IPs"
  type        = string
}

variable "subnet_cidr" {
  description = "The CIDR of the private subnet created for the virtual machines"
  type        = string
  default     = "10.20.0.0/24"
}

variable "dns_nameservers" {
  description = "The DNS servers of the private subnet"
  type        = list(string)
  default     = ["1.1.1.1", "8.8.8.8"]
}

variable "public_key_path" {
  description = "The path of the SSH public key uploaded as a key pair"
  type        = string
  default     = "~/.ssh/id_ed25519.pub"
}

variable "private_key_path" {
  description = "The path of the SSH private key used by Ansible (written in the generated inventory)"
  type        = string
  default     = "~/.ssh/id_ed25519"
}

variable "admin_cidrs" {
  description = "The networks allowed to access SSH, the dashboard and the API (the dashboard opens a WebSocket connection to the API from the browser)"
  type        = list(string)
  default     = ["0.0.0.0/0"]
}

variable "honeypot_ports" {
  description = "The ports of the agents open to the Internet"
  type        = list(number)
  default     = [80]
}

variable "agents" {
  description = "The agent virtual machines and the vulnerable application protected by each one (dvwa or juiceshop)"
  type = map(object({
    application = string
  }))
  default = {
    dvwa = {
      application = "dvwa"
    }
    juiceshop = {
      application = "juiceshop"
    }
  }

  validation {
    condition     = alltrue([for agent in values(var.agents) : contains(["dvwa", "juiceshop"], agent.application)])
    error_message = "The application of an agent must be dvwa or juiceshop."
  }
}

variable "boot_from_volume" {
  description = "Boot the virtual machines from volumes of the given sizes. Set to false on clouds without block storage (the disk of the flavor is used)"
  type        = bool
  default     = true
}
