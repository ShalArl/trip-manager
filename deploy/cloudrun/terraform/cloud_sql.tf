# Setup Postgresql

resource "google_sql_database_instance" "main" {
  name             = "${var.app_name}-db"
  region           = var.region
  database_version = "POSTGRES_16"
  settings {
    tier                        = var.db_tier
    deletion_protection_enabled = false
    disk_size                   = 10
    disk_type                   = "PD_SSD"
    edition                     = "ENTERPRISE"
    availability_type           = "ZONAL"
    
  }
  deletion_protection = false

  depends_on = [google_project_service.services["run.googleapis.com"]]
}

resource "google_sql_database" "app_db" {
  name     = var.db_name
  instance = google_sql_database_instance.main.name
  deletion_policy = ""
}

resource "google_sql_user" "db_user" {
  name     = var.db_user
  instance = google_sql_database_instance.main.name
  password = random_password.db_password.result
}

resource "random_password" "db_password" {
  length  = 16
  special = false
}

