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

variable "service_project_id" {
  type        = string
  default     = null
  description = "Project Id of the Service GCP Project attached to the Host GCP project."
}

variable "cloudsql_instance_name" {
  type        = string
  description = "Name of the cloud sql instance which will be created."
}

variable "region" {
  type        = string
  description = "Name of a GCP region."
}

variable "zone" {
  type        = string
  description = "Name of a GCP zone, should be in the same region as specified in the region variable."
}

variable "database_version" {
  type        = string
  description = "Database version of the mysql in Cloud SQL ."
}

variable "subnetwork_ip_cidr" {
  type        = string
  description = "CIDR range for the subnet to be created if var.create_subnetwork is set to true."

}

variable "gce_tags" {
  type        = list(string)
  default     = ["cloudsql"]
  description = "List of tags to be applied to gce instance."
}

variable "terraform_sa_name" {
  type        = string
  default     = "terraform-sa"
  description = "Name of the service account to be used by the terraform for automation & execution in build jobs."
}

variable "gce_sa_name" {
  type        = string
  default     = "gce-sa"
  description = "Name of the service account to be used by the compute instance."
}

variable "startup_script" {
  description = "User startup script to run when instances spin up"
  type        = string
  default     = null
}

variable "network_tier" {
  type        = string
  default     = "STANDARD"
  description = "Networking tier to be used."
}

variable "network_routing_mode" {
  type        = string
  default     = "GLOBAL"
  description = "Network Routing Mode to be used, Could be REGIONAL or GLOBAL."
}

variable "network_name" {
  type        = string
  description = "Name of the VPC network to be created if var.create_network is marked as true or Name of the already existing network if var.create_network is false."
}

