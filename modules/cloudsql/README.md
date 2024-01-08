<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_random"></a> [random](#provider\_random) | n/a |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_mssql"></a> [mssql](#module\_mssql) | GoogleCloudPlatform/sql-db/google//modules/mssql | 17.1.0 |
| <a name="module_mysql"></a> [mysql](#module\_mysql) | GoogleCloudPlatform/sql-db/google//modules/mysql | 17.1.0 |
| <a name="module_postgresql"></a> [postgresql](#module\_postgresql) | GoogleCloudPlatform/sql-db/google//modules/postgresql | 17.1.0 |

## Resources

| Name | Type |
|------|------|
| [random_string.mssql_random_suffix](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/string) | resource |
| [random_string.mysql_random_suffix](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/string) | resource |
| [random_string.postgresql_random_suffix](https://registry.terraform.io/providers/hashicorp/random/latest/docs/resources/string) | resource |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_create_mssql_db"></a> [create\_mssql\_db](#input\_create\_mssql\_db) | Flag to check if an mssql db needs to be created | `bool` | `false` | no |
| <a name="input_create_mysql_db"></a> [create\_mysql\_db](#input\_create\_mysql\_db) | Flag to check if an mysql db needs to be created | `bool` | `true` | no |
| <a name="input_create_postgresql_db"></a> [create\_postgresql\_db](#input\_create\_postgresql\_db) | Flag to check if an postgresql db needs to be created | `bool` | `false` | no |
| <a name="input_database_version"></a> [database\_version](#input\_database\_version) | Database version of the mysql in Cloud SQL . | `string` | n/a | yes |
| <a name="input_deletion_protection"></a> [deletion\_protection](#input\_deletion\_protection) | If resources created, should be protected with deletion\_protection | `bool` | `true` | no |
| <a name="input_ip_configuration"></a> [ip\_configuration](#input\_ip\_configuration) | The ip configuration for the master instances. | <pre>object({<br>    authorized_networks                           = optional(list(map(string)))<br>    ipv4_enabled                                  = bool<br>    private_network                               = optional(string)<br>    require_ssl                                   = optional(bool)<br>    allocated_ip_range                            = optional(string)<br>    enable_private_path_for_google_cloud_services = optional(bool)<br>    psc_enabled                                   = optional(bool)<br>    psc_allowed_consumer_projects                 = optional(list(string))<br>  })</pre> | <pre>{<br>  "allocated_ip_range": null,<br>  "authorized_networks": [],<br>  "enable_private_path_for_google_cloud_services": false,<br>  "ipv4_enabled": true,<br>  "private_network": null,<br>  "psc_allowed_consumer_projects": [],<br>  "psc_enabled": false,<br>  "require_ssl": null<br>}</pre> | no |
| <a name="input_name"></a> [name](#input\_name) | Name of the cloud sql instance which will be created. | `string` | n/a | yes |
| <a name="input_project_id"></a> [project\_id](#input\_project\_id) | Project Id of the GCP Project. | `string` | n/a | yes |
| <a name="input_zone"></a> [zone](#input\_zone) | Name of a GCP zone, should be in the same region as specified in the region variable. | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_cloudsql_instance_psc_attachment"></a> [cloudsql\_instance\_psc\_attachment](#output\_cloudsql\_instance\_psc\_attachment) | The psc\_service\_attachment\_link created for the master instance |
| <a name="output_mssql_cloudsql_instance_name"></a> [mssql\_cloudsql\_instance\_name](#output\_mssql\_cloudsql\_instance\_name) | Name of the cloud sql instance created in the service project. |
| <a name="output_mysql_cloudsql_instance_details"></a> [mysql\_cloudsql\_instance\_details](#output\_mysql\_cloudsql\_instance\_details) | Details of the cloud sql instance created in the service project. |
| <a name="output_mysql_cloudsql_instance_name"></a> [mysql\_cloudsql\_instance\_name](#output\_mysql\_cloudsql\_instance\_name) | Name of the cloud sql instance created in the service project. |
| <a name="output_postgres_cloudsql_instance_name"></a> [postgres\_cloudsql\_instance\_name](#output\_postgres\_cloudsql\_instance\_name) | Name of the cloud sql instance created in the service project. |
<!-- END_TF_DOCS -->