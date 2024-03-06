# Copyright 2023-2024 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


###############################################################################
#######################   VPC details for the Host Project ####################
###############################################################################
module "host_vpc" {
  source     = "../../modules/net-vpc"
  project_id = var.host_project_id
  name       = var.network_name
  vpc_create = var.create_network
  subnets = var.create_subnetwork == true ? [
    {
      name                  = var.subnetwork_name
      region                = var.region
      ip_cidr_range         = var.subnetwork_ip_cidr
      enable_private_access = var.enable_private_access
      flow_logs_config = {
        flow_sampling        = var.flow_sampling
        aggregation_interval = var.aggregation_interval
        metadata             = var.metadata
      }
    },
  ] : []
  shared_vpc_host = local.shared_vpc_host
  shared_vpc_service_projects = [
    try(module.service_project[0].project_id, "")
  ]
  psa_config = {
    ranges = {
      "${var.cloudsql_private_range_name}" = "${var.cloudsql_private_range_cidr}/${var.cloudsql_private_range_prefix_length}"
    }
    export_routes = var.enable_export_routes
    import_routes = var.enable_import_routes
  }
  depends_on = [
    module.host_project,
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

module "nat" {
  count                 = var.create_nat ? 1 : 0
  source                = "../../modules/net-cloudnat"
  project_id            = var.host_project_id
  region                = var.region
  name                  = var.nat_name
  router_network        = var.network_name
  router_create         = true
  config_source_subnets = "LIST_OF_SUBNETWORKS"
  router_name           = var.router_name
  config_port_allocation = {
    enable_endpoint_independent_mapping = false
    enable_dynamic_port_allocation      = true
  }
  subnetworks = [{
    self_link            = local.subnetwork_id,
    config_source_ranges = ["PRIMARY_IP_RANGE"],
    secondary_ranges     = []
  }]
  depends_on = [
    module.host_vpc
  ]
}
