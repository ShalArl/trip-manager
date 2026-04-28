# IAM Resources for Cloud Run deployment
resource "google_service_account" "deploy_sa" {
  account_id   = var.deploy_sa_name
  display_name = "Deploy Service Account"
  description  = "Used by GitHub Actions to deploy to Cloud Run"
}

resource "google_project_iam_member" "deploy_sa_artifact_writer" {
  project = var.project_id
  role    = "roles/artifactregistry.writer"
  member  = local.sa_deploy_member
}

resource "google_project_iam_member" "deploy_sa_run_admin" {
  project = var.project_id
  role    = "roles/run.admin"
  member  = local.sa_deploy_member
}

resource "google_project_iam_member" "deploy_sa_service_account_user" {
  project = var.project_id
  role    = "roles/iam.serviceAccountUser"
  member  = local.sa_deploy_member
}

resource "google_service_account_iam_member" "deploy_sa_workload_identity" {
  service_account_id = google_service_account.deploy_sa.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/projects/${data.google_project.project.number}/locations/global/workloadIdentityPools/${google_iam_workload_identity_pool.identity_pool.workload_identity_pool_id}/attribute.repository/${var.github_repo}"
}

resource "google_service_account" "runtime_sa" {
  account_id   = "${var.app_name}-runtime"
  display_name = "Runtime Service Account"
}

resource "google_project_iam_member" "runtime_cloudsql_client" {
  project = var.project_id
  role    = "roles/cloudsql.client"
  member  = local.sa_runtime_member
}

resource "google_storage_bucket_iam_member" "runtime_object_admin" {
  bucket = google_storage_bucket.uploads.name
  role   = "roles/storage.objectAdmin"
  member = local.sa_runtime_member
}

resource "google_service_account_iam_member" "runtime_can_impersonate_signer" {
  service_account_id = google_service_account.signed_url_signer.name
  role               = "roles/iam.serviceAccountTokenCreator"
  member             = local.sa_runtime_member
}

resource "google_project_iam_member" "runtime_secret_accessor" {
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
  member  = local.sa_runtime_member
}

resource "google_project_iam_member" "runtime_firebase_auth" {
  project = var.project_id
  role    = "roles/firebaseauth.admin"
  member  = local.sa_runtime_member
}

resource "google_project_iam_member" "runtime_firestore" {
  project = var.project_id
  role    = "roles/datastore.user"
  member  = local.sa_runtime_member
}
