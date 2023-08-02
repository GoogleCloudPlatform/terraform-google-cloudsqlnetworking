locals {
  network_name  = var.create_network == true ? try(module.host-vpc.name,"") : var.network_name
  network_id    = var.create_network == true ? try(module.host-vpc.id,"") : data.google_compute_network.host_vpc[0].id
  subnetwork_id = var.create_subnetwork == true ? try(module.host-vpc.subnet_ids["${var.region}/${var.subnetwork_name}"], "") : data.google_compute_subnetwork.host_vpc_subnetwork[0].id
  vm_service_account = {
    email  = module.gce_sa.email
    scopes = ["cloud-platform"]
  }
  ip_configuration = {
    authorized_networks                           = []
    ipv4_enabled                                  = false
    private_network                               = local.network_id
    require_ssl                                   = null
    allocated_ip_range                            = null
    enable_private_path_for_google_cloud_services = true
  }
  uservpc_network_name  = var.create_user_vpc_network == true ? try(module.user-vpc.name, "") : var.uservpc_network_name
  uservpc_network_id    = var.create_user_vpc_network == true ? try(module.user-vpc.id, "") : data.google_compute_network.user_vpc[0].id
  uservpc_subnetwork_id = var.create_user_vpc_subnetwork == true ? try(module.user-vpc.subnet_ids["${var.user_region}/${var.uservpc_subnetwork_name}"], "") : data.google_compute_subnetwork.user_vpc_subnetwork[0].id
  user_vm_service_account = {
    email  = module.user_gce_sa.email
    scopes = ["cloud-platform"]
  }
  private_ip_address               = module.host-vpc.psa_ranges["${var.cloudsql_private_range_name}"].address
  private_ip_address_prefix_length = module.host-vpc.psa_ranges["${var.cloudsql_private_range_name}"].prefix_length
}
