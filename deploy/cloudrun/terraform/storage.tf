resource "google_storage_bucket" "uploads" {
  location = var.region
  name     = var.gcs_bucket

  force_destroy               = true
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  cors {
    origin          = [local.cors_origins]
    method          = ["GET", "PUT"]
    response_header = ["Content-Type", "x-goog-resumable"]
    max_age_seconds = 3600
  }

  depends_on = [google_project_service.services["run.googleapis.com"]]
}

resource "google_service_account" "signed_url_signer" {
  account_id   = "signed-url-signer"
  display_name = "Service Account for signing GCS URLs"
}

resource "google_storage_bucket_iam_member" "signer_object_admin" {
  bucket = google_storage_bucket.uploads.name
  role   = "roles/storage.objectAdmin"
  member = "serviceAccount:${google_service_account.signed_url_signer.email}"
}