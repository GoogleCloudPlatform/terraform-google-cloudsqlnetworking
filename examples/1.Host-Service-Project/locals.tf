locals {
  network_name  = var.create_network == true ? module.host-vpc[0].name : var.network_name
  network_id    = var.create_network == true ? module.host-vpc[0].id : var.network_id
  subnetwork_id = var.create_subnetwork == true ? module.host-vpc[0].subnet_ids["${var.region}/${var.subnetwork_name}"] : var.subnetwork_id
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
