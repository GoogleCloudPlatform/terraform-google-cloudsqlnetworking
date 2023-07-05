module "cloudsql" {
  source              = "GoogleCloudPlatform/sql-db/google//modules/mysql"
  version             = "8.0.0"
  name                = var.cloudsql_instance_name
  database_version    = var.database_version
  zone                = var.zone
  project_id          = var.service_project_id
  ip_configuration    = local.ip_configuration
  deletion_protection = false
  depends_on = [
    google_service_networking_connection.private_vpc_connection
  ]
}