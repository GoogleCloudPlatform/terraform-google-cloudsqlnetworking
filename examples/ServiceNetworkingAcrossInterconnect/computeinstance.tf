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

###############################################################################
#########################  GCE for the Host Project ###########################
###############################################################################

module "google_compute_instance" {
  source               = "../../modules/computeinstance"
  project_id           = module.host_vpc.project_id
  subnetwork_id        = local.subnetwork_id
  vm_service_account   = local.vm_service_account
  region               = var.region
  zone                 = var.zone
  subnetwork_project   = var.host_project_id
  gce_tags             = var.gce_tags
  source_image         = var.source_image
  source_image_project = var.source_image_project
  source_image_family  = var.source_image_family
  deletion_protection  = var.deletion_protection
  metadata = {
    "enable-oslogin" : true
  }
  access_config = var.access_config
}
