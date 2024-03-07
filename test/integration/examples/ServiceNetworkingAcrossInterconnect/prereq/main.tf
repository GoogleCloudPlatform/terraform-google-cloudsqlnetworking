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

/*
1. Create a gce instance and execute 2 commands as a startup script in this instance to route the traffic
2. Create a route to route all traffic from 0.0.0.0/0 to this nat instance
3. Additional route for default-internet-gateway
4. Additional route for private google range
5. Firewall rule to allow for ingress traffic from onprem to vpc
6. Advertise 0.0.0.0/0  to bgp peers
7. Create a cloud nat(to enable internet connectivity to the NAT VM)
*/

locals {
  vm_service_account = {
    email  = module.nat_gce_sa.email
    scopes = ["cloud-platform"]
  }
}

// Routes to route the internet traffic
resource "google_compute_route" "nat-egress" {
  project                = var.host_project_id
  name                   = var.nat_egress_name
  dest_range             = var.egress_destination_range
  network                = var.network_name
  next_hop_instance      = module.google_compute_instance.name
  priority               = 500
  next_hop_instance_zone = var.zone
}

resource "google_compute_route" "internet-egress" {
  project          = var.host_project_id
  name             = var.internet_egress_name
  dest_range       = var.egress_destination_range
  network          = var.network_name
  next_hop_gateway = var.default-internet-gateway
  priority         = 400
  tags             = ["internet-egress"]
}

resource "google_compute_route" "google-api" {
  project          = var.host_project_id
  name             = var.google_api_route_name
  dest_range       = "199.36.153.8/30"
  network          = var.network_name
  next_hop_gateway = var.default-internet-gateway
  priority         = 100
}

// Firewall rules to enable reach of the compute instance to internet

module "firewall_rules" {
  source       = "../../../../../modules/firewall-rules"
  project_id   = var.host_project_id
  network_name = var.network_name
  rules = [
    {
      name          = var.firewall_name
      priority      = 1000
      description   = var.firewall_description
      direction     = "INGRESS"
      source_ranges = ["10.12.160.0/24"] //source range of VM from onprem
      allow = [{
        protocol = "all"
      }]
      log_config = {
        metadata = var.metadata
      }
  }]
}

module "google_compute_instance" {
  source                = "../../../../../modules/computeinstance"
  project_id            = var.host_project_id
  compute_instance_name = var.compute_instance_name
  subnetwork_id         = "projects/${var.host_project_id}/regions/${var.region}/subnetworks/${var.subnetwork_name}"
  vm_service_account    = local.vm_service_account
  region                = var.region
  zone                  = var.zone
  subnetwork_project    = var.host_project_id
  gce_tags              = var.gce_tags
  source_image          = var.source_image
  source_image_project  = var.source_image_project
  source_image_family   = var.source_image_family
  deletion_protection   = var.deletion_protection
  startup_script        = data.template_file.mysql_installer.rendered
  can_ip_forward        = true
  metadata = {
    "enable-oslogin" : true
  }
  access_config = var.access_config
}

data "template_file" "mysql_installer" {
  template = file("../prereq/setupsql-interconnect.sh")
  vars = {
    bastion_ip       = var.bastion_ip
    host_ip          = var.private_sql_ip
    default_username = var.default_username
    default_password = var.default_password
    database_name    = var.test_dbname
  }
}

module "nat_gce_sa" {
  source     = "../../../../../modules/iam-service-account"
  project_id = var.host_project_id
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
