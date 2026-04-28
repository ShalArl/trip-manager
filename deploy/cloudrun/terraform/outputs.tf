output "wif_provider" {
  description = "Workload Identity Provider name (für GitHub Actions Secret GCP_WIF_PROVIDER)"
  value       = google_iam_workload_identity_pool_provider.github_provider.name
}

output "deploy_sa_email" {
  description = "Deploy Service Account email (für GitHub Actions Secret GCP_DEPLOY_SA)"
  value       = google_service_account.deploy_sa.email
}

output "runtime_sa_email" {
  description = "Runtime Service Account email (für GitHub Actions Secret GCP_RUNTIME_SA)"
  value       = google_service_account.runtime_sa.email
}

output "signed_url_signer_email" {
  description = "Signed URL Signer Service Account email (für GitHub Actions Secret GCP_SIGNED_URL_SIGNER_SA)"
  value       = google_service_account.signed_url_signer
}

# Artifact Registry outputs for GitHub Actions secrets

output "ar_host" {
  description = "Artifact Registry host (für GitHub Actions Secret GCP_AR_HOST)"
  value       = "${var.region}-docker.pkg.dev"
}

output "ar_repo" {
  description = "Artifact Registry repo name (für GitHub Actions Secret GCP_AR_REPO)"
  value       = google_artifact_registry_repository.main.repository_id
}

output "backend_service_name" {
  value = google_cloud_run_v2_service.backend.name
}

output "frontend_service_name" {
  value = google_cloud_run_v2_service.frontend.name
}

output "sql_connection_name" {
  value = google_sql_database_instance.main.connection_name
}

output "gcs_bucket" {
  value = google_storage_bucket.uploads.name
}

output "dns_nameservers" {
  value = google_dns_managed_zone.main_zone.name_servers
}

output "backend_url" {
  value = "https://api.${var.domain}"
}

output "frontend_url" {
  value = "https://app.${var.domain}"
}