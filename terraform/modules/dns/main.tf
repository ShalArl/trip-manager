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
  rrdatas      = ["8.232.170.119"]
}

resource "google_dns_record_set" "api" {
  name         = "api.${var.domain}."
  type         = "A"
  ttl          = 300
  managed_zone = google_dns_managed_zone.primary.name
  project      = var.project_id
  rrdatas      = ["8.232.170.119"] # TODO: Add the IP address of the API server here
}

resource "google_dns_record_set" "www" {
  name         = "www.neatnode.xyz."
  type         = "A"
  ttl          = 300
  managed_zone = google_dns_managed_zone.primary.name
  project      = var.project_id
  rrdatas      = ["8.232.170.119"]
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

# Resend related
resource "google_dns_record_set" "resend_dkim" {
  name         = "resend._domainkey.neatnode.xyz."
  type         = "TXT"
  ttl          = 300
  managed_zone = google_dns_managed_zone.primary.name
  project      = var.project_id
  rrdatas      = ["\"v=DKIM1; k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCdhAAXtqNesK3awbDhBarUs3xTzVOP9iP2fjv+KsfMV/E8Jkz+YIDo6+xXVEJg+rMeU4ERNj29GyH9cQ0HzkuLQ3mc8NbIaNDo02FjWZI2n1LGsVLkhmOmwgqJCzVY/kBkmlfi1K2yFstT+29BPaCB07LhWqz3m7bmL7rKWUBWiQIDAQAB\""]
}

resource "google_dns_record_set" "resend_mx" {
  name         = "send.neatnode.xyz."
  type         = "MX"
  ttl          = 300
  managed_zone = google_dns_managed_zone.primary.name
  project      = var.project_id
  rrdatas      = ["10 feedback-smtp.eu-west-1.amazonses.com."]
}

resource "google_dns_record_set" "resend_spf" {
  name         = "send.neatnode.xyz."
  type         = "TXT"
  ttl          = 300
  managed_zone = google_dns_managed_zone.primary.name
  project      = var.project_id
  rrdatas      = ["\"v=spf1 include:amazonses.com ~all\""]
}

resource "google_dns_record_set" "resend_dmarc" {
  name         = "_dmarc.neatnode.xyz."
  type         = "TXT"
  ttl          = 300
  managed_zone = google_dns_managed_zone.primary.name
  project      = var.project_id
  rrdatas      = ["\"v=DMARC1; p=none;\""]
}

# ── Staging ────────────────────────────────────────────────────────────────────

resource "google_compute_global_address" "staging" {
  name    = "trip-manager-staging-ip"
  project = var.project_id
}

resource "google_dns_record_set" "staging" {
  name         = "staging.${var.domain}."
  type         = "A"
  ttl          = 300
  managed_zone = google_dns_managed_zone.primary.name
  project      = var.project_id
  rrdatas      = [google_compute_global_address.staging.address]
}

resource "google_dns_record_set" "api_staging" {
  name         = "api.staging.${var.domain}."
  type         = "A"
  ttl          = 300
  managed_zone = google_dns_managed_zone.primary.name
  project      = var.project_id
  rrdatas      = [google_compute_global_address.staging.address]
}

resource "google_certificate_manager_certificate" "staging" {
  name    = "trip-manager-staging-cert"
  project = var.project_id

  managed {
    domains = [
      "staging.${var.domain}",
      "api.staging.${var.domain}",
    ]
  }
}

resource "google_certificate_manager_certificate_map" "staging" {
  name    = "trip-manager-staging-cert-map"
  project = var.project_id
}

resource "google_certificate_manager_certificate_map_entry" "staging" {
  name         = "trip-manager-staging-cert-map-entry"
  project      = var.project_id
  map          = google_certificate_manager_certificate_map.staging.name
  certificates = [google_certificate_manager_certificate.staging.id]
  hostname     = "staging.${var.domain}"
}

resource "google_certificate_manager_certificate_map_entry" "api_staging" {
  name         = "trip-manager-staging-cert-map-entry-api"
  project      = var.project_id
  map          = google_certificate_manager_certificate_map.staging.name
  certificates = [google_certificate_manager_certificate.staging.id]
  hostname     = "api.staging.${var.domain}"
}

output "staging_ip" {
  value = google_compute_global_address.staging.address
}