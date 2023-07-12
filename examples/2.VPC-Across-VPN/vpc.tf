
###############################################################################
#######################   VPC details for the Host Project ####################
###############################################################################
module "host-vpc" {
  source     = "../../modules/net-vpc"
  project_id = var.host_project_id
  name       = var.network_name
  vpc_create = var.create_network
  subnets = var.create_subnetwork == true ? [
    {
      name                  = var.subnetwork_name
      region                = var.region
      ip_cidr_range         = var.subnetwork_ip_cidr
      enable_private_access = true
      flow_logs_config = {
        flow_sampling        = 0.5
        aggregation_interval = "INTERVAL_10_MIN"
        metadata             = "INCLUDE_ALL_METADATA"
      }
    },
  ] : []
  shared_vpc_host = true
  shared_vpc_service_projects = [
    var.service_project_id
  ]
  psa_config = {
    ranges = {
      "${var.cloudsql_private_range_name}" = "${var.cloudsql_private_range_cidr}/${var.cloudsql_private_range_prefix_length}"
    }
    export_routes = true
    import_routes = true
  }
  depends_on = [
    module.host_project,
    module.project_services
  ]
}

data "google_compute_network" "host_vpc" {
  count   = var.create_network == false ? 1 : 0
  name    = var.network_name
  project = var.host_project_id
}

data "google_compute_subnetwork" "host_vpc_subnetwork" {
  count   = var.create_subnetwork == false ? 1 : 0
  name    = var.subnetwork_name
  project = var.host_project_id
  region  = var.region
}

###############################################################################
#######################   VPC details for the User Project ####################
###############################################################################

module "user-vpc" {
  source     = "../../modules/net-vpc"
  project_id = var.user_project_id
  name       = var.uservpc_network_name
  vpc_create = var.create_user_vpc_network
  subnets = var.create_user_vpc_subnetwork == true ? [
    {
      name                  = var.uservpc_subnetwork_name
      region                = var.user_region
      ip_cidr_range         = var.uservpc_subnetwork_ip_cidr
      enable_private_access = true
      flow_logs_config = {
        flow_sampling        = 0.5
        aggregation_interval = "INTERVAL_10_MIN"
        metadata             = "INCLUDE_ALL_METADATA"
      }
    },
  ] : []
  depends_on = [
    module.host_project,
    module.project_services,
    module.user_project_services
  ]
}

data "google_compute_network" "user_vpc" {
  count   = var.create_user_vpc_network == false ? 1 : 0
  name    = var.uservpc_network_name
  project = var.user_project_id
}

data "google_compute_subnetwork" "user_vpc_subnetwork" {
  count   = var.create_user_vpc_subnetwork == false ? 1 : 0
  name    = var.uservpc_subnetwork_name
  project = var.user_project_id
  region  = var.user_region
}
