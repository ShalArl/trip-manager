terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "7.7.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "3.7.2"
    }
  }

  backend "gcs" {
    bucket = "project-32c60644-299b-4b05-8cf-tf-state"
    prefix = "terraform/state"
  }

}


