resource "google_secret_manager_secret" "db_url" {
  secret_id = "database-url"
  replication {
    auto {}
  }

  depends_on = [google_project_service.services["run.googleapis.com"]]
}

resource "google_secret_manager_secret_version" "db_url_version" {
  secret      = google_secret_manager_secret.db_url.id
  secret_data = local.db_connection
}
