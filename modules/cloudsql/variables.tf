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

variable "project_id" {
  type        = string
  description = "Project Id of the GCP Project."
}

variable "name" {
  type        = string
  description = "Name of the cloud sql instance which will be created."
}

variable "zone" {
  type        = string
  description = "Name of a GCP zone, should be in the same region as specified in the region variable."
}

variable "database_version" {
  type        = string
  description = "Database version of the mysql in Cloud SQL ."
}


variable "ip_configuration" {
  description = "The ip configuration for the master instances."
  type = object({
    authorized_networks                           = optional(list(map(string)))
    ipv4_enabled                                  = bool
    private_network                               = optional(string)
    require_ssl                                   = optional(bool)
    allocated_ip_range                            = optional(string)
    enable_private_path_for_google_cloud_services = optional(bool)
    psc_enabled                                   = optional(bool)
    psc_allowed_consumer_projects                 = optional(list(string))
  })
  default = {
    authorized_networks                           = []
    ipv4_enabled                                  = true
    private_network                               = null
    require_ssl                                   = null
    allocated_ip_range                            = null
    enable_private_path_for_google_cloud_services = false
    psc_enabled                                   = false
    psc_allowed_consumer_projects                 = []
  }
}

variable "create_mysql_db" {
  type        = bool
  default     = true
  description = "Flag to check if an mysql db needs to be created"
}

variable "create_mssql_db" {
  type        = bool
  default     = false
  description = "Flag to check if an mssql db needs to be created"
}

variable "create_postgresql_db" {
  type        = bool
  default     = false
  description = "Flag to check if an postgresql db needs to be created"
}

variable "deletion_protection" {
  type        = bool
  default     = true
  description = "If resources created, should be protected with deletion_protection"
}
