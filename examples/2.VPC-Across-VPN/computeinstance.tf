# Copyright 2023 Google LLC
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

###############################################################################
#########################  GCE for the Host Project ###########################
###############################################################################

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
  #startup_script       = data.template_file.mysql_installer.rendered
  metadata = {
    "enable-oslogin" : true
  }
  access_config = var.access_config
  depends_on = [
    module.host-vpc,
  ]
}


###############################################################################
#########################  GCE for the User Project ###########################
###############################################################################

module "user_google_compute_instance" {
  source               = "../../modules/computeinstance"
  project_id           = var.user_project_id
  subnetwork_id        = local.uservpc_subnetwork_id
  vm_service_account   = local.user_vm_service_account
  region               = var.user_region
  zone                 = var.user_zone
  subnetwork_project   = var.user_project_id
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
    module.user-vpc,
    module.user-nat
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
