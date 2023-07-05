variable "host_project_id" {
  type        = string
  description = "Project Id of the Host GCP Project."
}

variable "host_network_id" {
  type        = string
  description = "Complete network Id. This is required when var.create_network is set of false. e.g. : projects/pm-singleproject-20/global/networks/cloudsql-easy"
}

variable "user_project_id" {
  type        = string
  description = "Project Id of the User GCP Project."
}

variable "user_network_id" {
  type        = string
  default     = ""
  description = "Complete network Id. This is required when var.create_network is set of false. e.g. : projects/pm-singleproject-20/global/networks/cloudsql-easy"
}

variable "region" {
  type        = string
  description = "Name of a GCP region."
}

variable "private_ip_address" {
  type        = string
  description = "The IP address or beginning of the address range represented by this resource."
}

variable "private_ip_address_prefix_length" {
  type        = number
  description = " The prefix length of the IP range. If not present, it means the address field is a single IP address. "
}

variable "ha_vpn_gateway1_name" {
  type    = string
  default = "ha-vpn-1"
}

variable "ha_vpn_router1_name" {
  type    = string
  default = "ha-vpn-router1"
}

variable "tunnel1_name" {
  type    = string
  default = "ha-vpn-tunnel1"
}

variable "shared_secret_mesasge1" {
  type    = string
  default = "a secret message"
}

variable "tunnel2_name" {
  type    = string
  default = "ha-vpn-tunnel2"
}

variable "router1_asn" {
  type    = number
  default = 64514
}

variable "router1_interface1_name" {
  type    = string
  default = "router1-interface1"
}

variable "router1_interface2_name" {
  type    = string
  default = "router1-interface2"
}

variable "router1_peer1_name" {
  type    = string
  default = "router1-peer1"
}

variable "router1_peer2_name" {
  type    = string
  default = "router1-peer2"
}

variable "ha_vpn_gateway2_name" {
  type    = string
  default = "ha-vpn-2"
}

variable "ha_vpn_router2_name" {
  type    = string
  default = "ha-vpn-router2"
}

variable "router2_asn" {
  type    = number
  default = 64515
}

variable "tunnel3_name" {
  type    = string
  default = "ha-vpn-tunnel3"
}

variable "tunnel4_name" {
  type    = string
  default = "ha-vpn-tunnel4"
}

variable "shared_secret_mesasge2" {
  type    = string
  default = "a secret message"
}

variable "router2_interface1_name" {
  type    = string
  default = "router2-interface1"
}

variable "router2_interface2_name" {
  type    = string
  default = "router2-interface2"
}

variable "router2_peer1_name" {
  type    = string
  default = "router2-peer1"
}

variable "router2_peer2_name" {
  type    = string
  default = "router2-peer2"
}

variable "advertised_groups" {
  type    = string
  default = "ALL_SUBNETS"
}

variable "advertised_mode" {
  type    = string
  default = "CUSTOM"
}

variable "advertised_route_priority" {
  type    = number
  default = 100
}
