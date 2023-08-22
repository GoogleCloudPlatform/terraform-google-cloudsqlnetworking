# Copyright 2022 Google LLC. This software is provided as is, without
# warranty or representation for any use or purpose. Your use of it is
# subject to your agreement with Google.

provider "google" {
  impersonate_service_account = "iac-sa-test@pm-singleproject-20.iam.gserviceaccount.com"
}
provider "google-beta" {
  impersonate_service_account = "iac-sa-test@pm-singleproject-20.iam.gserviceaccount.com"
}

terraform {
  backend "gcs" {
    bucket                      = "pm-cncs-cloudsql-easy-networking"
    prefix                      = "test/example1"
    impersonate_service_account = "iac-sa-test@pm-singleproject-20.iam.gserviceaccount.com"
  }
}

