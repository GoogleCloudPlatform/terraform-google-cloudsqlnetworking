// Enable the service in host project
module "host_project" {
  source                      = "terraform-google-modules/project-factory/google//modules/project_services"
  project_id                  = var.host_project_id
  enable_apis                 = true
  disable_services_on_destroy = false
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
  source                      = "terraform-google-modules/project-factory/google//modules/project_services"
  project_id                  = var.service_project_id
  enable_apis                 = true
  disable_services_on_destroy = false
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

// Enable the service in user project
module "user_project_services" {
  source                      = "terraform-google-modules/project-factory/google//modules/project_services"
  project_id                  = var.user_project_id
  enable_apis                 = true
  disable_services_on_destroy = false
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
