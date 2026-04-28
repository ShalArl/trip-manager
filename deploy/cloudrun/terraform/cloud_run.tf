resource "google_cloud_run_v2_service" "backend" {
  name     = local.backend_service_name
  location = var.region

  template {
    service_account = google_service_account.runtime_sa.email

    # Limit to 10 instances to save cost, bottleneck is database anyway
    scaling {
      min_instance_count = 0
      max_instance_count = 10
    }

    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      ports {
        container_port = 8081
      }
      resources {
        limits = {
          memory = "512Mi"
          cpu    = "1"
        }
      }

      env {
        name  = "ENVIRONMENT"
        value = "production"
      }

      env {
        name  = "SERVER_PORT"
        value = "8081"
      }

      env {
        name  = "CORS_ALLOWED_ORIGINS"
        value = local.cors_origins
      }

      env {
        name  = "STORAGE_TYPE"
        value = "gcs"
      }

      env {
        name  = "GCS_BUCKET"
        value = google_storage_bucket.uploads.name
      }

      env {
        name  = "GCS_SIGNER_SA"
        value = google_service_account.signed_url_signer.email
      }

      env {
        name  = "SIGNED_URL_TTL_SECONDS"
        value = "900"
      }

      env {
        name  = "FIREBASE_PROJECT_ID"
        value = var.project_id
      }

      env {
        name = "DATABASE_URL"
        value_source {
          secret_key_ref {
            secret  = google_secret_manager_secret.db_url.secret_id
            version = "latest"
          }
        }
      }

      volume_mounts {
        name       = "cloudsql"
        mount_path = "/cloudsql"
      }
    }

    volumes {
      name = "cloudsql"
      cloud_sql_instance {
        instances = [google_sql_database_instance.main.connection_name]
      }
    }
  }

  lifecycle {
    ignore_changes = [
      template[0].containers[0].image,
    ]
  }


  depends_on = [
    google_project_service.services["run.googleapis.com"],
    google_secret_manager_secret_version.db_url_version
  ]

  deletion_protection = false
}

resource "google_cloud_run_v2_service_iam_member" "backend_public" {
  location = google_cloud_run_v2_service.backend.location
  name     = google_cloud_run_v2_service.backend.name
  member   = "allUsers"
  role     = "roles/run.invoker"
}

resource "google_cloud_run_v2_service" "frontend" {
  name     = local.frontend_service_name
  location = var.region

  template {
    service_account = google_service_account.runtime_sa.email
    # Not part of load testing therefore no scaling expected
    scaling {
      min_instance_count = 0
      max_instance_count = 5
    }

    containers {
      image = "us-docker.pkg.dev/cloudrun/container/hello"
      ports {
        container_port = 3000
      }

      resources {
        limits = {
          memory = "512Mi"
          cpu    = "1"
        }
      }

    }

  }

  lifecycle {
    ignore_changes = [
      template[0].containers[0].image,
    ]
  }

  depends_on = [google_project_service.services["run.googleapis.com"]]

  deletion_protection = false
}

resource "google_cloud_run_v2_service_iam_member" "frontend_public" {
  location = google_cloud_run_v2_service.frontend.location
  name     = google_cloud_run_v2_service.frontend.name
  member   = "allUsers"
  role     = "roles/run.invoker"
}
