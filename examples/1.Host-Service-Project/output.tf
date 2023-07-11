output "host_vpc_name" {
  value       = module.host-vpc[0].name
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


