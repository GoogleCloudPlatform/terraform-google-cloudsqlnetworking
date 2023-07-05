<!-- BEGIN_TF_DOCS -->
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
| <a name="input_host_network_id"></a> [host\_network\_id](#input\_host\_network\_id) | Complete network Id. This is required when var.create\_network is set of false. e.g. : projects/pm-singleproject-20/global/networks/cloudsql-easy | `string` | n/a | yes |
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
| <a name="input_user_network_id"></a> [user\_network\_id](#input\_user\_network\_id) | Complete network Id. This is required when var.create\_network is set of false. e.g. : projects/pm-singleproject-20/global/networks/cloudsql-easy | `string` | `""` | no |
| <a name="input_user_project_id"></a> [user\_project\_id](#input\_user\_project\_id) | Project Id of the User GCP Project. | `string` | n/a | yes |

## Outputs

No outputs.
<!-- END_TF_DOCS -->