variable "project_id" {
  description = "Project ID of the GCP project"
  type        = string
  default     = "project-32c60644-299b-4b05-8cf"
}

variable "region" {
  description = "Region to deploy resources in"
  type        = string
  default     = "europe-west1"
}

variable "domain" {
  description = "Domain name for the application"
  type        = string
  default     = "neatnode.xyz"
}

variable "environment" {
  description = "Deployment environment (e.g., dev, staging, prod)"
  type        = string
  default     = "prod"
}

variable "secrets" {
    description = "Map of secret names to their values"
    type        = map(string)
    default     = {}
}

variable "github_repo" {
  description = "GitHub repository in format owner/repo e.g. ShalArl/trip-manager"
  type        = string
  default     = "ShalArl/trip-manager"
}
