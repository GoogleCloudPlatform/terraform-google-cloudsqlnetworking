locals {
  network_name  = var.create_network == true ? google_compute_network.private_network[0].name : var.network_name
  network_id    = var.create_network == true ? google_compute_network.private_network[0].id : var.network_id
  subnetwork_id = var.create_subnetwork == true ? google_compute_subnetwork.subnet[0].id : var.subnetwork_id
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
  uservpc_network_id    = var.create_user_vpc_network == true ? google_compute_network.uservpc_private_network[0].id : var.user_network_id
  uservpc_subnetwork_id = var.create_user_vpc_subnetwork == true ? google_compute_subnetwork.uservpc_subnet[0].id : var.user_subnetwork_id
  user_vm_service_account = {
    email  = module.user_gce_sa.email
    scopes = ["cloud-platform"]
  }
}
