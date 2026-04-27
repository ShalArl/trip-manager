# Initialize Workload Identity Federation (primarily for deployment with GitHub Actions)

resource "google_iam_workload_identity_pool" "identity_pool" {
  workload_identity_pool_id = var.wip_pool
  display_name              = "GitHub Actions Pool"
  description               = "Identity pool for GitHub Actions deployments"

  depends_on = [google_project_service.services["run.googleapis.com"]]
}

resource "google_iam_workload_identity_pool_provider" "github_provider" {
  workload_identity_pool_id          = google_iam_workload_identity_pool.identity_pool.workload_identity_pool_id
  workload_identity_pool_provider_id = var.wip_provider
  display_name                       = "GitHub Actions Provider"

  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.actor"      = "assertion.actor"
    "attribute.repository" = "assertion.repository"
    "attribute.ref"        = "assertion.ref"
  }

  attribute_condition = "assertion.repository == \"${var.github_repo}\""

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}
