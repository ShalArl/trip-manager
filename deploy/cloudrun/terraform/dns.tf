# New: setup shared

resource "google_dns_managed_zone" "main_zone" {
  dns_name    = local.dns_name
  name        = local.dns_zone_name
  description = "Managed zone for ${var.domain}"

  depends_on = [google_project_service.services["run.googleapis.com"]]
}

# Backend
resource "google_dns_record_set" "api" {
  name         = "${local.api_domain}."
  type         = "CNAME"
  ttl          = 300
  managed_zone = google_dns_managed_zone.main_zone.name
  rrdatas      = ["ghs.googlehosted.com."]
}

# Frontend
resource "google_dns_record_set" "app" {
  managed_zone = google_dns_managed_zone.main_zone.name
  name         = "${local.app_domain}."
  type         = "CNAME"
  ttl          = "300"
  rrdatas      = ["ghs.googlehosted.com."]
}

resource "google_dns_record_set" "google_verification" {
  name         = "${var.domain}."
  managed_zone = google_dns_managed_zone.main_zone.name
  type         = "TXT"
  ttl          = 300
  rrdatas      = ["\"google-site-verification=${var.domain_verification}\""]
}

# IAAS

resource "google_dns_record_set" "hetzner" {
  name         = "hetzner.${var.domain}."
  managed_zone = google_dns_managed_zone.main_zone.name
  type         = "A"
  ttl          = 300
  rrdatas      = ["167.235.66.0"]
}

resource "google_dns_record_set" "hetzner_app" {
  name         = "hetzner-app.${var.domain}."
  managed_zone = google_dns_managed_zone.main_zone.name
  type         = "A"
  ttl          = 300
  rrdatas      = ["167.235.66.0"]
}

resource "google_dns_record_set" "hetzner_storage" {
  name         = "hetzner-storage.${var.domain}."
  managed_zone = google_dns_managed_zone.main_zone.name
  type         = "A"
  ttl          = 300
  rrdatas      = ["167.235.66.0"]
}