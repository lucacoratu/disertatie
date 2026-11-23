data "openstack_networking_network_v2" "external" {
  name = var.external_network_name
}

data "openstack_images_image_v2" "image" {
  name        = var.image_name
  most_recent = true
}

# ---------------------------------------------------------------- network
resource "openstack_networking_network_v2" "private" {
  name           = "${var.name_prefix}-network"
  admin_state_up = true
}

resource "openstack_networking_subnet_v2" "private" {
  name            = "${var.name_prefix}-subnet"
  network_id      = openstack_networking_network_v2.private.id
  cidr            = var.subnet_cidr
  ip_version      = 4
  dns_nameservers = var.dns_nameservers
}

resource "openstack_networking_router_v2" "router" {
  name                = "${var.name_prefix}-router"
  admin_state_up      = true
  external_network_id = data.openstack_networking_network_v2.external.id
}

resource "openstack_networking_router_interface_v2" "router" {
  router_id = openstack_networking_router_v2.router.id
  subnet_id = openstack_networking_subnet_v2.private.id
}

# ---------------------------------------------------------------- security groups
resource "openstack_networking_secgroup_v2" "central" {
  name        = "${var.name_prefix}-central"
  description = "Central virtual machine: SSH, dashboard and API for the administrators, API for the agents"
}

resource "openstack_networking_secgroup_rule_v2" "central_admin" {
  for_each = {
    for pair in setproduct(var.admin_cidrs, [22, 3000, 8081]) : "${pair[0]}-${pair[1]}" => pair
  }
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = each.value[1]
  port_range_max    = each.value[1]
  remote_ip_prefix  = each.value[0]
  security_group_id = openstack_networking_secgroup_v2.central.id
}

# The agents send the logs to the API through the private network
resource "openstack_networking_secgroup_rule_v2" "central_agents" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 8081
  port_range_max    = 8081
  remote_group_id   = openstack_networking_secgroup_v2.agent.id
  security_group_id = openstack_networking_secgroup_v2.central.id
}

resource "openstack_networking_secgroup_v2" "agent" {
  name        = "${var.name_prefix}-agent"
  description = "Agent virtual machines: SSH for the administrators, honeypot ports for everyone"
}

resource "openstack_networking_secgroup_rule_v2" "agent_ssh" {
  for_each          = toset(var.admin_cidrs)
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 22
  port_range_max    = 22
  remote_ip_prefix  = each.value
  security_group_id = openstack_networking_secgroup_v2.agent.id
}

resource "openstack_networking_secgroup_rule_v2" "agent_honeypot" {
  for_each          = toset([for port in var.honeypot_ports : tostring(port)])
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = tonumber(each.value)
  port_range_max    = tonumber(each.value)
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.agent.id
}

# ---------------------------------------------------------------- key pair
resource "openstack_compute_keypair_v2" "key" {
  name       = "${var.name_prefix}-key"
  public_key = file(pathexpand(var.public_key_path))
}

# ---------------------------------------------------------------- ports (fixed private addresses)
resource "openstack_networking_port_v2" "central" {
  name               = "${var.name_prefix}-central"
  network_id         = openstack_networking_network_v2.private.id
  security_group_ids = [openstack_networking_secgroup_v2.central.id]

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.private.id
  }
}

resource "openstack_networking_port_v2" "agent" {
  for_each           = var.agents
  name               = "${var.name_prefix}-agent-${each.key}"
  network_id         = openstack_networking_network_v2.private.id
  security_group_ids = [openstack_networking_secgroup_v2.agent.id]

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.private.id
  }
}

# ---------------------------------------------------------------- virtual machines
resource "openstack_compute_instance_v2" "central" {
  name        = "${var.name_prefix}-central"
  flavor_name = var.central_flavor
  image_id    = var.boot_from_volume ? null : data.openstack_images_image_v2.image.id
  key_pair    = openstack_compute_keypair_v2.key.name

  dynamic "block_device" {
    for_each = var.boot_from_volume ? [1] : []
    content {
      uuid                  = data.openstack_images_image_v2.image.id
      source_type           = "image"
      destination_type      = "volume"
      volume_size           = var.central_volume_size
      boot_index            = 0
      delete_on_termination = true
    }
  }

  network {
    port = openstack_networking_port_v2.central.id
  }

  depends_on = [openstack_networking_router_interface_v2.router]
}

resource "openstack_compute_instance_v2" "agent" {
  for_each    = var.agents
  name        = "${var.name_prefix}-agent-${each.key}"
  flavor_name = var.agent_flavor
  image_id    = var.boot_from_volume ? null : data.openstack_images_image_v2.image.id
  key_pair    = openstack_compute_keypair_v2.key.name

  dynamic "block_device" {
    for_each = var.boot_from_volume ? [1] : []
    content {
      uuid                  = data.openstack_images_image_v2.image.id
      source_type           = "image"
      destination_type      = "volume"
      volume_size           = var.agent_volume_size
      boot_index            = 0
      delete_on_termination = true
    }
  }

  network {
    port = openstack_networking_port_v2.agent[each.key].id
  }

  depends_on = [openstack_networking_router_interface_v2.router]
}

# ---------------------------------------------------------------- floating IPs
resource "openstack_networking_floatingip_v2" "central" {
  pool = var.external_network_name
}

resource "openstack_networking_floatingip_associate_v2" "central" {
  floating_ip = openstack_networking_floatingip_v2.central.address
  port_id     = openstack_networking_port_v2.central.id
  depends_on  = [openstack_compute_instance_v2.central]
}

resource "openstack_networking_floatingip_v2" "agent" {
  for_each = var.agents
  pool     = var.external_network_name
}

resource "openstack_networking_floatingip_associate_v2" "agent" {
  for_each    = var.agents
  floating_ip = openstack_networking_floatingip_v2.agent[each.key].address
  port_id     = openstack_networking_port_v2.agent[each.key].id
  depends_on  = [openstack_compute_instance_v2.agent]
}

# ---------------------------------------------------------------- Ansible inventory
resource "local_file" "inventory" {
  filename        = "${path.module}/../ansible/inventory.ini"
  file_permission = "0644"
  content = templatefile("${path.module}/inventory.tftpl", {
    ssh_user           = var.ssh_user
    private_key_path   = pathexpand(var.private_key_path)
    central_public_ip  = openstack_networking_floatingip_v2.central.address
    central_private_ip = openstack_networking_port_v2.central.all_fixed_ips[0]
    agents = {
      for name, agent in var.agents : name => {
        public_ip   = openstack_networking_floatingip_v2.agent[name].address
        private_ip  = openstack_networking_port_v2.agent[name].all_fixed_ips[0]
        application = agent.application
      }
    }
  })
}
