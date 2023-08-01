output "mysql_cloudsql_instance_name" {
  value       = var.create_mysql_db == true ? module.mysql[0].instance_name : null
  description = "Name of the cloud sql instance created in the service project."
}

output "mysql_cloudsql_instance_details" {
  value       = var.create_mysql_db == true ? module.mysql[0] : null
  description = "Details of the cloud sql instance created in the service project."
}

output "mssql_cloudsql_instance_name" {
  value       = var.create_mssql_db == true ? module.mssql[0].instance_name : null
  description = "Name of the cloud sql instance created in the service project."
}

output "postgres_cloudsql_instance_name" {
  value       = var.create_postgressql_db == true ? module.postgresql[0].instance_name : null
  description = "Name of the cloud sql instance created in the service project."
}
