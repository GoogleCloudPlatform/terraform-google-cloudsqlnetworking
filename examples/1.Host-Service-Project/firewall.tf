# Copyright 2023 Google LLC
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

module "firewall_rules" {
  source       = "../../modules/firewall-rules"
  project_id   = var.host_project_id
  network_name = local.network_name
  rules = [{
    name          = "allow-ssh-ingress"
    description   = "Creates firewall rule targeting tagged instances"
    direction     = "INGRESS"
    source_ranges = ["0.0.0.0/0"]
    target_tags   = var.gce_tags
    allow = [{
      protocol = "tcp"
      ports    = ["22"]
    }]
    log_config = {
      metadata = "INCLUDE_ALL_METADATA"
    }
    },
    {
      name                    = "cloudsql-ingress"
      description             = "Creates firewall rule targeting tagged instances for cloudsql"
      direction               = "INGRESS"
      source_service_accounts = [module.gce_sa.email]
      allow = [{
        protocol = "tcp"
        ports    = ["3306"]
      }]
      log_config = {
        metadata = "INCLUDE_ALL_METADATA"
      }
  }]
}
