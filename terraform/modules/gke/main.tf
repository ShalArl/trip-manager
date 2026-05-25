resource "google_compute_network" "vpc" {
  name                    = "trip-manager-vpc"
  auto_create_subnetworks = false
  project                 = var.project_id
}

resource "google_compute_subnetwork" "subnet" {
  name          = "trip-manager-subnet"
  ip_cidr_range = "10.0.0.0/18"
  region        = var.region
  network       = google_compute_network.vpc.id
  project       = var.project_id

  secondary_ip_range {
    range_name    = "pods"
    ip_cidr_range = "10.48.0.0/14"
  }

  secondary_ip_range {
    range_name    = "services"
    ip_cidr_range = "10.52.0.0/20"
  }
}

resource "google_container_cluster" "primary" {
  provider = google-beta

  name     = "trip-manager-${var.environment}"
  location = var.region   # Autopilot ist immer regional
  project  = var.project_id

  enable_autopilot = true

  network    = google_compute_network.vpc.name
  subnetwork = google_compute_subnetwork.subnet.name

  ip_allocation_policy {
    cluster_secondary_range_name  = "pods"
    services_secondary_range_name = "services"
  }

  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  release_channel {
    channel = "REGULAR"
  }

  gateway_api_config {
    channel = "CHANNEL_STANDARD"
  }

  deletion_protection = false
}