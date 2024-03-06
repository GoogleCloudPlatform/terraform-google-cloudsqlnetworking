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

module "compute_address" {
  source       = "terraform-google-modules/address/google"
  names        = local.compute_address_name
  project_id   = var.consumer_project_id
  region       = var.region
  address_type = "INTERNAL"
  subnetwork   = var.consumer_subnetwork_name
  addresses    = var.endpoint_ip
  depends_on = [
    module.consumer_vpc
  ]
}

resource "google_compute_forwarding_rule" "cloudsql_forwarding_rule" {
  name                  = local.compute_forwarding_rule_name
  project               = var.consumer_project_id
  region                = var.region
  network               = module.consumer_vpc.name
  ip_address            = module.compute_address.self_links[0]
  load_balancing_scheme = var.load_balancing_scheme
  target                = module.sql_db.cloudsql_instance_psc_attachment
  depends_on = [
    module.sql_db,
    module.user_project_instance
  ]
}