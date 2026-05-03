resource "hcloud_ssh_key" "deployer" {
  name       = "${var.server_name}-key"
  public_key = var.deployer_ssh_public_key
}

resource "hcloud_firewall" "app" {
  name = "${var.server_name}-fw"

  rule {
    direction  = "in"
    protocol   = "tcp"
    port       = "22"
    source_ips = ["0.0.0.0/0", "::/0"]
  }

  rule {
    direction  = "in"
    protocol   = "tcp"
    port       = "80"
    source_ips = ["0.0.0.0/0", "::/0"]
  }

  rule {
    direction  = "in"
    protocol   = "tcp"
    port       = "443"
    source_ips = ["0.0.0.0/0", "::/0"]
  }
}

resource "hcloud_server" "app" {
  name        = var.server_name
  server_type = var.server_type
  image       = "ubuntu-24.04"
  location    = var.server_location

  ssh_keys     = [hcloud_ssh_key.deployer.id]
  firewall_ids = [hcloud_firewall.app.id]

  user_data = templatefile("${path.module}/cloud-init.yaml.tpl", {
    ssh_public_key        = var.deployer_ssh_public_key
    domain                = var.domain
    backend_image         = var.backend_image
    frontend_image        = var.frontend_image
    firebase_project_id   = var.firebase_project_id
    cors_allowed_origins  = var.cors_allowed_origins
    postgres_image        = var.postgres_image
    postgres_user         = var.postgres_user
    postgres_password     = var.postgres_password
    postgres_db           = var.postgres_db
    github_username       = var.github_username
    github_registry_token = var.github_registry_token
    minio_access_key      = var.minio_access_key
    minio_secret_key      = var.minio_secret_key
  })

  lifecycle {
    ignore_changes = [user_data]
  }

  labels = {
    managed_by = "terraform"
    project    = "trip-manager"
  }
}
