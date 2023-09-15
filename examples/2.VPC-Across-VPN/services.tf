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

// Enable the service in user project
module "user_project_services" {
  source     = "../../modules/services"
  project_id = var.user_project_id
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
