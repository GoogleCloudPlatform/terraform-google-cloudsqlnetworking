// Create service account to be used by Terraform when integrated with IaC
module "terraform_service_accounts" {
  source     = "../../modules/iam-service-account"
  project_id = var.host_project_id
  name       = "terraform-sa"
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
    module.host_project,
    module.project_services
  ]
}

// Create service account to be used by Compute Instance
module "gce_sa" {
  source     = "../../modules/iam-service-account"
  project_id = var.service_project_id
  name       = "gce-sa"
  # non-authoritative roles granted *to* the service accounts on other resources
  iam_project_roles = {
    "${var.host_project_id}" = [
      "roles/compute.networkUser",
    ],
    "${var.service_project_id}" = [
      "roles/iam.serviceAccountUser",
      "roles/cloudsql.client",
    ]
  }
  depends_on = [
    module.host_project,
    module.project_services
  ]
}
