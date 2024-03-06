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

locals {
  network_name           = var.create_network == true ? try(module.host_vpc.name, "") : var.network_name
  network_id             = var.create_network == true ? try(module.host_vpc.id, "") : data.google_compute_network.host_vpc[0].id
  subnetwork_id          = var.create_subnetwork == true ? try(module.host_vpc.subnet_ids["${var.region}/${var.subnetwork_name}"], "") : data.google_compute_subnetwork.host_vpc_subnetwork[0].id
  cloudsql_instance_name = compact(tolist([module.sql_db.mssql_cloudsql_instance_name, module.sql_db.postgres_cloudsql_instance_name, module.sql_db.mysql_cloudsql_instance_name]))
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
  private_ip_address               = module.host_vpc.psa_ranges["${var.cloudsql_private_range_name}"].address
  private_ip_address_prefix_length = module.host_vpc.psa_ranges["${var.cloudsql_private_range_name}"].prefix_length
  vlan_attachment_project_id       = var.host_project_id
  interconnect_project_id          = var.interconnect_project_id
  first_interconnect_name          = var.first_interconnect_name
  second_interconnect_name         = var.second_interconnect_name
  shared_vpc_host                  = var.service_project_id == null ? false : true
}
