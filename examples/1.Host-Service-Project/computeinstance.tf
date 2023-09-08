
module "google_compute_instance" {
  source               = "../../modules/computeinstance"
  project_id           = var.service_project_id
  subnetwork_id        = local.subnetwork_id
  vm_service_account   = local.vm_service_account
  region               = var.region
  zone                 = var.zone
  subnetwork_project   = var.service_project_id
  gce_tags             = var.gce_tags
  source_image         = var.source_image
  source_image_project = var.source_image_project
  source_image_family  = var.source_image_family
  deletion_protection  = var.deletion_protection
  startup_script       = data.template_file.mysql_installer.rendered
  metadata = {
    "enable-oslogin" : true
  }
  access_config = var.access_config
  depends_on = [
    module.host-vpc,
    module.nat,
  ]
}

data "template_file" "mysql_installer" {
  template = file("../startupscripts/setupsql.sh")
  vars = {
    host_ip          = lookup(module.sql-db.mysql_cloudsql_instance_details, "private_ip_address", "")
    default_username = "default"
    default_password = lookup(module.sql-db.mysql_cloudsql_instance_details, "generated_user_password", "")
    database_name    = var.test_dbname
  }
}
