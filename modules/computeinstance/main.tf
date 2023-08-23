

module "instance_template" {
  source               = "terraform-google-modules/vm/google//modules/instance_template"
  version              = "8.0.1"
  project_id           = var.project_id
  subnetwork           = var.subnetwork_id
  service_account      = var.vm_service_account
  subnetwork_project   = var.project_id
  tags                 = var.gce_tags
  source_image_project = var.source_image_project
  source_image_family  = var.source_image_family
  metadata             = var.metadata
  startup_script       = var.startup_script
  access_config        = var.access_config
}

module "compute_instance" {
  source              = "terraform-google-modules/vm/google//modules/compute_instance"
  version             = "8.0.1"
  region              = var.region
  zone                = var.zone
  num_instances       = var.target_size
  instance_template   = module.instance_template.self_link
  deletion_protection = var.deletion_protection
}
