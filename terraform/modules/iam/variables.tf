variable "project_id" {
  type = string
}

variable "github_repo" {
  type = string
}

variable "gke_cluster_id" {
  description = "GKE Cluster ID – verify resource exists before applying"
  type        = string
}