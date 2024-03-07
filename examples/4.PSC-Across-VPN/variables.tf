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

variable "create_user_vpc_network" {
  type        = bool
  default     = true
  description = "Variable to determine if a new network should be created or not."
}

variable "create_user_vpc_subnetwork" {
  type        = bool
  default     = true
  description = "Variable to determine if a new sub network should be created or not."
}

variable "create_nat" {
  description = "Boolean variable to create the Cloud NAT for allowing the VM to connect to external Internet."
  type        = bool
  default     = true
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

variable "cloudsql_instance_name" {
  type        = string
  description = "Name of the cloud sql instance which will be created."
  default     = "cloudsql"
}

variable "consumer_project_id" {
  type        = string
  description = "Project ID of the Consumer GCP Project."
}

variable "producer_project_id" {
  type        = string
  description = "Project ID of the Producer GCP Project."
}

variable "user_project_id" {
  type        = string
  description = "Project ID of the User GCP Project."
}

variable "database_version" {
  type        = string
  description = "Database version of the mysql in Cloud SQL ."
  default     = "MYSQL_8_0"
}

variable "region" {
  type        = string
  description = "Name of a GCP region."
}

variable "zone" {
  type        = string
  description = "Name of a GCP zone, should be in the same region as specified in the region variable."
}

variable "user_network_name" {
  type        = string
  description = "Name of the User VPC network to be created if var.create_network is marked as true or Name of the already existing network if var.create_network is false."
}

variable "user_subnetwork_name" {
  type        = string
  description = "Name of the User sub network to be created if var.create_subnetwork is marked as true or Name of the already existing sub network if var.create_subnetwork is false."
}

variable "user_cidr" {
  type        = string
  description = "CIDR range of the user VPC Network"
}

variable "nat_name" {
  type        = string
  description = "Name of the NAT connection for the VM instance to communicate to the internet."
}

variable "router_name" {
  type        = string
  description = "Name of the router for Cloud NAT."
}

variable "consumer_gateway_name" {
  type        = string
  description = "Name of the consumer Gateway."
  default     = "gcp-vpc-gateway1"
}

variable "user_gateway_name" {
  type        = string
  description = "Name of the user Gateway."
  default     = "gcp-vpc-gateway2"
}

variable "user_region" {
  type        = string
  description = "Region for the user cloud entities"
}

variable "user_zone" {
  type        = string
  description = "Zone for the user cloud entities"
}

variable "consumer_cidr" {
  type        = string
  description = "CIDR range of the consumer VPC Network"
}

variable "consumer_network_name" {
  type        = string
  description = "Name of the Consumer VPC network to be created if var.create_network is marked as true or Name of the already existing network if var.create_network is false."
}

variable "consumer_subnetwork_name" {
  type        = string
  description = "Name of the Consumer sub network to be created if var.create_subnetwork is marked as true or Name of the already existing sub network if var.create_subnetwork is false."
}

variable "gce_tags" {
  type        = list(string)
  default     = ["cloudsql"]
  description = "List of tags to be applied to gce instance."
}

variable "endpoint_ip" {
  type        = list(string)
  description = "Endpoint IP address to be reserved for PSC connection"
}

variable "terraform-sa" {
  type        = string
  description = "Service Account to be used by Terraform"
  default     = "terraform-sa"
}

variable "test_dbname" {
  type        = string
  description = "Name of the DB to be created inside the SQL instance"
  default     = "test_db"

}

variable "load_balancing_scheme" {
  description = "Load Balacing Scheme for the ILB/Forwarding Rule"
  type        = string
  default     = ""

}

variable "create_mysql_db" {
  type        = bool
  default     = true
  description = "Flag to check if an mysql db needs to be created"
}

variable "create_postgresql_db" {
  description = "Bool value to create Postgres DB"
  type        = bool
  default     = false
}