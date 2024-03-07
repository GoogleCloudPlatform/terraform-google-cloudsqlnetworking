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
#######################   VPC details for the Consumer Project ####################
###############################################################################

module "consumer_vpc" {
  source     = "../../modules/net-vpc"
  project_id = var.consumer_project_id
  name       = var.consumer_network_name
  vpc_create = var.create_network
  subnets = var.create_subnetwork == true ? [
    {
      name                  = var.consumer_subnetwork_name
      region                = var.region
      ip_cidr_range         = var.consumer_cidr
      enable_private_access = true
      flow_logs_config = {
        flow_sampling        = 0.5
        aggregation_interval = "INTERVAL_10_MIN"
        metadata             = "INCLUDE_ALL_METADATA"
      }
    },
  ] : []
}

data "google_compute_network" "consumer_vpc" {
  count   = var.create_network == false ? 1 : 0
  name    = var.consumer_network_name
  project = var.consumer_project_id
}

data "google_compute_subnetwork" "consumer_vpc_subnetwork" {
  count   = var.create_subnetwork == false ? 1 : 0
  name    = var.consumer_subnetwork_name
  project = var.consumer_project_id
  region  = var.region
}

module "consumer_nat" {
  count                 = var.create_nat ? 1 : 0
  source                = "../../modules/net-cloudnat"
  project_id            = var.consumer_project_id
  region                = var.region
  name                  = var.nat_name
  router_network        = var.consumer_network_name
  router_create         = true
  config_source_subnets = "LIST_OF_SUBNETWORKS"
  router_name           = var.router_name
  config_port_allocation = {
    enable_endpoint_independent_mapping = false
    enable_dynamic_port_allocation      = true
  }
  subnetworks = [{
    self_link            = local.consumer_subnetwork_id,
    config_source_ranges = ["PRIMARY_IP_RANGE"],
    secondary_ranges     = []
  }]
  depends_on = [
    module.consumer_vpc
  ]
}