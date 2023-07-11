// Enable the service in host project
module "project_services" {
  source                      = "terraform-google-modules/project-factory/google//modules/project_services"
  project_id                  = var.project_id
  enable_apis                 = var.enable_apis
  disable_services_on_destroy = var.disable_services_on_destroy
  activate_apis               = var.activate_apis
}
