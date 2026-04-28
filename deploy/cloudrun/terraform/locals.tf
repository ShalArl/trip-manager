locals {
  api_domain = "api.${var.domain}"
  app_domain = "app.${var.domain}"

  # Cloud Run
  backend_service_name  = "${var.app_name}-backend"
  frontend_service_name = "${var.app_name}-frontend"

  # Domain Mapping
  dm_backend_name = "api.${var.domain}"
  dm_frontend_name = "app.${var.domain}"

  # DNS
  dns_zone_name = "${replace(var.domain, ".", "-")}-zone"
  dns_name      = "${var.domain}."

  # Database
  db_name       = "${var.app_name}-db"
  db_connection = "postgres://${var.db_user}:${random_password.db_password.result}@/${var.db_name}?host=/cloudsql/${google_sql_database_instance.main.connection_name}&sslmode=disable"

  cors_origins = "https://${local.app_domain}"

  # SA
  sa_deploy_member  = "serviceAccount:${google_service_account.deploy_sa.email}"
  sa_runtime_member = "serviceAccount:${google_service_account.runtime_sa.email}"
  sa_signer_member = "serviceAccount:${google_service_account.signed_url_signer.email}"

}
