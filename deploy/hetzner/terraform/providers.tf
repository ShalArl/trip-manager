terraform {
  required_version = ">= 1.5"

  required_providers {
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = "~> 1.45"
    }
  }

  backend "gcs" {
    bucket = "project-32c60644-299b-4b05-8cf-tf-state"
    prefix = "terraform/state/hetzner"
  }
}

provider "hcloud" {
  token = var.hcloud_token
}