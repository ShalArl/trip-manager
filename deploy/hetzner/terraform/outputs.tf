output "server_ipv4" {
  value = hcloud_server.app.ipv4_address
}

output "server_ipv6" {
  value = hcloud_server.app.ipv6_address
}

output "server_id" {
  value = hcloud_server.app.id
}

output "ssh_command" {
  value = "ssh deployer@${hcloud_server.app.ipv4_address}"
}

output "deploy_command" {
  description = "Run this to deploy a new image version"
  value       = "ssh deployer@${hcloud_server.app.ipv4_address} '/home/deployer/scripts/update.sh'"
}