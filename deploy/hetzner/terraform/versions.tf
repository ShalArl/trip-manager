terraform {
  required_providers {
    hcloud = {
      source  = "hetznercloud/hcloud"
      version = "1.37.0"
    }

    random = {
      source  = "hashicorp/random"
      version = "3.7.2"
    }
  }
}
variable "hcloud_token" {
  sensitive = true
}

provider "hcloud" {
  token = var.hcloud_token
}

resource "hcloud_ssh_key" "main_key" {
  name       = "deploy-key"
  public_key = "ssh-rsa AAAAB3Nza... dein_key_hier ..."
}

resource "hcloud_server" "web_server" {
  name        = "mein-vserver"
  server_type = "cx23"
  image       = "ubuntu-24.04"
  location    = "nbg1"

  ssh_keys = [hcloud_ssh_key.main_key.id]

  # Da du Docker-Compose nutzt, stellen wir sicher, dass Docker da ist
  user_data = <<-EOT
    #cloud-config
    packages:
      - docker.io
      - docker-compose
    runcmd:
      - systemctl enable --now docker
  EOT

  # Verhindert, dass Terraform den Server löscht, falls du dich beim Import vertippst
  lifecycle {
    prevent_destroy = false
  }
}

# Output der IP, damit du siehst, ob sie sich nach einem Apply ändert
output "server_ip" {
  value = hcloud_server.web_server.ipv4_address
}