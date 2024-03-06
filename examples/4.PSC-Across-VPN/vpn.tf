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

module "consumer_project_vpn" {
  source     = "../../modules/net-vpn-ha"
  project_id = var.consumer_project_id
  region     = var.region
  network    = var.consumer_network_name
  name       = var.consumer_gateway_name
  peer_gateways = {
    default = {
      gcp = module.user_project_vpn.self_link
    }
  }
  router_config = {
    asn            = 64514
    advertise_mode = "DEFAULT"
  }
  tunnels = {
    remote-0 = {
      bgp_peer = {
        address = "169.254.1.1"
        asn     = 64513
      }
      bgp_peer_options                = null
      bgp_session_range               = "169.254.1.2/30"
      ike_version                     = 2
      vpn_gateway_interface           = 0
      peer_external_gateway_interface = null
      shared_secret                   = module.consumer_project_vpn.random_secret
    }

    remote-1 = {
      bgp_peer = {
        address = "169.254.2.1"
        asn     = 64513
      }
      bgp_peer_options                = null
      bgp_session_range               = "169.254.2.2/30"
      ike_version                     = 2
      vpn_gateway_interface           = 1
      peer_external_gateway_interface = null
      shared_secret                   = module.consumer_project_vpn.random_secret
    }

  }
  depends_on = [
    module.user_vpc,
    module.consumer_vpc
  ]
}

module "user_project_vpn" {
  source     = "../../modules/net-vpn-ha"
  project_id = var.user_project_id
  region     = var.user_region
  network    = var.user_network_name
  name       = var.user_gateway_name
  router_config = {
    asn = 64513
  }
  peer_gateways = {
    default = {
      gcp = module.consumer_project_vpn.self_link
    }
  }
  tunnels = {
    remote-0 = {
      bgp_peer = {
        address = "169.254.1.2"
        asn     = 64514
      }
      bgp_session_range     = "169.254.1.1/30"
      ike_version           = 2
      vpn_gateway_interface = 0
      shared_secret         = module.consumer_project_vpn.random_secret
    }

    remote-1 = {
      bgp_peer = {
        address = "169.254.2.2"
        asn     = 64514
      }
      bgp_session_range     = "169.254.2.1/30"
      ike_version           = 2
      vpn_gateway_interface = 1
      shared_secret         = module.consumer_project_vpn.random_secret
    }
  }
  depends_on = [
    module.user_vpc,
    module.consumer_vpc
  ]
}