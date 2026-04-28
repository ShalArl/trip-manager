resource "google_cloud_run_domain_mapping" "backend" {
  location = var.region
  name     = local.dm_backend_name

  metadata {
    namespace = var.project_id
  }

  spec {
    route_name = google_cloud_run_v2_service.backend.name
  }
}

resource "google_cloud_run_domain_mapping" "frontend" {
  location = var.region
  name     = local.dm_frontend_name

  metadata {
    namespace = var.project_id
  }

  spec {
    route_name = google_cloud_run_v2_service.frontend.name
  }
}
