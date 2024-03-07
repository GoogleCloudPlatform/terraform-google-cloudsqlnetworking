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

output "host_vpc_name" {
  value       = local.network_name
  description = "Name of the host VPC created in the host project."
}

output "host_network_id" {
  value       = local.network_id
  description = "Network ID for the host VPC network created in the host project."
}

output "host_subnetwork_id" {
  value       = local.subnetwork_id
  description = "Subnetwork ID created inside the host VPC network created in the host project."
}

output "host_psa_ranges" {
  value       = module.host_vpc.psa_ranges["${var.cloudsql_private_range_name}"]
  description = "PSA range allocated for the private service connection in the host vpc."
}

output "cloudsql_instance_name" {
  value       = local.cloudsql_instance_name[0]
  description = "Name of the cloud sql instance created in the project."
}

output "cloudsql_generated_password" {
  value       = lookup(module.sql_db.cloudsql_instance_details, "generated_user_password", "")
  description = "CloudSQL instance generated password."
  sensitive   = true
}

output "cloudsql_private_ip" {
  value       = lookup(module.sql_db.cloudsql_instance_details, "private_ip_address", "")
  description = "CloudSQL instance private ip address."
}
