variable "project_id" {
  type        = string
  description = "The ID of the project in which to deploy the Cloud Run service."
}

variable "region" {
  type        = string
  description = "The region in which to deploy the Cloud Run service."
}

variable "app_name" {
  type        = string
  description = "The name of the application to deploy."
}

variable "github_repo" {
  type        = string
  description = "The GitHub repository containing the application code. And from which the Cloud Build will be triggered and or images pushed."
}

variable "wip_pool" {
  type        = string
  description = "The name of the Workload Identity Pool to create for the Cloud Build service account to impersonate and access other GCP resources."
}

variable "wip_provider" {
  type        = string
  description = "The name of the Workload Identity Pool Provider to create for the GitHub Actions OIDC provider to authenticate and access GCP resources."
}

variable "db_tier" {
  type        = string
  description = "The tier of the Cloud SQL instance to create for the application database."
}

variable "db_name" {
  type        = string
  description = "The name of the database to create in the Cloud SQL instance."
}

variable "domain" {
  type        = string
  description = "The custom domain to use for the Cloud Run services. The Terraform configuration will create a DNS managed zone and record sets for the API and frontend subdomains, and map them to the respective Cloud Run services."
}

variable "db_user" {
  type        = string
  description = "The username for the database user to create in the Cloud SQL instance."
}

variable "deploy_sa_name" {
  type        = string
  description = "The name of the service account to create for deployment with GitHub Actions."
}

variable "ar_repo_name" {
  type        = string
  description = "The name of the Artifact Registry repository to create for storing container images."
}

variable "gcs_bucket" {
  type        = string
  description = "The name of the GCS bucket for uploading image files."
}

variable "domain_verification" {
  type = string
  description = "The TXT record value for domain ownership verification from https://search.google.com/search-console/welcome. This is required to set up custom domains with Cloud Run."
}

variable "me" {
  type = string
  description = "Your email address to grant permissions for load testing the deployed application."
}