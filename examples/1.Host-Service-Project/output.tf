output "host_vpc_name" {
  value       = local.network_name
  description = "Name of the VPC created in the host project."
}

output "host_network_id" {
  value       = local.network_id
  description = "Network ID for the host VPC network created in the host project."
}

output "host_subnetwork_id" {
  value       = local.subnetwork_id
  description = "Sub Network ID created inside the host VPC network created in the host project."
}

output "cloudsql_instance_name" {
  value       = module.sql-db.mysql_cloudsql_instance_name
  description = "Name of the my cloud sql instance created in the service project."
}

