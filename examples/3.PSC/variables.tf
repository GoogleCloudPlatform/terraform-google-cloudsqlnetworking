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

variable "consumer_project_id" {
  type        = string
  description = "Project ID of the Consumer GCP Project."
}

variable "producer_project_id" {
  type        = string
  description = "Project ID of the Producer GCP Project."
}

variable "consumer_subnetwork_name" {
  type        = string
  description = "Name of the subnetwork to be created in Consumer project"
}

variable "consumer_network_name" {
  type        = string
  description = "Name of the network to be created in Consumer project"
}

variable "region" {
  type        = string
  description = "Name of a GCP region."
}

variable "zone" {
  type        = string
  description = "Name of a GCP zone, should be in the same region as specified in the region variable."
}

variable "gce_tags" {
  type        = list(string)
  default     = ["cloudsql"]
  description = "List of tags to be applied to gce instance."
}

variable "source_image" {
  description = "Source disk image. If neither source_image nor source_image_family is specified, defaults to the latest public image."
  type        = string
  default     = ""
}

variable "source_image_family" {
  description = "Source image family. If neither source_image nor source_image_family is specified, defaults to the latest public image."
  type        = string
  default     = "ubuntu-2204-lts"
}

variable "source_image_project" {
  description = "Project where the source image comes from. The default project contains images."
  type        = string
  default     = "ubuntu-os-cloud"
}

variable "deletion_protection" {
  description = "Enable delete protection."
  type        = bool
  default     = true
}

variable "nat_name" {
  description = "Name of the cloud nat to be created."
  default     = "sqleasy-nat"
}

variable "router_name" {
  description = "Name of the router to be used by the NAT."
  default     = "sqleasynatrouter-1"
}

variable "consumer_cidr" {
  type        = string
  description = "CIDR Range of the Consumer VPC"
}

variable "cloudsql_instance_name" {
  type        = string
  description = "Name of the cloud sql instance which will be created."
}

variable "database_version" {
  type        = string
  description = "Database version of the mysql in Cloud SQL ."
}

variable "reserved_ips" {
  type        = list(string)
  description = "IP to reserve for connecting to the SQL instance."
  default     = ["10.0.0.5"]
}

variable "create_network" {
  type        = bool
  default     = true
  description = "Variable to determine if a new network should be created or not."
}

variable "create_subnetwork" {
  type        = bool
  default     = true
  description = "Variable to determine if a new sub network should be created or not."
}

variable "create_nat" {
  description = "Boolean variable to create the Cloud NAT for allowing the VM to connect to external Internet."
  type        = bool
  default     = true
}

variable "test_dbname" {
  description = "Test DB being created inside the Cloud SQL isntance"
  type        = string
  default     = "test_db"
}

variable "create_mysql_db" {
  description = "Bool value to create MySQL DB"
  type        = bool
  default     = true
}

variable "create_postgresql_db" {
  description = "Bool value to create Postgres DB"
  type        = bool
  default     = false
}

variable "load_balancing_scheme" {
  description = "Load Balacing Scheme for the ILB/Forwarding Rule"
  type        = string
  default     = ""

}
