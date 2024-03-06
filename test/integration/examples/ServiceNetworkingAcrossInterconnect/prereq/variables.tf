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

variable "host_project_id" {
  type        = string
  description = "Project Id of the Host GCP Project."
}


variable "region" {
  type        = string
  description = "Name of a GCP region."
}

variable "zone" {
  type        = string
  description = "Name of a GCP zone, should be in the same region as specified in the region variable."
}

variable "network_name" {
  type        = string
  description = "Name of the VPC network to be created if var.create_network is marked as true or Name of the already existing network if var.create_network is false."
}

variable "nat_egress_name" {
  type        = string
  default     = "nat-egress"
  description = "Name of the nat egress resource."
}

variable "internet_egress_name" {
  type        = string
  default     = "internet-egress"
  description = "Name of the nat egress resource."
}

variable "subnetwork_name" {
  type        = string
  description = "Name of the sub network to be created if var.create_subnetwork is marked as true or Name of the already existing sub network if var.create_subnetwork is false."
}

variable "test_dbname" {
  type        = string
  default     = "test_db"
  description = "Database Name to be created from startup script."
}

variable "gce_tags" {
  type        = list(string)
  default     = ["internet-egress"]
  description = "List of tags to be applied to gce instance."
}

variable "private_sql_ip" {
  type        = string
  description = "SQL instance private ip address"
}

variable "default_username" {
  type        = string
  default     = "default"
  description = "Username for the SQL instance"
}

variable "default_password" {
  type        = string
  description = "Password for the SQL instance"
}

variable "compute_instance_name" {
  type        = string
  default     = "test-nat"
  description = "Name of the compute instance."
}

variable "gce_sa_name" {
  type        = string
  default     = "nat-gce-sa"
  description = "Name of the service account being used by the compute instance."
}

variable "default-internet-gateway" {
  type        = string
  default     = "default-internet-gateway"
  description = "Name of the default internet gateway."
}

variable "google_api_route_name" {
  type        = string
  default     = "google-api"
  description = "Name of the google api route."
}

variable "egress_destination_range" {
  type        = string
  default     = "0.0.0.0/0"
  description = "Internet egress range."
}

variable "firewall_name" {
  type        = string
  default     = "allow-all-from-onprem-to-vpc-1"
  description = "Name of the firewall."
}

variable "firewall_description" {
  type        = string
  default     = "Creates firewall rule targeting tagged instances for cloudsql"
  description = "Description of the firewall."
}

variable "metadata" {
  type        = string
  default     = "INCLUDE_ALL_METADATA"
  description = "The metadata annotations that you want to include in the logs."
}


#######
# disk
#######

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
  default     = false
}

variable "bastion_ip" {
  description = "IP address of the bastion ip present in the onprem environment."
  type        = string
  default     = "10.12.160.10"
}

###########################
# Public IP
###########################
variable "access_config" {
  description = "Access configurations, i.e. IPs via which the VM instance can be accessed via the Internet."
  type = object({
    nat_ip       = string
    network_tier = string
  })
  default = null
}
