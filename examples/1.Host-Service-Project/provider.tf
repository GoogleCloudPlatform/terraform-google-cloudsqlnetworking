# Copyright 2022 Google LLC. This software is provided as is, without
# warranty or representation for any use or purpose. Your use of it is
# subject to your agreement with Google.

provider "google" {
  impersonate_service_account = ""
}
provider "google-beta" {
  impersonate_service_account = ""
}

terraform {
  backend "gcs" {
    bucket                      = ""
    prefix                      = ""
    impersonate_service_account = ""
  }
}

