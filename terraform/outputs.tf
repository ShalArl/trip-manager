output "gke_cluster_name" {
  value = module.gke.cluster_name
}

output "name_servers" {
  value = module.dns.name_servers
  description = "To be added in namecheap config"
}

output "service_accounts" {
  value = module.iam.service_account_emails
}