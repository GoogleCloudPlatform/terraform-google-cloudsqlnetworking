variable "project_id" {
  type        = string
  description = "Project Id of the Service GCP Project attached to the Host GCP project."
}

variable "subnetwork_project" {
  type        = string
  description = "The ID of the project in which the subnetwork belongs. If it is not provided."
}

variable "subnetwork_id" {
  type        = string
  default     = ""
  description = "Complete subnetwork Id. This is required when var.create_subnetwork is set of false. e.g. : projects/pm-singleproject-20/regions/us-central1/subnetworks/cloudsql-easy-subnet"
}

variable "machine_type" {
  type        = string
  default     = "e2-small"
  description = "Google compute instance machine type to be created."
}

variable "compute_instance_name" {
  type        = string
  default     = "cloudsql-easy-networking-001"
  description = "Name of the compute instance machine that will be created."
}

variable "vm_service_account" {
  type = object({
    email  = string
    scopes = set(string)
  })
  description = "Service account to attach to the instance. See https://www.terraform.io/docs/providers/google/r/compute_instance_template#service_account."
}
###########
# metadata
###########

variable "startup_script" {
  description = "User startup script to run when instances spin up"
  type        = string
  default     = ""
}

variable "metadata" {
  type        = map(string)
  description = "Metadata, provided as a map"
  default     = {}
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

variable "region" {
  type        = string
  description = "Name of a GCP region."
}

variable "zone" {
  type        = string
  description = "Name of a GCP zone, should be in the same region as specified in the region variable."
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
