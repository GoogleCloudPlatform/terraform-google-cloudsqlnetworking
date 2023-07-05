<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | 3.90.1 |
| <a name="provider_google-beta"></a> [google-beta](#provider\_google-beta) | 4.71.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_cloudsql"></a> [cloudsql](#module\_cloudsql) | GoogleCloudPlatform/sql-db/google//modules/mysql | 8.0.0 |
| <a name="module_compute_instance"></a> [compute\_instance](#module\_compute\_instance) | terraform-google-modules/vm/google//modules/compute_instance | n/a |
| <a name="module_gce_sa"></a> [gce\_sa](#module\_gce\_sa) | terraform-google-modules/service-accounts/google | ~> 3.0 |
| <a name="module_ha-vpn"></a> [ha-vpn](#module\_ha-vpn) | ../../../modules/ha-vpn | n/a |
| <a name="module_host_project"></a> [host\_project](#module\_host\_project) | terraform-google-modules/project-factory/google//modules/project_services | n/a |
| <a name="module_instance_template"></a> [instance\_template](#module\_instance\_template) | terraform-google-modules/vm/google//modules/instance_template | n/a |
| <a name="module_project_services"></a> [project\_services](#module\_project\_services) | terraform-google-modules/project-factory/google//modules/project_services | n/a |
| <a name="module_terraform_service_accounts"></a> [terraform\_service\_accounts](#module\_terraform\_service\_accounts) | terraform-google-modules/service-accounts/google | ~> 3.0 |
| <a name="module_user_compute_instance"></a> [user\_compute\_instance](#module\_user\_compute\_instance) | terraform-google-modules/vm/google//modules/compute_instance | n/a |
| <a name="module_user_gce_sa"></a> [user\_gce\_sa](#module\_user\_gce\_sa) | terraform-google-modules/service-accounts/google | ~> 3.0 |
| <a name="module_user_instance_template"></a> [user\_instance\_template](#module\_user\_instance\_template) | terraform-google-modules/vm/google//modules/instance_template | n/a |
| <a name="module_user_project_services"></a> [user\_project\_services](#module\_user\_project\_services) | terraform-google-modules/project-factory/google//modules/project_services | n/a |

## Resources

| Name | Type |
|------|------|
| [google-beta_google_compute_global_address.private_ip_address](https://registry.terraform.io/providers/hashicorp/google-beta/latest/docs/resources/google_compute_global_address) | resource |
| [google-beta_google_compute_network.private_network](https://registry.terraform.io/providers/hashicorp/google-beta/latest/docs/resources/google_compute_network) | resource |
| [google-beta_google_service_networking_connection.private_vpc_connection](https://registry.terraform.io/providers/hashicorp/google-beta/latest/docs/resources/google_service_networking_connection) | resource |
| [google_compute_firewall.enable_ssh](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_firewall) | resource |
| [google_compute_firewall.rules](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_firewall) | resource |
| [google_compute_firewall.uservpc_enable_ssh](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_firewall) | resource |
| [google_compute_network.uservpc_private_network](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_network) | resource |
| [google_compute_network_peering_routes_config.peering_routes](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_network_peering_routes_config) | resource |
| [google_compute_shared_vpc_host_project.host](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_shared_vpc_host_project) | resource |
| [google_compute_shared_vpc_service_project.service1](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_shared_vpc_service_project) | resource |
| [google_compute_subnetwork.subnet](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_subnetwork) | resource |
| [google_compute_subnetwork.uservpc_subnet](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_subnetwork) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_cloudsql_instance_name"></a> [cloudsql\_instance\_name](#input\_cloudsql\_instance\_name) | Name of the cloud sql instance which will be created. | `string` | n/a | yes |
| <a name="input_create_network"></a> [create\_network](#input\_create\_network) | Variable to determine if a new network should be created or not. | `bool` | `true` | no |
| <a name="input_create_subnetwork"></a> [create\_subnetwork](#input\_create\_subnetwork) | Variable to determine if a new sub network should be created or not. | `bool` | `true` | no |
| <a name="input_create_user_vpc_network"></a> [create\_user\_vpc\_network](#input\_create\_user\_vpc\_network) | Variable to determine if a new network should be created or not. | `bool` | `true` | no |
| <a name="input_create_user_vpc_subnetwork"></a> [create\_user\_vpc\_subnetwork](#input\_create\_user\_vpc\_subnetwork) | Variable to determine if a new sub network should be created or not. | `bool` | `true` | no |
| <a name="input_database_version"></a> [database\_version](#input\_database\_version) | Database version of the mysql in Cloud SQL . | `string` | n/a | yes |
| <a name="input_gce_tags"></a> [gce\_tags](#input\_gce\_tags) | List of tags to be applied to gce instance. | `list(string)` | <pre>[<br>  "cloudsql"<br>]</pre> | no |
| <a name="input_host_project_id"></a> [host\_project\_id](#input\_host\_project\_id) | Project Id of the Host GCP Project. | `string` | n/a | yes |
| <a name="input_network_id"></a> [network\_id](#input\_network\_id) | Complete network Id. This is required when var.create\_network is set of false. e.g. : projects/pm-singleproject-20/global/networks/cloudsql-easy | `string` | `""` | no |
| <a name="input_network_name"></a> [network\_name](#input\_network\_name) | Name of the VPC network to be created if var.create\_network is marked as true. If set as false, then provide the name of the existing network to be used. | `string` | `"cloudsql-easy"` | no |
| <a name="input_network_routing_mode"></a> [network\_routing\_mode](#input\_network\_routing\_mode) | Network Routing Mode to be used, Could be REGIONAL or GLOBAL. | `string` | `"GLOBAL"` | no |
| <a name="input_network_tier"></a> [network\_tier](#input\_network\_tier) | Networking tier to be used. | `string` | `"STANDARD"` | no |
| <a name="input_region"></a> [region](#input\_region) | Name of a GCP region. | `string` | n/a | yes |
| <a name="input_service_project_id"></a> [service\_project\_id](#input\_service\_project\_id) | Project Id of the Service GCP Project attached to the Host GCP project. | `string` | n/a | yes |
| <a name="input_subnetwork_id"></a> [subnetwork\_id](#input\_subnetwork\_id) | Complete subnetwork Id. This is required when var.create\_subnetwork is set of false. e.g. : projects/pm-singleproject-20/regions/us-central1/subnetworks/cloudsql-easy-subnet | `string` | `""` | no |
| <a name="input_subnetwork_ip_cidr"></a> [subnetwork\_ip\_cidr](#input\_subnetwork\_ip\_cidr) | CIDR range for the subnet to be created if var.create\_subnetwork is set to true. | `string` | n/a | yes |
| <a name="input_subnetwork_name"></a> [subnetwork\_name](#input\_subnetwork\_name) | Name of the sub network to be created if var.create\_subnetwork is marked as true. | `string` | `"cloudsql-easy-subnet"` | no |
| <a name="input_target_size"></a> [target\_size](#input\_target\_size) | Number of GCE instances to be created. | `number` | `1` | no |
| <a name="input_user_network_id"></a> [user\_network\_id](#input\_user\_network\_id) | Complete network Id. This is required when var.create\_network is set of false. e.g. : projects/pm-singleproject-20/global/networks/cloudsql-easy | `string` | `""` | no |
| <a name="input_user_network_routing_mode"></a> [user\_network\_routing\_mode](#input\_user\_network\_routing\_mode) | Network Routing Mode to be used, Could be REGIONAL or GLOBAL. | `string` | `"GLOBAL"` | no |
| <a name="input_user_project_id"></a> [user\_project\_id](#input\_user\_project\_id) | Project Id of the User GCP Project. | `string` | n/a | yes |
| <a name="input_user_region"></a> [user\_region](#input\_user\_region) | Name of a GCP region. | `string` | n/a | yes |
| <a name="input_user_subnetwork_id"></a> [user\_subnetwork\_id](#input\_user\_subnetwork\_id) | Complete subnetwork Id. This is required when var.create\_subnetwork is set of false. e.g. : projects/pm-singleproject-20/regions/us-central1/subnetworks/cloudsql-easy-subnet | `string` | `""` | no |
| <a name="input_uservpc_network_name"></a> [uservpc\_network\_name](#input\_uservpc\_network\_name) | Name of the VPC network to be created if var.create\_user\_vpc\_network is marked as true. | `string` | `"cloudsql-user"` | no |
| <a name="input_uservpc_subnetwork_ip_cidr"></a> [uservpc\_subnetwork\_ip\_cidr](#input\_uservpc\_subnetwork\_ip\_cidr) | CIDR range for the subnet to be created if var.create\_subnetwork is set to true. | `string` | n/a | yes |
| <a name="input_uservpc_subnetwork_name"></a> [uservpc\_subnetwork\_name](#input\_uservpc\_subnetwork\_name) | Name of the sub network to be created if var.create\_user\_vpc\_subnetwork is marked as true. | `string` | `"cloudsql-user-subnet"` | no |
| <a name="input_zone"></a> [zone](#input\_zone) | Name of a GCP zone, should be in the same region as specified in the region variable. | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_host_network_id"></a> [host\_network\_id](#output\_host\_network\_id) | Network ID for the host VPC network created in the host project. |
| <a name="output_host_subnetwork_id"></a> [host\_subnetwork\_id](#output\_host\_subnetwork\_id) | Sub Network ID created inside the host VPC network created in the host project. |
| <a name="output_private_ip_address"></a> [private\_ip\_address](#output\_private\_ip\_address) | The IP address or beginning of the address range represented by this resource. |
| <a name="output_private_ip_address_prefix_length"></a> [private\_ip\_address\_prefix\_length](#output\_private\_ip\_address\_prefix\_length) | The prefix length of the IP range. If not present, it means the address field is a single IP address. |
| <a name="output_uservpc_network_id"></a> [uservpc\_network\_id](#output\_uservpc\_network\_id) | Network ID for the User VPC network created in the user project. |
| <a name="output_uservpc_subnetwork_id"></a> [uservpc\_subnetwork\_id](#output\_uservpc\_subnetwork\_id) | Sub Network ID created inside the User VPC network created in the User project. |
<!-- END_TF_DOCS -->