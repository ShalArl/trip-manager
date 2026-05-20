variable "project_id" {
  type = string
}

variable "iam_sa_email" {
  type = string
}

variable "secret_names" {
  description = "List of secret names to create"
  type        = list(string)
}

variable "secret_values" {
  description = "Map of secret name to value"
  type        = map(string)
  sensitive   = true
}