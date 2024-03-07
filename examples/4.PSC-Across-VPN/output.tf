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

output "cloudsql_instance_name" {
  value       = local.cloudsql_instance_name[0]
  description = "Name of the SQL instance created in the producer project."
}

output "compute_instance_name" {
  value       = module.user_project_instance
  description = "Name of the compute instance created"
}

output "reserved_ip" {
  value       = module.compute_address
  description = "IP Address reserved as service endpoint"
}

output "consumer_network_name" {
  value       = module.consumer_vpc.name
  description = "Name of the consumer VPC"
}

output "consumer_network_id" {
  value       = local.consumer_network_id
  description = "Network ID for the consumer VPC network created in the consumer project."
}

output "consumer_subnetwork_id" {
  value       = local.consumer_subnetwork_id
  description = "Sub Network ID created inside the consumer VPC network created in the consumer project."
}

output "cloudsql_instance_psc_attachment" {
  value       = module.sql_db.cloudsql_instance_psc_attachment
  description = "The psc_service_attachment_link created for the master instance"
  sensitive   = false
}