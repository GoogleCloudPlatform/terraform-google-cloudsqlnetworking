module "ha-vpn" {
  source                           = "../../modules/ha-vpn"
  region                           = var.region
  host_project_id                  = var.host_project_id
  host_network_id                  = local.network_id
  private_ip_address               = module.host-vpc[0].psa_ranges["${var.cloudsql_private_range}"].address       //google_compute_global_address.private_ip_address.address
  private_ip_address_prefix_length = module.host-vpc[0].psa_ranges["${var.cloudsql_private_range}"].prefix_length //google_compute_global_address.private_ip_address.prefix_length
  user_project_id                  = var.user_project_id
  user_network_id                  = local.uservpc_network_id
}
