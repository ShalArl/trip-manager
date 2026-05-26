resource "google_dns_managed_zone" "primary" {
  name        = "trip-manager-zone"
  dns_name    = "${var.domain}."
  description = "Trip Manager DNS Zone"
  project     = var.project_id
}

resource "google_dns_record_set" "root" {
  name         = "${var.domain}."
  type         = "A"
  ttl          = 300
  managed_zone = google_dns_managed_zone.primary.name
  project      = var.project_id
  rrdatas      = ["8.233.101.229"]
}

resource "google_dns_record_set" "api" {
  name         = "api.${var.domain}."
  type         = "A"
  ttl          = 300
  managed_zone = google_dns_managed_zone.primary.name
  project      = var.project_id
  rrdatas      = ["8.233.101.229"] # TODO: Add the IP address of the API server here
}

resource "google_dns_record_set" "www" {
  name         = "www.neatnode.xyz."
  type         = "A"
  ttl          = 300
  managed_zone = google_dns_managed_zone.primary.name
  project      = var.project_id
  rrdatas      = ["8.233.101.229"]
}

resource "google_certificate_manager_certificate" "primary" {
  name    = "trip-manager-cert"
  project = var.project_id

  managed {
    domains = [
      var.domain,
      "api.${var.domain}",
      "www.${var.domain}"
    ]
  }
}

resource "google_certificate_manager_certificate_map" "primary" {
  name    = "trip-manager-cert-map"
  project = var.project_id
}

resource "google_certificate_manager_certificate_map_entry" "www" {
  name         = "trip-manager-cert-map-entry-www"
  project      = var.project_id
  map          = google_certificate_manager_certificate_map.primary.name
  certificates = [google_certificate_manager_certificate.primary.id]
  hostname     = "www.${var.domain}"
}

resource "google_certificate_manager_certificate_map_entry" "primary" {
  name         = "trip-manager-cert-map-entry"
  project      = var.project_id
  map          = google_certificate_manager_certificate_map.primary.name
  certificates = [google_certificate_manager_certificate.primary.id]
  hostname     = var.domain
}

resource "google_certificate_manager_certificate_map_entry" "api" {
  name         = "trip-manager-cert-map-entry-api"
  project      = var.project_id
  map          = google_certificate_manager_certificate_map.primary.name
  certificates = [google_certificate_manager_certificate.primary.id]
  hostname     = "api.${var.domain}"
}

