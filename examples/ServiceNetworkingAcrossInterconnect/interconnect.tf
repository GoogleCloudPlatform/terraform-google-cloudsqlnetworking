# Copyright 2023-2024 Google LLC

# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at

#     http://www.apache.org/licenses/LICENSE-2.0

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

resource "google_compute_router" "interconnect-router" {
  name    = var.ic_router_name
  network = local.network_name
  project = var.host_project_id
  region  = var.region
  bgp {
    asn               = var.ic_router_bgp_asn
    advertise_mode    = var.ic_router_advertise_mode
    advertised_groups = var.ic_router_advertise_groups
    dynamic "advertised_ip_ranges" {
      for_each = concat(var.user_specified_ip_range, ["${local.private_ip_address}/${local.private_ip_address_prefix_length}"])
      iterator = item
      content {
        range = item.value
      }
    }
  }
}


module "vlan_attachment_a" {
  source      = "../../modules/net-vlan-attachment"
  network     = local.network_name
  project_id  = var.host_project_id
  region      = var.region
  name        = var.first_va_name
  description = var.first_va_description
  peer_asn    = var.first_va_asn
  router_config = {
    create = var.create_first_vc_router
    name   = var.ic_router_name
  }
  dedicated_interconnect_config = {
    bandwidth = var.first_va_bandwidth
    bgp_range = var.first_va_bgp_range
    //URL of the underlying Interconnect object that this attachment's traffic will traverse through.
    interconnect = "projects/${local.interconnect_project_id}/global/interconnects/${local.first_interconnect_name}"
    vlan_tag     = var.first_vlan_tag
    project      = local.interconnect_project_id
  }
  depends_on = [
    google_compute_router.interconnect-router
  ]
}

module "vlan_attachment_b" {
  source      = "../../modules/net-vlan-attachment"
  network     = local.network_name
  project_id  = var.host_project_id
  region      = var.region
  name        = var.second_va_name
  description = var.second_va_description
  peer_asn    = var.second_va_asn
  router_config = {
    create = var.create_second_vc_router
    name   = var.ic_router_name
  }
  dedicated_interconnect_config = {
    bandwidth = var.second_va_bandwidth
    bgp_range = var.second_va_bgp_range //candidate_subnet
    //URL of the underlying Interconnect object that this attachment's traffic will traverse through.
    interconnect = "projects/${local.interconnect_project_id}/global/interconnects/${local.second_interconnect_name}"
    vlan_tag     = var.second_vlan_tag
    project      = local.interconnect_project_id
  }
  depends_on = [
    google_compute_router.interconnect-router
  ]
}

