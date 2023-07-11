variable "project_id" {
  type        = string
  description = "Project Id of the GCP Project for which services needs to be enabled."
}

variable "enable_apis" {
  description = "Whether to actually enable the APIs. If false, this module is a no-op."
  type        = bool
  default     = true
}

variable "activate_apis" {
  description = "The list of apis to activate within the project"
  type        = list(string)
  default     = []
}

variable "disable_services_on_destroy" {
  description = "Whether project services will be disabled when the resources are destroyed. https://www.terraform.io/docs/providers/google/r/google_project_service.html#disable_on_destroy"
  default     = false
  type        = bool
}

variable "disable_dependent_services" {
  description = "Whether services that are enabled and which depend on this service should also be disabled when this service is destroyed. https://www.terraform.io/docs/providers/google/r/google_project_service.html#disable_dependent_services"
  default     = false
  type        = bool
}
