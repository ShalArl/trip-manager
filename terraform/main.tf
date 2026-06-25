resource "google_project_service" "apis" {
  for_each = toset([
    "monitoring.googleapis.com",
    "cloudtrace.googleapis.com",
    "logging.googleapis.com",
    "container.googleapis.com",
    "secretmanager.googleapis.com",
    "dns.googleapis.com",
    "artifactregistry.googleapis.com",
    "iam.googleapis.com",
    "cloudresourcemanager.googleapis.com",
    "compute.googleapis.com",
    "certificatemanager.googleapis.com",
    "pubsub.googleapis.com",
    "firestore.googleapis.com",

  ])

  project            = var.project_id
  service            = each.value
  disable_on_destroy = false
}

// Comment that block out to save costs
module "gke" {
   source      = "./modules/gke"
   project_id  = var.project_id
   region      = var.region
   environment = var.environment
   depends_on = [
     google_project_service.apis
   ]
 }

module "iam" {
  source      = "./modules/iam"
  project_id  = var.project_id
  github_repo = var.github_repo
  gke_cluster_id = "gke_project-32c60644-299b-4b05-8cf_europe-west1_trip-manager-prod"
  depends_on = [
    google_project_service.apis
  ]
}

module "storage" {
  source = "./modules/storage"
}

module "dns" {
  source     = "./modules/dns"
  project_id = var.project_id
  domain     = var.domain
  depends_on = [google_project_service.apis]
}

module "secrets" {
  source        = "./modules/secrets"
  project_id    = var.project_id
  iam_sa_email  = module.iam.external_secrets_sa_email
  secret_names  = keys(var.secrets)
  secret_values = var.secrets
  depends_on    = [google_project_service.apis]
}

module "pubsub" {
  source     = "./modules/pubsub"
  project_id = var.project_id
  depends_on = [
    google_project_service.apis
  ]
}

module "firestore" {
  source     = "./modules/firestore"
  project_id = var.project_id
  depends_on = [
    google_project_service.apis
  ]
}
