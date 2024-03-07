# Copyright 2023-2024 Google LLC
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

module "user_project_instance" {
  source               = "../../modules/computeinstance"
  project_id           = var.user_project_id
  subnetwork_id        = var.user_subnetwork_name
  subnetwork_project   = var.user_project_id
  vm_service_account   = local.vm_service_account
  region               = var.user_region
  zone                 = var.user_zone
  gce_tags             = var.gce_tags
  source_image         = var.source_image
  source_image_project = var.source_image_project
  source_image_family  = var.source_image_family
  deletion_protection  = var.deletion_protection
  startup_script       = data.template_file.mysql_installer.rendered # should be omitted/commented if not creating a MySQL instance
  metadata = {
    "enable-oslogin" : true
  }
  depends_on = [
    module.user_vpc,
    module.user_nat,
    module.user_project_vpn
  ]
}

# should be omitted/commented if not creating a MySQL instance

data "template_file" "mysql_installer" {
  template = file("../startupscripts/setupsql.sh")
  vars = {
    host_ip          = module.compute_address.addresses[0]
    default_username = "default"
    default_password = lookup(module.sql_db.mysql_cloudsql_instance_details, "generated_user_password", "")
    database_name    = var.test_dbname
  }
}