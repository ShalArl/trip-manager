terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.7.0"
    }
    google-beta = {
      source  = "hashicorp/google-beta"
      version = "7.7.0"
    }
  }
}
locals {
  services = [
    "auth", "social", "users", "trips", "external-secrets", "frontend", "travel-info",
    "feed", "newsletter", "otel-collector"
  ]
}

resource "google_service_account" "services" {
  for_each = toset(local.services)

  account_id   = "${each.value}-sa"
  display_name = "${each.value} Service Account"
  project      = var.project_id
}

# Firestore + Firebase Related Roles
resource "google_project_iam_member" "auth_firebase" {
  member  = "serviceAccount:${google_service_account.services["auth"].email}"
  project = var.project_id
  role    = "roles/firebase.sdkAdminServiceAgent"
}

# GCS Related Roles
resource "google_project_iam_member" "social_firestore" {
  member  = "serviceAccount:${google_service_account.services["social"].email}"
  project = var.project_id
  role    = "roles/datastore.user"
}

resource "google_project_iam_member" "social_storage" {
  member  = "serviceAccount:${google_service_account.services["social"].email}"
  project = var.project_id
  role    = "roles/storage.objectAdmin"
}

resource "google_project_iam_member" "social_token_creator" {
  member  = "serviceAccount:${google_service_account.services["social"].email}"
  project = var.project_id
  role    = "roles/iam.serviceAccountTokenCreator"
}

resource "google_project_iam_member" "locations_storage" {
  member  = "serviceAccount:${google_service_account.services["trips"].email}"
  project = var.project_id
  role    = "roles/storage.objectAdmin"
}

resource "google_project_iam_member" "locations_token_creator" {
  member  = "serviceAccount:${google_service_account.services["trips"].email}"
  project = var.project_id
  role    = "roles/iam.serviceAccountTokenCreator"
}

# General Roles + Infrastructure
resource "google_project_iam_member" "external_secrets_secretmanager" {
  member  = "serviceAccount:${google_service_account.services["external-secrets"].email}"
  project = var.project_id
  role    = "roles/secretmanager.secretAccessor"
}

resource "google_service_account_iam_member" "workload_identity" {
  for_each = toset([
    "auth", "social", "users", "trips", "newsletter", "otel-collector"
  ])
  member             = "serviceAccount:${var.project_id}.svc.id.goog[trip-manager-prod/${each.value}]"
  role               = "roles/iam.workloadIdentityUser"
  service_account_id = google_service_account.services[each.value].name

  depends_on = [var.gke_cluster_id]
}

resource "google_service_account_iam_member" "external_secrets_workload_identity" {
  member             = "serviceAccount:${var.project_id}.svc.id.goog[external-secrets/external-secrets]"
  role               = "roles/iam.workloadIdentityUser"
  service_account_id = google_service_account.services["external-secrets"].name

  depends_on = [var.gke_cluster_id]
}

resource "google_artifact_registry_repository" "trip_manager" {
  project       = var.project_id
  format        = "DOCKER"
  location      = "europe-west1"
  repository_id = "trip-manager"
}

# WIF Pool (GitHub)
resource "google_iam_workload_identity_pool" "github" {
  project                   = var.project_id
  workload_identity_pool_id = "github-pool"
  display_name              = "GitHub Workload Identity Pool"
}

resource "google_iam_workload_identity_pool_provider" "github" {
  project                            = var.project_id
  workload_identity_pool_id          = google_iam_workload_identity_pool.github.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-provider"
  display_name                       = "GitHub Actions Provider"
  description                        = "Workload Identity Provider for GitHub Actions"

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }

  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.actor"      = "assertion.actor"
    "attribute.repository" = "assertion.repository"
  }

  attribute_condition = "assertion.repository == '${var.github_repo}'"
}


resource "google_service_account" "github_actions" {
  account_id   = "github-actions-sa"
  display_name = "GitHub Actions Service Account"
  project      = var.project_id
}

resource "google_project_iam_member" "github_actions_roles" {
  for_each = toset([
    "roles/container.developer",
    "roles/artifactregistry.writer",
    "roles/iam.serviceAccountTokenCreator",
  ])

  project = var.project_id
  role    = each.value
  member  = "serviceAccount:${google_service_account.github_actions.email}"
}

resource "google_service_account_iam_member" "github_wif" {
  service_account_id = google_service_account.github_actions.name
  role               = "roles/iam.workloadIdentityUser"
  member             = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github.name}/attribute.repository/${var.github_repo}"
}

# Artifact Registry Reader Role for all Services otherwise they can't pull the images from the Artifact Registry Repository
resource "google_project_iam_member" "ar_reader" {
  for_each = toset(local.services)

  project = var.project_id
  role    = "roles/artifactregistry.reader"
  member  = "serviceAccount:${each.value}-sa@${var.project_id}.iam.gserviceaccount.com"
}


resource "google_project_iam_member" "trips_pubsub" {
  project = var.project_id
  role    = "roles/pubsub.editor"
  member  = "serviceAccount:${google_service_account.services["trips"].email}"
}

resource "google_project_iam_member" "newsletter_pubsub" {
  project = var.project_id
  role    = "roles/pubsub.editor"
  member  = "serviceAccount:${google_service_account.services["newsletter"].email}"
}


resource "google_project_iam_member" "feed_pub_sub" {
  project = var.project_id
  role    = "roles/pubsub.subscriber"
  member  = "serviceAccount:${google_service_account.services["feed"].email}"
}

resource "google_project_iam_member" "users_firebase" {
  member  = "serviceAccount:${google_service_account.services["users"].email}"
  project = var.project_id
  role    = "roles/firebase.sdkAdminServiceAgent"
}

resource "google_project_iam_member" "otel_monitoring" {
  project = var.project_id
  role    = "roles/monitoring.metricWriter"
  member  = "serviceAccount:${google_service_account.services["otel-collector"].email}"
}

resource "google_project_iam_member" "otel_trace" {
  project = var.project_id
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.services["otel-collector"].email}"
}

resource "google_project_iam_member" "users_monitoring" {
  project = var.project_id
  role    = "roles/monitoring.viewer"
  member  = "serviceAccount:${google_service_account.services["users"].email}"
}