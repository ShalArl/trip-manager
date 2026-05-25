terraform {
  backend "gcs" {
    bucket = "project-32c60644-299b-4b05-8cf-terraform-state"
    prefix = "terraform/state"
  }
}