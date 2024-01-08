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
  vm_service_account = {
    email  = module.gce_sa.email
    scopes = ["cloud-platform"]
  }
  consumer_network_id          = var.create_network == true ? try(module.consumer_vpc.id, "") : data.google_compute_network.consumer_vpc[0].id
  consumer_subnetwork_id       = var.create_subnetwork == true ? try(module.consumer_vpc.subnet_ids["${var.region}/${var.consumer_subnetwork_name}"], "") : data.google_compute_subnetwork.consumer_vpc_subnetwork[0].id
  compute_address_name         = ["psc-compute-address-${var.cloudsql_instance_name}"]
  compute_forwarding_rule_name = "psc-forwarding-rule-${var.cloudsql_instance_name}"
  ip_configuration = {
    authorized_networks                           = []
    private_network                               = ""
    ipv4_enabled                                  = false
    require_ssl                                   = false
    psc_enabled                                   = true
    allocated_ip_range                            = null
    enable_private_path_for_google_cloud_services = false
    psc_allowed_consumer_projects                 = [var.consumer_project_id]
  }
  cloudsql_instance_name = compact(tolist([module.sql_db.mssql_cloudsql_instance_name, module.sql_db.postgres_cloudsql_instance_name, module.sql_db.mysql_cloudsql_instance_name]))
}