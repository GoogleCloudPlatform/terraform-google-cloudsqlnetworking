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

output "mysql_cloudsql_instance_name" {
  value       = var.create_mysql_db == true ? module.mysql[0].instance_name : null
  description = "Name of the cloud sql instance created in the service project."
}

output "mysql_cloudsql_instance_details" {
  value       = var.create_mysql_db == true ? module.mysql[0] : null
  description = "Details of the cloud sql instance created in the service project."
}

output "mssql_cloudsql_instance_name" {
  value       = var.create_mssql_db == true ? module.mssql[0].instance_name : null
  description = "Name of the cloud sql instance created in the service project."
}

output "postgres_cloudsql_instance_name" {
  value       = var.create_postgresql_db == true ? module.postgresql[0].instance_name : null
  description = "Name of the cloud sql instance created in the service project."
}

output "cloudsql_instance_psc_attachment" {
  value       = var.create_mysql_db == true ? module.mysql[0].instance_psc_attachment : module.postgresql[0].instance_psc_attachment
  description = "The psc_service_attachment_link created for the master instance"
  sensitive   = false
}