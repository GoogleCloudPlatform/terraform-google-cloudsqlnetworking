
// Create service account to be used by terraform
module "terraform_service_accounts" {
  source     = "terraform-google-modules/service-accounts/google"
  version    = "~> 3.0"
  project_id = var.host_project_id
  prefix     = ""
  names      = ["terraform-sa"]
  project_roles = [
    "${var.host_project_id}=>roles/compute.networkAdmin",
    "${var.host_project_id}=>roles/compute.securityAdmin",
    "${var.host_project_id}=>roles/iam.serviceAccountAdmin",
    "${var.host_project_id}=>roles/serviceusage.serviceUsageAdmin",
    "${var.host_project_id}=>roles/resourcemanager.projectIamAdmin",
    "${var.service_project_id}=>roles/cloudsql.admin",
    //add xpnhost permission
    "${var.service_project_id}=>roles/compute.instanceAdmin",
    "${var.service_project_id}=>roles/iam.serviceAccountAdmin",
    "${var.service_project_id}=>roles/serviceusage.serviceUsageAdmin",
    "${var.service_project_id}=>roles/iam.serviceAccountUser",
    "${var.service_project_id}=>roles/resourcemanager.projectIamAdmin",

    "${var.user_project_id}=>roles/compute.instanceAdmin",
    "${var.user_project_id}=>roles/iam.serviceAccountAdmin",
    "${var.user_project_id}=>roles/serviceusage.serviceUsageAdmin",
    "${var.user_project_id}=>roles/iam.serviceAccountUser",
    "${var.user_project_id}=>roles/resourcemanager.projectIamAdmin",
  ]
  depends_on = [
    module.host_project,
    module.project_services
  ]
}

// Create service account to be used by compute instance created inside the Service Project
module "gce_sa" {
  source     = "terraform-google-modules/service-accounts/google"
  version    = "~> 3.0"
  project_id = var.service_project_id
  prefix     = ""
  names      = ["gce-sa"]
  project_roles = [
    "${var.host_project_id}=>roles/compute.networkUser",
    "${var.service_project_id}=>roles/iam.serviceAccountUser",
    "${var.service_project_id}=>roles/cloudsql.client",
  ]
  depends_on = [
    module.host_project,
    module.project_services
  ]
}

// Create service account to be used by compute instance created inside the User Project
module "user_gce_sa" {
  source     = "terraform-google-modules/service-accounts/google"
  version    = "~> 3.0"
  project_id = var.user_project_id
  prefix     = ""
  names      = ["user-gce-sa"]
  project_roles = [
    "${var.user_project_id}=>roles/compute.networkUser",
    "${var.user_project_id}=>roles/iam.serviceAccountUser",
    "${var.service_project_id}=>roles/cloudsql.client",
  ]
  depends_on = [
    module.host_project,
    module.project_services,
    module.module.user_project_services
  ]
}
