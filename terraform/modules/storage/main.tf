resource "google_storage_bucket" "uploads" {
  location = var.region
  name     = var.gcs_bucket

  force_destroy               = true
  uniform_bucket_level_access = true
  public_access_prevention    = "enforced"

  cors {
    origin          = var.cors_origins
    method          = ["GET", "PUT", "POST", "DELETE", "HEAD"]
    response_header = ["Content-Type", "Authorization"]
    max_age_seconds = 3600
  }

}