variable "subnetwork_name" {
  type        = string
  description = "Name of the sub network to be created if var.create_subnetwork is marked as true or Name of the already existing sub network if var.create_subnetwork is false."
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

variable "flow_sampling" {
  type        = number
  default     = 0.5
  description = "Flow sampling rate of logs."
}

variable "aggregation_interval" {
  type        = string
  default     = "INTERVAL_10_MIN"
  description = "The aggregation interval for flow logs in that subnet. The interval can be set to any of the following: 5-sec (default), 30-sec, 1-min, 5-min, 10-min, or 15-min."
}

variable "metadata" {
  type        = string
  default     = "INCLUDE_ALL_METADATA"
  description = "The metadata annotations that you want to include in the logs."
}

variable "enable_private_access" {
  type        = bool
  default     = true
  description = "Enable private access in the subnet."
}

variable "enable_export_routes" {
  type        = bool
  default     = true
  description = "Enable export routes."
}

variable "enable_import_routes" {
  type        = bool
  default     = true
  description = "Enable import routes."
}

variable "test_dbname" {
  type        = string
  default     = "test_db"
  description = "Database Name to be created from startup script."
}

#######
# nat
######

variable "router_name" {
  description = "Name of the router to be used by the NAT."
  default     = "sqleasynatrouter"
}

variable "nat_name" {
  description = "Name of the cloud nat to be created."
  default     = "sqleasy-nat"
}

variable "create_nat" {
  description = "Boolean variable to create the Cloud NAT for allowing the VM to connect to external Internet."
  type        = bool
  default     = true
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
  default     = true
}

variable "cloudsql_private_range_name" {
  description = "Name of the default IP range."
  default     = "privateip-range"
}

variable "cloudsql_private_range_cidr" {
  description = "Cidr of the private IP range."
  default     = ""
}

variable "cloudsql_private_range_prefix_length" {
  description = "Prefix length of the private IP range."
  default     = "20"
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

###########################
# Interconnect IP
###########################

variable "interconnect_project_id" {
  description = "The ID of the project in which the resource(physical connection at colocation facilitity) belongs."
  type        = string
}

variable "first_interconnect_name" {
  description = "Name of the first interconnect object. This will be used to populate the URL of the underlying Interconnect object that this attachment's traffic will traverse through."
  type        = string
}

variable "second_interconnect_name" {
  description = "Name of the second interconnect object. This will be used to populate the URL of the underlying Interconnect object that this attachment's traffic will traverse through."
  type        = string
}

variable "ic_router_name" {
  description = "Name of the interconnect router."
  type        = string
  default     = "interconnect-router"
}

variable "ic_router_bgp_asn" {
  description = "Local BGP Autonomous System Number (ASN). Must be an RFC6996 private ASN, either 16-bit or 32-bit. The value will be fixed for this router resource."
  type        = string
}

variable "ic_router_advertise_mode" {
  description = "User-specified flag to indicate which mode to use for advertisement. Default value is DEFAULT. Possible values are: DEFAULT, CUSTOM"
  type        = string
  default     = "CUSTOM"
}

variable "ic_router_advertise_groups" {
  description = "User-specified list of prefix groups to advertise in custom mode. This field can only be populated if advertiseMode is CUSTOM and is advertised to all peers of the router."
  type        = list(string)
  default     = ["ALL_SUBNETS"]
}

variable "user_specified_ip_range" {
  description = "User-specified list of individual IP ranges to advertise in custom mode. This range specifies google private api address."
  type        = list(string)
  default     = ["199.36.154.8/30"]
}

variable "create_mysql_db" {
  type        = bool
  default     = true
  description = "Flag to check if an mysql db needs to be created"
}

## First Vlan Attachment

variable "first_va_name" {
  description = "The name of the first interconnect attachment"
  type        = string
  default     = "vlan-attachment-a"
}

variable "first_va_description" {
  description = "The description of the first interconnect attachment"
  type        = string
  default     = "interconnect-a vlan attachment 0"
}

variable "first_va_asn" {
  description = "(Required) Local BGP Autonomous System Number (ASN). Must be an RFC6996 private ASN, either 16-bit or 32-bit. The value will be fixed for this router resource."
  type        = string
}

variable "first_va_bandwidth" {
  description = "Provisioned bandwidth capacity for the first interconnect attachment."
  type        = string
  default     = "BPS_10G"
}

variable "first_va_bgp_range" {
  description = "Up to 16 candidate prefixes that can be used to restrict the allocation of cloudRouterIpAddress and customerRouterIpAddress for this attachment. All prefixes must be within link-local address space (169.254.0.0/16) and must be /29 or shorter (/28, /27, etc)."
  type        = string
}

variable "first_vlan_tag" {
  description = "The IEEE 802.1Q VLAN tag for this attachment, in the range 2-4094."
  type        = number
}

variable "create_first_vc_router" {
  description = "Select 'true' to create a separate router for this VLAN attachment, or 'false' to use the current router configuration."
  type        = bool
  default     = false
}

## Second Vlan Attachment

variable "second_va_name" {
  description = "The name of the Second interconnect attachment."
  type        = string
  default     = "vlan-attachment-b"
}

variable "second_va_description" {
  description = "The description of the second interconnect attachment"
  type        = string
  default     = "interconnect-b vlan attachment 1"
}

variable "second_va_asn" {
  description = "(Required) Local BGP Autonomous System Number (ASN). Must be an RFC6996 private ASN, either 16-bit or 32-bit. The value will be fixed for this router resource."
  type        = string
}

variable "second_va_bandwidth" {
  description = "Provisioned bandwidth capacity for the second interconnect attachment."
  type        = string
  default     = "BPS_10G"
}

variable "second_va_bgp_range" {
  description = "Up to 16 candidate prefixes that can be used to restrict the allocation of cloudRouterIpAddress and customerRouterIpAddress for this attachment. All prefixes must be within link-local address space (169.254.0.0/16) and must be /29 or shorter (/28, /27, etc)."
  type        = string
}

variable "second_vlan_tag" {
  description = "The IEEE 802.1Q VLAN tag for this attachment, in the range 2-4094."
  type        = number
}

variable "create_second_vc_router" {
  description = "Select 'true' to create a separate router for this VLAN attachment, or 'false' to use the current router configuration."
  type        = bool
  default     = false
}
