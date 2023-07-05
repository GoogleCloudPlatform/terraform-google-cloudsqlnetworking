output "host_network_id" {
  value       = local.network_id
  description = "Network ID for the host VPC network created in the host project."
}

output "host_subnetwork_id" {
  value       = local.subnetwork_id
  description = "Sub Network ID created inside the host VPC network created in the host project."
}

output "uservpc_network_id" {
  value       = local.uservpc_network_id
  description = "Network ID for the User VPC network created in the user project."
}

output "uservpc_subnetwork_id" {
  value       = local.uservpc_subnetwork_id
  description = "Sub Network ID created inside the User VPC network created in the User project."
}

output "private_ip_address" {
  value       = google_compute_global_address.private_ip_address.address
  description = "The IP address or beginning of the address range represented by this resource."
}

output "private_ip_address_prefix_length" {
  value       = google_compute_global_address.private_ip_address.prefix_length
  description = " The prefix length of the IP range. If not present, it means the address field is a single IP address. "
}
