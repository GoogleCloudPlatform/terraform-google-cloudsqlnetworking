variable "host_project_id" {
  type        = string
  description = "Project Id of the Host GCP Project."
}

variable "service_project_id" {
  type        = string
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

variable "network_tier" {
  type        = string
  default     = "PREMIUM"
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

variable "test_dbname" {
  type        = string
  default     = "test_db"
  description = "Database Name to be created from startup script."
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

#####################
# Private Range
#####################

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
