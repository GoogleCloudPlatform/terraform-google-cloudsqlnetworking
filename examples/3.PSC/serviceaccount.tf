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

// Create service account to be used by Terraform when integrated with IaC
module "terraform_service_accounts" {
  source     = "../../modules/iam-service-account"
  project_id = var.consumer_project_id
  name       = "terraform-sa"
  # non-authoritative roles granted *to* the service accounts on other resources
  iam_project_roles = {
    "${var.consumer_project_id}" = [
      "roles/compute.networkAdmin",
      "roles/compute.networkUser",
      "roles/compute.securityAdmin",
      "roles/iam.serviceAccountAdmin",
      "roles/serviceusage.serviceUsageAdmin",
      "roles/resourcemanager.projectIamAdmin",
      "roles/compute.instanceAdmin",
    ],
    "${var.producer_project_id}" = [
      "roles/iam.serviceAccountAdmin",
      "roles/serviceusage.serviceUsageAdmin",
      "roles/resourcemanager.projectIamAdmin",
      "roles/iam.serviceAccountUser",
      "roles/cloudsql.admin",
    ]
  }
  depends_on = [
    module.consumer_project,
    module.producer_project
  ]
}

// Service account to be used by Compute Instance
module "gce_sa" {
  source     = "../../modules/iam-service-account"
  project_id = var.consumer_project_id
  name       = "gce-sa"
  # non-authoritative roles granted *to* the service accounts on other resources
  iam_project_roles = {
    "${var.producer_project_id}" = [
      "roles/compute.networkUser",
    ],
    "${var.consumer_project_id}" = [
      "roles/iam.serviceAccountUser",
      "roles/cloudsql.client",
    ]
  }
  depends_on = [
    module.consumer_project,
    module.producer_project
  ]
}