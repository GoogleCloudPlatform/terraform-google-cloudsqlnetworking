resource "google_compute_ha_vpn_gateway" "ha_gateway1" {
  region  = var.region
  name    = var.ha_vpn_gateway1_name
  project = var.host_project_id
  network = var.host_network_id
}

resource "google_compute_ha_vpn_gateway" "ha_gateway2" {
  region  = var.region
  name    = var.ha_vpn_gateway2_name
  project = var.user_project_id
  network = var.user_network_id
}

resource "google_compute_router" "router1" {
  name    = var.ha_vpn_router1_name
  region  = var.region
  project = var.host_project_id
  network = var.host_network_id
  bgp {
    asn               = var.router1_asn
    advertise_mode    = var.advertised_mode
    advertised_groups = [var.advertised_groups]
    advertised_ip_ranges {
      range = "${var.private_ip_address}/${var.private_ip_address_prefix_length}"
    }
  }
}

resource "google_compute_router" "router2" {
  name    = var.ha_vpn_router2_name
  region  = var.region
  project = var.user_project_id
  network = var.user_network_id
  bgp {
    asn               = var.router2_asn
    advertise_mode    = var.advertised_mode
    advertised_groups = [var.advertised_groups]
  }
}

resource "google_compute_vpn_tunnel" "tunnel1" {
  name                  = var.tunnel1_name
  region                = var.region
  project               = var.host_project_id
  vpn_gateway           = google_compute_ha_vpn_gateway.ha_gateway1.id
  peer_gcp_gateway      = google_compute_ha_vpn_gateway.ha_gateway2.id
  shared_secret         = var.shared_secret_mesasge1
  router                = google_compute_router.router1.id
  vpn_gateway_interface = 0
}

resource "google_compute_vpn_tunnel" "tunnel2" {
  name                  = var.tunnel2_name
  region                = var.region
  project               = var.host_project_id
  vpn_gateway           = google_compute_ha_vpn_gateway.ha_gateway1.id
  peer_gcp_gateway      = google_compute_ha_vpn_gateway.ha_gateway2.id
  shared_secret         = var.shared_secret_mesasge1
  router                = google_compute_router.router1.id
  vpn_gateway_interface = 1
}

resource "google_compute_vpn_tunnel" "tunnel3" {
  name                  = var.tunnel3_name
  region                = var.region
  project               = var.user_project_id
  vpn_gateway           = google_compute_ha_vpn_gateway.ha_gateway2.id
  peer_gcp_gateway      = google_compute_ha_vpn_gateway.ha_gateway1.id
  shared_secret         = var.shared_secret_mesasge2
  router                = google_compute_router.router2.id
  vpn_gateway_interface = 0
}

resource "google_compute_vpn_tunnel" "tunnel4" {
  name                  = var.tunnel4_name
  region                = var.region
  project               = var.user_project_id
  vpn_gateway           = google_compute_ha_vpn_gateway.ha_gateway2.id
  peer_gcp_gateway      = google_compute_ha_vpn_gateway.ha_gateway1.id
  shared_secret         = var.shared_secret_mesasge2
  router                = google_compute_router.router2.id
  vpn_gateway_interface = 1
}

resource "google_compute_router_interface" "router1_interface1" {
  name       = var.router1_interface1_name
  router     = google_compute_router.router1.name
  project    = var.host_project_id
  region     = var.region
  ip_range   = "169.254.0.1/30"
  vpn_tunnel = google_compute_vpn_tunnel.tunnel1.name
}

resource "google_compute_router_peer" "router1_peer1" {
  name                      = var.router1_peer1_name
  router                    = google_compute_router.router1.name
  region                    = var.region
  project                   = var.host_project_id
  peer_ip_address           = "169.254.0.2"
  peer_asn                  = var.router2_asn
  advertised_route_priority = var.advertised_route_priority
  interface                 = google_compute_router_interface.router1_interface1.name
}

resource "google_compute_router_interface" "router1_interface2" {
  name       = var.router1_interface2_name
  router     = google_compute_router.router1.name
  region     = var.region
  project    = var.host_project_id
  ip_range   = "169.254.1.2/30"
  vpn_tunnel = google_compute_vpn_tunnel.tunnel2.name
}

resource "google_compute_router_peer" "router1_peer2" {
  name                      = var.router1_peer2_name
  router                    = google_compute_router.router1.name
  region                    = var.region
  project                   = var.host_project_id
  peer_ip_address           = "169.254.1.1"
  peer_asn                  = var.router2_asn
  advertised_route_priority = var.advertised_route_priority
  interface                 = google_compute_router_interface.router1_interface2.name
}

resource "google_compute_router_interface" "router2_interface1" {
  name       = var.router2_interface1_name
  router     = google_compute_router.router2.name
  region     = var.region
  project    = var.user_project_id
  ip_range   = "169.254.0.2/30"
  vpn_tunnel = google_compute_vpn_tunnel.tunnel3.name
}

resource "google_compute_router_peer" "router2_peer1" {
  name                      = var.router2_peer1_name
  router                    = google_compute_router.router2.name
  region                    = var.region
  project                   = var.user_project_id
  peer_ip_address           = "169.254.0.1"
  peer_asn                  = var.router1_asn
  advertised_route_priority = var.advertised_route_priority
  interface                 = google_compute_router_interface.router2_interface1.name
}

resource "google_compute_router_interface" "router2_interface2" {
  name       = var.router2_interface2_name
  router     = google_compute_router.router2.name
  region     = var.region
  project    = var.user_project_id
  ip_range   = "169.254.1.1/30"
  vpn_tunnel = google_compute_vpn_tunnel.tunnel4.name
}

resource "google_compute_router_peer" "router2_peer2" {
  name                      = var.router2_peer2_name
  router                    = google_compute_router.router2.name
  region                    = var.region
  project                   = var.user_project_id
  peer_ip_address           = "169.254.1.2"
  peer_asn                  = var.router1_asn
  advertised_route_priority = var.advertised_route_priority
  interface                 = google_compute_router_interface.router2_interface2.name
}
