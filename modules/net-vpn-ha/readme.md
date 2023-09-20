# Cloud HA VPN Module

This module makes it easy to deploy either GCP-to-GCP or GCP-to-On-prem [Cloud HA VPN](https://cloud.google.com/network-connectivity/docs/vpn/concepts/overview#ha-vpn).

## Examples

### GCP to GCP

```hcl
module "vpn-1" {
  source     = "./fabric/modules/net-vpn-ha"
  project_id = var.project_id
  region     = "europe-west4"
  network    = var.vpc1.self_link
  name       = "net1-to-net-2"
  peer_gateways = {
    default = { gcp = module.vpn-2.self_link }
  }
  router_config = {
    asn = 64514
    custom_advertise = {
      all_subnets = true
      ip_ranges = {
        "10.0.0.0/8" = "default"
      }
    }
  }
  tunnels = {
    remote-0 = {
      bgp_peer = {
        address = "169.254.1.1"
        asn     = 64513
      }
      bgp_session_range     = "169.254.1.2/30"
      vpn_gateway_interface = 0
    }
    remote-1 = {
      bgp_peer = {
        address = "169.254.2.1"
        asn     = 64513
      }
      bgp_session_range     = "169.254.2.2/30"
      vpn_gateway_interface = 1
    }
  }
}

module "vpn-2" {
  source        = "./fabric/modules/net-vpn-ha"
  project_id    = var.project_id
  region        = "europe-west4"
  network       = var.vpc2.self_link
  name          = "net2-to-net1"
  router_config = { asn = 64513 }
  peer_gateways = {
    default = { gcp = module.vpn-1.self_link }
  }
  tunnels = {
    remote-0 = {
      bgp_peer = {
        address = "169.254.1.2"
        asn     = 64514
      }
      bgp_session_range     = "169.254.1.1/30"
      shared_secret         = module.vpn-1.random_secret
      vpn_gateway_interface = 0
    }
    remote-1 = {
      bgp_peer = {
        address = "169.254.2.2"
        asn     = 64514
      }
      bgp_session_range     = "169.254.2.1/30"
      shared_secret         = module.vpn-1.random_secret
      vpn_gateway_interface = 1
    }
  }
}
# tftest modules=2 resources=18
```

Note: When using the `for_each` meta-argument you might experience a Cycle Error due to the multiple `net-vpn-ha` modules referencing each other. To fix this you can create the [google_compute_ha_vpn_gateway](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_ha_vpn_gateway) resources separately and reference them in the `net-vpn-ha` module via the `vpn_gateway` and `peer_gcp_gateway` variables.

### GCP to on-prem

```hcl
module "vpn_ha" {
  source     = "./fabric/modules/net-vpn-ha"
  project_id = var.project_id
  region     = var.region
  network    = var.vpc.self_link
  name       = "mynet-to-onprem"
  peer_gateways = {
    default = {
      external = {
        redundancy_type = "SINGLE_IP_INTERNALLY_REDUNDANT"
        interfaces      = ["8.8.8.8"] # on-prem router ip address
      }
    }
  }
  router_config = { asn = 64514 }
  tunnels = {
    remote-0 = {
      bgp_peer = {
        address = "169.254.1.1"
        asn     = 64513
      }
      bgp_session_range               = "169.254.1.2/30"
      peer_external_gateway_interface = 0
      shared_secret                   = "mySecret"
      vpn_gateway_interface           = 0
    }
    remote-1 = {
      bgp_peer = {
        address = "169.254.2.1"
        asn     = 64513
      }
      bgp_session_range               = "169.254.2.2/30"
      peer_external_gateway_interface = 0
      shared_secret                   = "mySecret"
      vpn_gateway_interface           = 1
    }
  }
}
# tftest modules=1 resources=10
```
<!-- BEGIN TFDOC -->

## Variables

| name | description | type | required | default |
|---|---|:---:|:---:|:---:|
| [name](variables.tf#L17) | VPN Gateway name (if an existing VPN Gateway is not used), and prefix used for dependent resources. | <code>string</code> | ✓ |  |
| [network](variables.tf#L22) | VPC used for the gateway and routes. | <code>string</code> | ✓ |  |
| [project_id](variables.tf#L46) | Project where resources will be created. | <code>string</code> | ✓ |  |
| [region](variables.tf#L51) | Region used for resources. | <code>string</code> | ✓ |  |
| [router_config](variables.tf#L56) | Cloud Router configuration for the VPN. If you want to reuse an existing router, set create to false and use name to specify the desired router. | <code title="object&#40;&#123;&#10;  create    &#61; optional&#40;bool, true&#41;&#10;  asn       &#61; number&#10;  name      &#61; optional&#40;string&#41;&#10;  keepalive &#61; optional&#40;number&#41;&#10;  custom_advertise &#61; optional&#40;object&#40;&#123;&#10;    all_subnets &#61; bool&#10;    ip_ranges   &#61; map&#40;string&#41;&#10;  &#125;&#41;&#41;&#10;&#125;&#41;">object&#40;&#123;&#8230;&#125;&#41;</code> | ✓ |  |
| [peer_gateways](variables.tf#L27) | Configuration of the (external or GCP) peer gateway. | <code title="map&#40;object&#40;&#123;&#10;  external &#61; optional&#40;object&#40;&#123;&#10;    redundancy_type &#61; string&#10;    interfaces      &#61; list&#40;string&#41;&#10;  &#125;&#41;&#41;&#10;  gcp &#61; optional&#40;string&#41;&#10;&#125;&#41;&#41;">map&#40;object&#40;&#123;&#8230;&#125;&#41;&#41;</code> |  | <code>&#123;&#125;</code> |
| [tunnels](variables.tf#L71) | VPN tunnel configurations. | <code title="map&#40;object&#40;&#123;&#10;  bgp_peer &#61; object&#40;&#123;&#10;    address        &#61; string&#10;    asn            &#61; number&#10;    route_priority &#61; optional&#40;number, 1000&#41;&#10;    custom_advertise &#61; optional&#40;object&#40;&#123;&#10;      all_subnets          &#61; bool&#10;      all_vpc_subnets      &#61; bool&#10;      all_peer_vpc_subnets &#61; bool&#10;      ip_ranges            &#61; map&#40;string&#41;&#10;    &#125;&#41;&#41;&#10;  &#125;&#41;&#10;  bgp_session_range               &#61; string&#10;  ike_version                     &#61; optional&#40;number, 2&#41;&#10;  peer_external_gateway_interface &#61; optional&#40;number&#41;&#10;  peer_gateway                    &#61; optional&#40;string, &#34;default&#34;&#41;&#10;  router                          &#61; optional&#40;string&#41;&#10;  shared_secret                   &#61; optional&#40;string&#41;&#10;  vpn_gateway_interface           &#61; number&#10;&#125;&#41;&#41;">map&#40;object&#40;&#123;&#8230;&#125;&#41;&#41;</code> |  | <code>&#123;&#125;</code> |
| [vpn_gateway](variables.tf#L99) | HA VPN Gateway Self Link for using an existing HA VPN Gateway. Ignored if `vpn_gateway_create` is set to `true`. | <code>string</code> |  | <code>null</code> |
| [vpn_gateway_create](variables.tf#L105) | Create HA VPN Gateway. | <code>bool</code> |  | <code>true</code> |

## Outputs

| name | description | sensitive |
|---|---|:---:|
| [bgp_peers](outputs.tf#L18) | BGP peer resources. |  |
| [external_gateway](outputs.tf#L25) | External VPN gateway resource. |  |
| [gateway](outputs.tf#L30) | VPN gateway resource (only if auto-created). |  |
| [id](outputs.tf#L35) | Fully qualified VPN gateway id. |  |
| [name](outputs.tf#L42) | VPN gateway name (only if auto-created). . |  |
| [random_secret](outputs.tf#L47) | Generated secret. |  |
| [router](outputs.tf#L52) | Router resource (only if auto-created). |  |
| [router_name](outputs.tf#L57) | Router name. |  |
| [self_link](outputs.tf#L62) | HA VPN gateway self link. |  |
| [tunnel_names](outputs.tf#L67) | VPN tunnel names. |  |
| [tunnel_self_links](outputs.tf#L75) | VPN tunnel self links. |  |
| [tunnels](outputs.tf#L83) | VPN tunnel resources. |  |

<!-- END TFDOC -->

<!-- BEGIN_TF_DOCS
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | n/a |

## Modules

No modules.

## Resources

| Name | Type |
|------|------|
| [google_compute_ha_vpn_gateway.ha_gateway1](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_ha_vpn_gateway) | resource |
| [google_compute_ha_vpn_gateway.ha_gateway2](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_ha_vpn_gateway) | resource |
| [google_compute_router.router1](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router) | resource |
| [google_compute_router.router2](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router) | resource |
| [google_compute_router_interface.router1_interface1](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router_interface) | resource |
| [google_compute_router_interface.router1_interface2](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router_interface) | resource |
| [google_compute_router_interface.router2_interface1](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router_interface) | resource |
| [google_compute_router_interface.router2_interface2](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router_interface) | resource |
| [google_compute_router_peer.router1_peer1](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router_peer) | resource |
| [google_compute_router_peer.router1_peer2](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router_peer) | resource |
| [google_compute_router_peer.router2_peer1](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router_peer) | resource |
| [google_compute_router_peer.router2_peer2](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router_peer) | resource |
| [google_compute_vpn_tunnel.tunnel1](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_vpn_tunnel) | resource |
| [google_compute_vpn_tunnel.tunnel2](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_vpn_tunnel) | resource |
| [google_compute_vpn_tunnel.tunnel3](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_vpn_tunnel) | resource |
| [google_compute_vpn_tunnel.tunnel4](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_vpn_tunnel) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_advertised_groups"></a> [advertised\_groups](#input\_advertised\_groups) | n/a | `string` | `"ALL_SUBNETS"` | no |
| <a name="input_advertised_mode"></a> [advertised\_mode](#input\_advertised\_mode) | n/a | `string` | `"CUSTOM"` | no |
| <a name="input_advertised_route_priority"></a> [advertised\_route\_priority](#input\_advertised\_route\_priority) | n/a | `number` | `100` | no |
| <a name="input_ha_vpn_gateway1_name"></a> [ha\_vpn\_gateway1\_name](#input\_ha\_vpn\_gateway1\_name) | n/a | `string` | `"ha-vpn-1"` | no |
| <a name="input_ha_vpn_gateway2_name"></a> [ha\_vpn\_gateway2\_name](#input\_ha\_vpn\_gateway2\_name) | n/a | `string` | `"ha-vpn-2"` | no |
| <a name="input_ha_vpn_router1_name"></a> [ha\_vpn\_router1\_name](#input\_ha\_vpn\_router1\_name) | n/a | `string` | `"ha-vpn-router1"` | no |
| <a name="input_ha_vpn_router2_name"></a> [ha\_vpn\_router2\_name](#input\_ha\_vpn\_router2\_name) | n/a | `string` | `"ha-vpn-router2"` | no |
| <a name="input_host_network_id"></a> [host\_network\_id](#input\_host\_network\_id) | Complete network Id. This is required when var.create\_network is set of false. e.g. : projects/<GCP-HOST-PROJECT-ID>/global/networks/cloudsql-easy | `string` | n/a | yes |
| <a name="input_host_project_id"></a> [host\_project\_id](#input\_host\_project\_id) | Project Id of the Host GCP Project. | `string` | n/a | yes |
| <a name="input_private_ip_address"></a> [private\_ip\_address](#input\_private\_ip\_address) | The IP address or beginning of the address range represented by this resource. | `string` | n/a | yes |
| <a name="input_private_ip_address_prefix_length"></a> [private\_ip\_address\_prefix\_length](#input\_private\_ip\_address\_prefix\_length) | The prefix length of the IP range. If not present, it means the address field is a single IP address. | `number` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | Name of a GCP region. | `string` | n/a | yes |
| <a name="input_router1_asn"></a> [router1\_asn](#input\_router1\_asn) | n/a | `number` | `64514` | no |
| <a name="input_router1_interface1_name"></a> [router1\_interface1\_name](#input\_router1\_interface1\_name) | n/a | `string` | `"router1-interface1"` | no |
| <a name="input_router1_interface2_name"></a> [router1\_interface2\_name](#input\_router1\_interface2\_name) | n/a | `string` | `"router1-interface2"` | no |
| <a name="input_router1_peer1_name"></a> [router1\_peer1\_name](#input\_router1\_peer1\_name) | n/a | `string` | `"router1-peer1"` | no |
| <a name="input_router1_peer2_name"></a> [router1\_peer2\_name](#input\_router1\_peer2\_name) | n/a | `string` | `"router1-peer2"` | no |
| <a name="input_router2_asn"></a> [router2\_asn](#input\_router2\_asn) | n/a | `number` | `64515` | no |
| <a name="input_router2_interface1_name"></a> [router2\_interface1\_name](#input\_router2\_interface1\_name) | n/a | `string` | `"router2-interface1"` | no |
| <a name="input_router2_interface2_name"></a> [router2\_interface2\_name](#input\_router2\_interface2\_name) | n/a | `string` | `"router2-interface2"` | no |
| <a name="input_router2_peer1_name"></a> [router2\_peer1\_name](#input\_router2\_peer1\_name) | n/a | `string` | `"router2-peer1"` | no |
| <a name="input_router2_peer2_name"></a> [router2\_peer2\_name](#input\_router2\_peer2\_name) | n/a | `string` | `"router2-peer2"` | no |
| <a name="input_shared_secret_mesasge1"></a> [shared\_secret\_mesasge1](#input\_shared\_secret\_mesasge1) | n/a | `string` | `"a secret message"` | no |
| <a name="input_shared_secret_mesasge2"></a> [shared\_secret\_mesasge2](#input\_shared\_secret\_mesasge2) | n/a | `string` | `"a secret message"` | no |
| <a name="input_tunnel1_name"></a> [tunnel1\_name](#input\_tunnel1\_name) | n/a | `string` | `"ha-vpn-tunnel1"` | no |
| <a name="input_tunnel2_name"></a> [tunnel2\_name](#input\_tunnel2\_name) | n/a | `string` | `"ha-vpn-tunnel2"` | no |
| <a name="input_tunnel3_name"></a> [tunnel3\_name](#input\_tunnel3\_name) | n/a | `string` | `"ha-vpn-tunnel3"` | no |
| <a name="input_tunnel4_name"></a> [tunnel4\_name](#input\_tunnel4\_name) | n/a | `string` | `"ha-vpn-tunnel4"` | no |
| <a name="input_user_network_id"></a> [user\_network\_id](#input\_user\_network\_id) | Complete network Id. This is required when var.create\_network is set of false. e.g. : projects/<GCP-HOST-PROJECT-ID>/global/networks/cloudsql-easy | `string` | `""` | no |
| <a name="input_user_project_id"></a> [user\_project\_id](#input\_user\_project\_id) | Project Id of the User GCP Project. | `string` | n/a | yes |

## Outputs

No outputs.-->
<!-- END_TF_DOCS -->
