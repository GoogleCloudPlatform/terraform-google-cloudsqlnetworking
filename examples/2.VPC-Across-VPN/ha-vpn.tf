module "ha-vpn" {
  source                           = "../../modules/ha-vpn"
  region                           = var.region
  host_project_id                  = var.host_project_id
  host_network_id                  = local.network_id
  private_ip_address               = module.host-vpc.psa_ranges["${var.cloudsql_private_range_name}"].address       //google_compute_global_address.private_ip_address.address
  private_ip_address_prefix_length = module.host-vpc.psa_ranges["${var.cloudsql_private_range_name}"].prefix_length //google_compute_global_address.private_ip_address.prefix_length
  user_project_id                  = var.user_project_id
  user_network_id                  = local.uservpc_network_id
}
