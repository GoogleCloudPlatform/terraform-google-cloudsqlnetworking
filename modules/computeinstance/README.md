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
| [google_compute_instance.compute_instance](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_instance) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | Project Id of the Service GCP Project attached to the Host GCP project. | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | Name of a GCP region. | `string` | n/a | yes |
| <a name="input_subnetwork_project"></a> [subnetwork\_project](#input\_subnetwork\_project) | The ID of the project in which the subnetwork belongs. If it is not provided. | `string` | n/a | yes |
| <a name="input_vm_service_account"></a> [vm\_service\_account](#input\_vm\_service\_account) | Service account to attach to the instance. See https://www.terraform.io/docs/providers/google/r/compute_instance_template#service_account. | <pre>object({<br>    email  = string<br>    scopes = set(string)<br>  })</pre> | n/a | yes |
| <a name="input_zone"></a> [zone](#input\_zone) | Name of a GCP zone, should be in the same region as specified in the region variable. | `string` | n/a | yes |
| <a name="input_access_config"></a> [access\_config](#input\_access\_config) | Access configurations, i.e. IPs via which the VM instance can be accessed via the Internet. | <pre>list(object({<br>    nat_ip       = string<br>    network_tier = string<br>  }))</pre> | `[]` | no |
| <a name="input_compute_instance_name"></a> [compute\_instance\_name](#input\_compute\_instance\_name) | Name of the compute instance machine that will be created. | `string` | `"cloudsql-easy-networking-001"` | no |
| <a name="input_deletion_protection"></a> [deletion\_protection](#input\_deletion\_protection) | Enable delete protection. | `bool` | `true` | no |
| <a name="input_gce_tags"></a> [gce\_tags](#input\_gce\_tags) | List of tags to be applied to gce instance. | `list(string)` | <pre>[<br>  "cloudsql"<br>]</pre> | no |
| <a name="input_machine_type"></a> [machine\_type](#input\_machine\_type) | Google compute instance machine type to be created. | `string` | `"e2-small"` | no |
| <a name="input_metadata"></a> [metadata](#input\_metadata) | Metadata, provided as a map | `map(string)` | `{}` | no |
| <a name="input_network_tier"></a> [network\_tier](#input\_network\_tier) | Networking tier to be used. | `string` | `"PREMIUM"` | no |
| <a name="input_source_image"></a> [source\_image](#input\_source\_image) | Source disk image. If neither source\_image nor source\_image\_family is specified, defaults to the latest public image. | `string` | `""` | no |
| <a name="input_source_image_family"></a> [source\_image\_family](#input\_source\_image\_family) | Source image family. If neither source\_image nor source\_image\_family is specified, defaults to the latest public image. | `string` | `"ubuntu-2204-lts"` | no |
| <a name="input_source_image_project"></a> [source\_image\_project](#input\_source\_image\_project) | Project where the source image comes from. The default project contains images. | `string` | `"ubuntu-os-cloud"` | no |
| <a name="input_startup_script"></a> [startup\_script](#input\_startup\_script) | User startup script to run when instances spin up | `string` | `""` | no |
| <a name="input_subnetwork_id"></a> [subnetwork\_id](#input\_subnetwork\_id) | Complete subnetwork Id. This is required when var.create\_subnetwork is set of false. e.g. : projects/pm-singleproject-20/regions/us-central1/subnetworks/cloudsql-easy-subnet | `string` | `""` | no |

## Outputs

No outputs.
<!-- END_TF_DOCS -->