output "name_servers" {
  value = google_dns_managed_zone.primary.name_servers
  description = "List of name servers which are to be added in namecheap config"
}