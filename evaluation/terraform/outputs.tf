output "central_public_ip" {
  description = "The floating IP of the central virtual machine"
  value       = openstack_networking_floatingip_v2.central.address
}

output "dashboard_url" {
  description = "The URL of the dashboard"
  value       = "http://${openstack_networking_floatingip_v2.central.address}:3000"
}

output "api_url" {
  description = "The URL of the API"
  value       = "http://${openstack_networking_floatingip_v2.central.address}:8081/api/v1"
}

output "agents" {
  description = "The agents, the application they protect and their addresses"
  value = {
    for name, agent in var.agents : name => {
      application = agent.application
      public_ip   = openstack_networking_floatingip_v2.agent[name].address
      private_ip  = openstack_networking_port_v2.agent[name].all_fixed_ips[0]
      url         = "http://${openstack_networking_floatingip_v2.agent[name].address}"
    }
  }
}
