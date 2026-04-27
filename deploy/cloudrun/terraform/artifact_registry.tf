resource "google_artifact_registry_repository" "main" {
  location      = var.region
  repository_id = var.ar_repo_name
  description   = "Docker images for ${var.app_name}"
  format        = "DOCKER"

  depends_on = [google_project_service.services["run.googleapis.com"]]
}
