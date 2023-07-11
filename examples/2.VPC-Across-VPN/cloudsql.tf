module "sql-db" {
  source              = "../../modules/cloudsql"
  name                = var.cloudsql_instance_name
  database_version    = var.database_version
  zone                = var.zone
  project_id          = var.service_project_id
  ip_configuration    = local.ip_configuration
  deletion_protection = var.deletion_protection
  create_mysql_db     = true
  depends_on = [
    module.host-vpc
  ]
}
