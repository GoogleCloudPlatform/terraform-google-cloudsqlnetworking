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


// Create service account to be used by terraform
module "tf_host_prjct_service_accounts" {
  source     = "../../modules/iam-service-account"
  project_id = var.host_project_id
  name       = var.terraform_sa_name
  # non-authoritative roles granted *to* the service accounts on other resources
  iam_project_roles = {
    "${var.host_project_id}" = [
      "roles/compute.networkAdmin",
      "roles/compute.networkUser",
      "roles/compute.securityAdmin",
      "roles/iam.serviceAccountAdmin",
      "roles/serviceusage.serviceUsageAdmin",
      "roles/resourcemanager.projectIamAdmin",
    ],
  }
  depends_on = [
    module.host_project,
  ]
}

// Assign permision for service project to the service account
module "tf_svc_prjct_service_accounts" {
  count                  = var.service_project_id == null ? 0 : 1
  source                 = "../../modules/iam-service-account"
  project_id             = var.host_project_id
  name                   = var.terraform_sa_name
  service_account_create = false
  # non-authoritative roles granted *to* the service accounts on other resources
  iam_project_roles = {
    "${var.service_project_id}" = [
      "roles/iam.serviceAccountAdmin",
      "roles/compute.instanceAdmin",
      "roles/serviceusage.serviceUsageAdmin",
      "roles/resourcemanager.projectIamAdmin",
      "roles/iam.serviceAccountUser",
      "roles/cloudsql.admin",
    ]
  }
  depends_on = [
    module.service_project,
    module.tf_host_prjct_service_accounts
  ]
}

// Create service account to be used by Compute Instance
module "gce_sa" {
  source     = "../../modules/iam-service-account"
  project_id = module.host_project.project_id
  name       = var.gce_sa_name
  # non-authoritative roles granted *to* the service accounts on other resources
  iam_project_roles = {
    "${var.host_project_id}" = [
      "roles/cloudsql.client",
      "roles/compute.networkUser",
      "roles/iam.serviceAccountUser",
    ],
  }
}
