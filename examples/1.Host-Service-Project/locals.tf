locals {
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
}
