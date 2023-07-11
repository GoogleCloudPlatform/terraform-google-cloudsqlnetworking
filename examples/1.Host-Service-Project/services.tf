// Enable the service in host project
module "host_project" {
  source     = "../../modules/services"
  project_id = var.host_project_id
  activate_apis = [
    "iam.googleapis.com",
    "compute.googleapis.com",
    "logging.googleapis.com",
    "sqladmin.googleapis.com",
    "monitoring.googleapis.com",
    "servicenetworking.googleapis.com",
    "cloudresourcemanager.googleapis.com",
  ]
}

// Enable the service in service project
module "project_services" {
  source     = "../../modules/services"
  project_id = var.service_project_id
  activate_apis = [
    "iam.googleapis.com",
    "compute.googleapis.com",
    "logging.googleapis.com",
    "sqladmin.googleapis.com",
    "monitoring.googleapis.com",
    "servicenetworking.googleapis.com",
    "cloudresourcemanager.googleapis.com",
  ]
}

