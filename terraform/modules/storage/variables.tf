variable "gcs_bucket" {
  type    = string
  default = "trip-manager-bucket"
}

variable "region" {
  type    = string
  default = "europe-west1"
}

variable "cors_origins" {
  type    = list(string)
  default = ["https://neatnode.xyz", "https://www.neatnode.xyz"]
}
