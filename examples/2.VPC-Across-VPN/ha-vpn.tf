
module "host_project_vpn" {
  source     = "../../modules/net-vpn-ha"
  project_id = var.host_project_id
  region     = var.region
  network    = local.network_id
  name       = var.ha_vpn_gateway1_name
  peer_gateways = {
    default = {
        gcp = module.user_project_vpn.self_link
      }
  }
  router_config = {
    asn = var.router2_asn
    custom_advertise = {
      all_subnets = false
      ip_ranges = {
        "${local.private_ip_address}/${local.private_ip_address_prefix_length}" = "privateiprange"
      }
    }
  }
  tunnels = {
    remote-0 = {
      bgp_peer = {
        address = "169.254.1.1"
        asn     = var.router1_asn
      }
      bgp_session_range     = "169.254.1.2/30"
      vpn_gateway_interface = 0
    }
    remote-1 = {
      bgp_peer = {
        address = "169.254.2.1"
        asn     = var.router1_asn
      }
      bgp_session_range     = "169.254.2.2/30"
      vpn_gateway_interface = 1
    }
  }
}

module "user_project_vpn" {
  source        = "../../modules/net-vpn-ha"
  project_id    = var.user_project_id
  region        = var.region
  network       = local.uservpc_network_id
  name          = var.ha_vpn_gateway2_name
  router_config = { asn = var.router1_asn }
  peer_gateways = {
    default = {
        gcp = module.host_project_vpn.self_link
      }
  }
  tunnels = {
    remote-0 = {
      bgp_peer = {
        address = "169.254.1.2"
        asn     = var.router2_asn
      }
      bgp_session_range     = "169.254.1.1/30"
      shared_secret         = module.host_project_vpn.random_secret
      vpn_gateway_interface = 0
    }
    remote-1 = {
      bgp_peer = {
        address = "169.254.2.2"
        asn     = var.router2_asn
      }
      bgp_session_range     = "169.254.2.1/30"
      shared_secret         = module.host_project_vpn.random_secret
      vpn_gateway_interface = 1
    }
  }
}
