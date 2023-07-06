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
    authorized_networks                           = list(map(string))
    ipv4_enabled                                  = bool
    private_network                               = string
    require_ssl                                   = bool
    allocated_ip_range                            = string
    enable_private_path_for_google_cloud_services = optional(bool)
  })
  default = {
    authorized_networks                           = []
    ipv4_enabled                                  = true
    private_network                               = null
    require_ssl                                   = null
    allocated_ip_range                            = null
    enable_private_path_for_google_cloud_services = false
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

variable "create_postgressql_db" {
  type        = bool
  default     = false
  description = "Flag to check if an postgressql db needs to be created"
}


variable "deletion_protection" {
  type        = bool
  default     = true
  description = "If resources created, should be protected with deletion_protection"
}
