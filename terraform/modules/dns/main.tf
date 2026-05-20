resource "google_dns_managed_zone" "primary" {
  name        = "trip-manager-zone"
  dns_name    = "${var.domain}."
  description = "Trip Manager DNS Zone"
  project     = var.project_id
}

# resource "google_dns_record_set" "api" {
#   name         = "api.${var.domain}."
#   type         = "A"
#   ttl          = 300
#   managed_zone = google_dns_managed_zone.primary.name
#   project      = var.project_id
#   rrdatas      = [""]  # TODO: Add the IP address of the API server here
# }


