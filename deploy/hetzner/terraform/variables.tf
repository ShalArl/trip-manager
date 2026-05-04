variable "hcloud_token" {
  type        = string
  description = "Hetzner Cloud API token"
  sensitive   = true
}

variable "deployer_ssh_public_key" {
  type        = string
  description = "Public SSH key for deployer user"
}

variable "server_name" {
  type    = string
  default = "trip-manager-iaas"
}

variable "server_type" {
  type    = string
  default = "cx23"
}

variable "server_location" {
  type    = string
  default = "nbg1"
}

variable "domain" {
  type        = string
  description = "Domain name for the server (e.g. iaas.neatnode.xyz)"
}

# === App-Konfiguration ===

variable "backend_image" {
  type        = string
  description = "Docker image for backend (e.g. ghcr.io/user/trip-manager/backend:latest)"
}

variable "frontend_image" {
  type        = string
  description = "Docker image for frontend (e.g. ghcr.io/user/trip-manager/frontend:hetzner-latest)"
}

variable "firebase_project_id" {
  type    = string
  default = "project-32c60644-299b-4b05-8cf"
}

variable "cors_allowed_origins" {
  type    = string
  default = "https://app.neatnode.xyz"
}

# === Postgres ===

variable "postgres_image" {
  type    = string
  default = "postgres:16-alpine"
}

variable "postgres_user" {
  type    = string
  default = "tripmanager"
}

variable "postgres_db" {
  type    = string
  default = "tripmanager"
}

variable "postgres_password" {
  type        = string
  description = "Postgres password (generate via: openssl rand -base64 32)"
  sensitive   = true
}

# === Infrastructure related ===

variable "github_registry_token" {
  type        = string
  description = "GitHub PAT for pulling private images from ghcr.io (optional)"
  sensitive   = true
  default     = ""
}

variable "github_username" {
  type        = string
  description = "GitHub username for ghcr.io login"
  default     = ""
}

variable "minio_access_key" {
  type        = string
  description = "MinIO root user / access key"
  sensitive   = true
}

variable "minio_secret_key" {
  type        = string
  description = "MinIO root password / secret key"
  sensitive   = true
}
