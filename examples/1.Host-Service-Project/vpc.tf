
resource "google_compute_network" "private_network" {
  count                   = var.create_network == true ? 1 : 0
  provider                = google-beta
  name                    = var.network_name
  auto_create_subnetworks = false
  project                 = var.host_project_id
  routing_mode            = var.network_routing_mode
  depends_on = [
    module.host_project,
    module.project_services
  ]
}

resource "google_compute_subnetwork" "subnet" {
  count         = var.create_subnetwork == true ? 1 : 0
  name          = var.subnetwork_name
  project       = var.host_project_id
  ip_cidr_range = var.subnetwork_ip_cidr
  region        = var.region
  network       = local.network_id
  log_config {
    aggregation_interval = "INTERVAL_10_MIN"
    flow_sampling        = 0.5
    metadata             = "INCLUDE_ALL_METADATA"
  }
  depends_on = [
    module.host_project,
    module.project_services
  ]
}

resource "google_compute_global_address" "private_ip_address" {
  provider      = google-beta
  project       = var.host_project_id
  name          = "private-ip-address"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 20
  network       = local.network_id
}

resource "google_service_networking_connection" "private_vpc_connection" {
  provider                = google-beta
  network                 = local.network_id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_address.name]
}

//set up shared project
# A host project provides network resources to associated service projects.
resource "google_compute_shared_vpc_host_project" "host" {
  project = var.host_project_id
}

//set up service project
# A service project gains access to network resources provided by its
# associated host project.
resource "google_compute_shared_vpc_service_project" "service1" {
  host_project    = google_compute_shared_vpc_host_project.host.project
  service_project = var.service_project_id
  depends_on      = [google_compute_shared_vpc_host_project.host]
}
