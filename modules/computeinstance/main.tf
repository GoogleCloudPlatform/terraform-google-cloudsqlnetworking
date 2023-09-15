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


resource "google_compute_instance" "compute_instance" {
  name                = var.compute_instance_name
  project             = var.project_id
  machine_type        = var.machine_type
  zone                = var.zone
  tags                = var.gce_tags
  deletion_protection = var.deletion_protection
  boot_disk {
    initialize_params {
      image = "${var.source_image_project}/${var.source_image_family}"
    }
  }
  network_interface {
    subnetwork = var.subnetwork_id
    dynamic "access_config" {
      for_each = var.access_config == null ? [] : [1]
      content {
        network_tier = try(var.access_config.network_tier, var.network_tier)
        nat_ip       = try(var.access_config.nat_ip, "")
      }
    }
  }
  metadata                = var.metadata
  metadata_startup_script = var.startup_script
  service_account {
    email  = var.vm_service_account.email
    scopes = var.vm_service_account.scopes
  }
}
