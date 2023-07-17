locals {
  network_name  = var.create_network == true ? module.host-vpc.name : var.network_name
  network_id    = var.create_network == true ? module.host-vpc.id : data.google_compute_network.host_vpc[0].id
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
}
