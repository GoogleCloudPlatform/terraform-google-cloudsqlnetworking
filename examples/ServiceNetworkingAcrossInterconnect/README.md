## Introduction

Dedicated Interconnect provides direct physical connections between customer on-premise network and Google network. Dedicated Interconnect enables you to transfer large amounts of data between networks which can be more cost-effective than purchasing additional bandwidth over the public internet. You can get a quick overview about dedicated interconnect from this [link](https://cloud.google.com/network-connectivity/docs/interconnect/concepts/dedicated-overview).

This solution guide assumes that the reader is comfortable with dedicated interconnect as described in the product page. This solution template helps the user to create a dedicated connection with redundancy & 99.9% availability between an on-premise environment and a GCP project present in Google cloud.The solution then establishes a Cloud SQL connection using the private IP address of a Cloud SQL instance created inside a GCP project and a VM instance present in the user on-premise environment.

This solution allows you to securely access Cloud SQL instances present in the host project from an onprem environment, without exposing the Cloud SQL instances to the public internet.

Here is a brief overview of the resources being managed by this terraform solution :

1. Creates a VPC Network and subnet in the host project.
2. Creates a VLAN attachment to existing physical connections established at a colocation facility and opts for redundant connection to ensure 99.9% availability. It can be in the same project or in a different project.
3. Creates a cloud router.
4. Configures VLAN ID and BGP IPs capacity.
5. Configures BGP Session.
6. Creates a Cloud SQL instance in the gcp project.
7. Establishes a connection between the Cloud SQL instance from the on-prem VM instance.

### Benefits:

- Securely access Cloud SQL instances using private IP of Cloud SQL Instance.
- Improves performance and reliability by dedicated interconnect.
- Reduces costs by avoiding public IP addresses.

### Use cases:

- Developing and testing applications from onprem that use Cloud SQL.
- Running production applications that use Cloud SQL.
- Connecting to Cloud SQL instances from on-premises networks.

## Architecture

![Service Networking Interconnect Scenario](../images/servicenetworking_across_interconnect.png)

## Pre-requisite


The expectation is to have customer accessibility to a common peering location(aka colocation facility) where google offers dedicated interconnect & customer routers can directly connect to google device via cross connect.

While we have described the set of steps involved below but this [link](https://cloud.google.com/network-connectivity/docs/interconnect/how-to/dedicated/provisioning-overview) describes the process in more detail.

1. User orders dedicated interconnect through cloud console.
2. User receives [LOA-CFA](https://cloud.google.com/network-connectivity/docs/interconnect/concepts/terminology#loa) (letter of authorisation and connecting facility assignment) information which includes working with colocation facilities to establish cross connection with google peering network.
3. After all the necessary connection have been established which involve setting up interconnect link which turns to `green` suggesting it is `PROVISIONED`.
4. At this point interconnect status is available, physical connection between customer onprem to the edge of google's network is set.
5. User should have terraform and gcloud installed in the machine from which they plan to execute this script. Here are the link that describes the [terraform installation](https://developer.hashicorp.com/terraform/tutorials/aws-get-started/install-cli) steps and [gcloud cli installation steps](https://cloud.google.com/sdk/docs/install) .
6.  User should have access to GCP Projects and access to use the physical connections established at colocation facilty which would be mapped to interconnect project. Following are two network design pattern that this solution template supports

  a.  Using One GCP Project : if using single host project topology.

  b.  Using two or more GCP projects : if using host-service project topology for setting up Host & Service projects. Compute XpnPermission described in Step7 will be required if following the Host-Service project topology.

7. User planning to run this script should have following permissions asssigned to them in the respective projects as described below. User can either use webconsole or gcloud cli to assign these permission to the user identity using which these scripts will be executed. User can follow either step `7.a.` or `7.b.` to complete this step.
   - **Interconnect Project**
      - compute.interconnectAttachments.use
      - compute.interconnectAttachments.get
      - compute.interconnectAttachments.list
   - **Host Project**
      - roles/compute.networkAdmin
      - roles/compute.securityAdmin
      - roles/iam.serviceAccountAdmin
      - roles/serviceusage.serviceUsageAdmin
      - roles/resourcemanager.projectIamAdmin
   - **Service Project**
      - roles/cloudsql.admin
      - roles/compute.instanceAdmin
      - roles/iam.serviceAccountAdmin
      - roles/serviceusage.serviceUsageAdmin
      - roles/resourcemanager.projectIamAdmin
   - [Optional] **Compute XpnPermission**
      - User should have `roles/compute.xpnAdmin` permission at a common folder owning the host and service project. Here is a [link](https://cloud.google.com/compute/docs/access/iam#compute.xpnAdmin) describing the same. This is required to associate the service project with the host project.


     a. **Using Webconsole** : User can either use [GCP web console](https://cloud.google.com/iam/docs/grant-role-console) to assign the IAM permission to the user who plans to run this script.

     b. **Using gcloud cli** : User can either use [gcloud cli](https://cloud.google.com/sdk/gcloud/reference/projects/add-iam-policy-binding) to assign IAM permission to the user who plans to run the script.

     User can then use this service account and impersonate this service account while running the terraform code. You can read more about the google cloud service account impersonsation [here](https://cloud.google.com/iam/docs/service-account-overview#impersonation).

     Once you have created this service account you can then update `providers.tf.template` file by updating the `impersonate_service_account` field with the service account you have created with appropriate permission as described above and renaming the `providers.tf.template` to `providers.tf` file.

     Once the providers.tf file is updated. The updated content of `providers.tf` file may look like

      ```
      provider "google" {
        impersonate_service_account = "iac-sa@<GCP-HOST-PROJECT-ID>.iam.gserviceaccount.com"
      }
      provider "google-beta" {
        impersonate_service_account = "iac-sa@<GCP-HOST-PROJECT-ID>.iam.gserviceaccount.com"
      }
      ```

      **Note :** User should have service account admin and project iam admin permissions in the respective GCP projects in order to assign the above mentioned permissions.


## Execution

1. User should have authenticated using gcloud command `gcloud auth application-default login` command in the cli/machine using which user plans to execute the terraform code. This [link](https://cloud.google.com/sdk/gcloud/reference/auth/application-default/login) describes more detail about the `gcloud auth` command mentioned above.
2. User can now change their working directory using `cd` to the example directory `cloudsql-easy-networking/examples/ServiceNetworkingAcrossInterconnect` in order to execute the terraform code.
3. Update the variables in **terraform.tfvars** as per your configuration like host_project_id, service_project_id etc. User can also go through the [Inputs](#inputs) section of this readme that describes the list of input variables that can be updated. Here are two [examples](#examples) of the terraform.tfvars file which can be referred while updating your terraform.tfvars file.
4. Run command `terraform init `. This command initializes the working directory containing terraform configuration files. More description about terraform init can be found at this [link](https://developer.hashicorp.com/terraform/cli/commands/init).
5. Run command `terraform validate` to validate the configuration files present in this directory. More description about terraform validate can be found at this [link](https://developer.hashicorp.com/terraform/cli/commands/validate).
6. Run command `terraform plan`.
This command creates an execution plan, which lets you preview the changes that terraform plans to make in your infrastructure. More details about this command can be found at [link](https://developer.hashicorp.com/terraform/cli/commands/plan). Review the content displayed in the plan stage and if all looks good then move to next step.
7. Run `terraform apply`and type `yes` when asked for confirmation/approval. This command executes the actions proposed in a terraform plan. More details about this command can be found at [link](https://developer.hashicorp.com/terraform/cli/commands/apply).
8. **Deleting resources** : Run `terraform destroy` and type `yes` when asked for confirmation/approval. This command will delete the resources created using the terraform. More details about this command can be found at [here](https://developer.hashicorp.com/terraform/cli/commands/destroy).

## Examples

1. This example leverages the existing network and subnetwork. Network information like network_name, network_id and subnetwork_id should be passed in the terraform.tfvars file. The terraform.tfvars should look like

    ```
      host_project_id        = "<GCP-HOST-PROJECT-ID>"
      service_project_id     = "<GCP-SERVICE-PROJECT-ID>"
      database_version       = "MYSQL_8_0"
      cloudsql_instance_name = "cn-sqlinstance10"
      region                 = "us-west2"
      zone                   = "us-west2-a"
      create_network         = false
      create_subnetwork      = false
      network_name           = "cloudsql-easy"
      subnetwork_name        = "cloudsql-easy-subnet"
      subnetwork_ip_cidr     = "10.2.0.0/16"

      # Variables for Interconnect
      interconnect_project_id  = "<GCP-INTERCONNECT-PROJECT-ID>"
      first_interconnect_name  = "interconnect-1"
      second_interconnect_name = "interconnect-2"
      ic_router_bgp_asn        = 65001

      //first vlan attachment configuration values
      first_va_name          = "vlan-attachment-a"
      first_va_asn           = "65418"
      create_first_vc_router = false
      first_va_bandwidth     = "BPS_1G"
      first_va_bgp_range     = "169.254.61.0/29"
      first_vlan_tag         = 601


      //second vlan attachment configuration values
      second_va_name          = "vlan-attachment-b"
      second_va_asn           = "65418"
      create_second_vc_router = false
      second_va_bandwidth     = "BPS_1G"
      second_va_bgp_range     = "169.254.61.8/29"
      second_vlan_tag         = 601

    ```

2. This example creates new network and subnetwork with the provided cidr range.

    ```
      host_project_id        = "<GCP-HOST-PROJECT-ID>"
      service_project_id     = "<GCP-SERVICE-PROJECT-ID>"
      database_version       = "MYSQL_8_0"
      cloudsql_instance_name = "cn-sqlinstance10"
      region                 = "us-west2"
      zone                   = "us-west2-a"
      create_network         = true
      create_subnetwork      = true
      network_name           = "cloudsql-easy"
      subnetwork_name        = "cloudsql-easy-subnet"
      subnetwork_ip_cidr     = "10.2.0.0/16"

      # Variables for Interconnect
      interconnect_project_id  = "<GCP-INTERCONNECT-PROJECT-ID>"
      first_interconnect_name  = "interconnect-1"
      second_interconnect_name = "interconnect-2"
      ic_router_bgp_asn        = 65001

      //first vlan attachment configuration values
      first_va_name          = "vlan-attachment-a"
      first_va_asn           = "65418"
      create_first_vc_router = false
      first_va_bandwidth     = "BPS_1G"
      first_va_bgp_range     = "169.254.61.0/29"
      first_vlan_tag         = 601


      //second vlan attachment configuration values
      second_va_name          = "vlan-attachment-b"
      second_va_asn           = "65418"
      create_second_vc_router = false
      second_va_bandwidth     = "BPS_1G"
      second_va_bgp_range     = "169.254.61.8/29"
      second_vlan_tag         = 601

    ```

<!-- BEGIN_TF_DOCS -->
## Requirements

No requirements.

## Providers

| Name | Version |
|------|---------|
| <a name="provider_google"></a> [google](#provider\_google) | 4.84.0 |

## Modules

| Name | Source | Version |
|------|--------|---------|
| <a name="module_firewall_rules"></a> [firewall\_rules](#module\_firewall\_rules) | ../../modules/firewall-rules | n/a |
| <a name="module_gce_sa"></a> [gce\_sa](#module\_gce\_sa) | ../../modules/iam-service-account | n/a |
| <a name="module_google_compute_instance"></a> [google\_compute\_instance](#module\_google\_compute\_instance) | ../../modules/computeinstance | n/a |
| <a name="module_host_project"></a> [host\_project](#module\_host\_project) | ../../modules/services | n/a |
| <a name="module_host_vpc"></a> [host\_vpc](#module\_host\_vpc) | ../../modules/net-vpc | n/a |
| <a name="module_nat"></a> [nat](#module\_nat) | ../../modules/net-cloudnat | n/a |
| <a name="module_service_project"></a> [service\_project](#module\_service\_project) | ../../modules/services | n/a |
| <a name="module_sql_db"></a> [sql\_db](#module\_sql\_db) | ../../modules/cloudsql | n/a |
| <a name="module_tf_host_prjct_service_accounts"></a> [tf\_host\_prjct\_service\_accounts](#module\_tf\_host\_prjct\_service\_accounts) | ../../modules/iam-service-account | n/a |
| <a name="module_tf_svc_prjct_service_accounts"></a> [tf\_svc\_prjct\_service\_accounts](#module\_tf\_svc\_prjct\_service\_accounts) | ../../modules/iam-service-account | n/a |
| <a name="module_vlan_attachment_a"></a> [vlan\_attachment\_a](#module\_vlan\_attachment\_a) | ../../modules/net-vlan-attachment | n/a |
| <a name="module_vlan_attachment_b"></a> [vlan\_attachment\_b](#module\_vlan\_attachment\_b) | ../../modules/net-vlan-attachment | n/a |

## Resources

| Name | Type |
|------|------|
| [google_compute_router.interconnect-router](https://registry.terraform.io/providers/hashicorp/google/latest/docs/resources/compute_router) | resource |
| [google_compute_network.host_vpc](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/compute_network) | data source |
| [google_compute_subnetwork.host_vpc_subnetwork](https://registry.terraform.io/providers/hashicorp/google/latest/docs/data-sources/compute_subnetwork) | data source |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:--------:|
| <a name="input_cloudsql_instance_name"></a> [cloudsql\_instance\_name](#input\_cloudsql\_instance\_name) | Name of the cloud sql instance which will be created. | `string` | n/a | yes |
| <a name="input_database_version"></a> [database\_version](#input\_database\_version) | Database version of the mysql in Cloud SQL . | `string` | n/a | yes |
| <a name="input_first_interconnect_name"></a> [first\_interconnect\_name](#input\_first\_interconnect\_name) | Name of the first interconnect object. This will be used to populate the URL of the underlying Interconnect object that this attachment's traffic will traverse through. | `string` | n/a | yes |
| <a name="input_first_va_asn"></a> [first\_va\_asn](#input\_first\_va\_asn) | (Required) Local BGP Autonomous System Number (ASN). Must be an RFC6996 private ASN, either 16-bit or 32-bit. The value will be fixed for this router resource. | `string` | n/a | yes |
| <a name="input_first_va_bgp_range"></a> [first\_va\_bgp\_range](#input\_first\_va\_bgp\_range) | Up to 16 candidate prefixes that can be used to restrict the allocation of cloudRouterIpAddress and customerRouterIpAddress for this attachment. All prefixes must be within link-local address space (169.254.0.0/16) and must be /29 or shorter (/28, /27, etc). | `string` | n/a | yes |
| <a name="input_first_vlan_tag"></a> [first\_vlan\_tag](#input\_first\_vlan\_tag) | The IEEE 802.1Q VLAN tag for this attachment, in the range 2-4094. | `number` | n/a | yes |
| <a name="input_host_project_id"></a> [host\_project\_id](#input\_host\_project\_id) | Project Id of the Host GCP Project. | `string` | n/a | yes |
| <a name="input_ic_router_bgp_asn"></a> [ic\_router\_bgp\_asn](#input\_ic\_router\_bgp\_asn) | Local BGP Autonomous System Number (ASN). Must be an RFC6996 private ASN, either 16-bit or 32-bit. The value will be fixed for this router resource. | `string` | n/a | yes |
| <a name="input_interconnect_project_id"></a> [interconnect\_project\_id](#input\_interconnect\_project\_id) | The ID of the project in which the resource(physical connection at colocation facilitity) belongs. | `string` | n/a | yes |
| <a name="input_network_name"></a> [network\_name](#input\_network\_name) | Name of the VPC network to be created if var.create\_network is marked as true or Name of the already existing network if var.create\_network is false. | `string` | n/a | yes |
| <a name="input_region"></a> [region](#input\_region) | Name of a GCP region. | `string` | n/a | yes |
| <a name="input_second_interconnect_name"></a> [second\_interconnect\_name](#input\_second\_interconnect\_name) | Name of the second interconnect object. This will be used to populate the URL of the underlying Interconnect object that this attachment's traffic will traverse through. | `string` | n/a | yes |
| <a name="input_second_va_asn"></a> [second\_va\_asn](#input\_second\_va\_asn) | (Required) Local BGP Autonomous System Number (ASN). Must be an RFC6996 private ASN, either 16-bit or 32-bit. The value will be fixed for this router resource. | `string` | n/a | yes |
| <a name="input_second_va_bgp_range"></a> [second\_va\_bgp\_range](#input\_second\_va\_bgp\_range) | Up to 16 candidate prefixes that can be used to restrict the allocation of cloudRouterIpAddress and customerRouterIpAddress for this attachment. All prefixes must be within link-local address space (169.254.0.0/16) and must be /29 or shorter (/28, /27, etc). | `string` | n/a | yes |
| <a name="input_second_vlan_tag"></a> [second\_vlan\_tag](#input\_second\_vlan\_tag) | The IEEE 802.1Q VLAN tag for this attachment, in the range 2-4094. | `number` | n/a | yes |
| <a name="input_subnetwork_ip_cidr"></a> [subnetwork\_ip\_cidr](#input\_subnetwork\_ip\_cidr) | CIDR range for the subnet to be created if var.create\_subnetwork is set to true. | `string` | n/a | yes |
| <a name="input_subnetwork_name"></a> [subnetwork\_name](#input\_subnetwork\_name) | Name of the sub network to be created if var.create\_subnetwork is marked as true or Name of the already existing sub network if var.create\_subnetwork is false. | `string` | n/a | yes |
| <a name="input_zone"></a> [zone](#input\_zone) | Name of a GCP zone, should be in the same region as specified in the region variable. | `string` | n/a | yes |
| <a name="input_access_config"></a> [access\_config](#input\_access\_config) | Access configurations, i.e. IPs via which the VM instance can be accessed via the Internet. | <pre>object({<br>    nat_ip       = string<br>    network_tier = string<br>  })</pre> | `null` | no |
| <a name="input_aggregation_interval"></a> [aggregation\_interval](#input\_aggregation\_interval) | The aggregation interval for flow logs in that subnet. The interval can be set to any of the following: 5-sec (default), 30-sec, 1-min, 5-min, 10-min, or 15-min. | `string` | `"INTERVAL_10_MIN"` | no |
| <a name="input_cloudsql_private_range_cidr"></a> [cloudsql\_private\_range\_cidr](#input\_cloudsql\_private\_range\_cidr) | Cidr of the private IP range. | `string` | `""` | no |
| <a name="input_cloudsql_private_range_name"></a> [cloudsql\_private\_range\_name](#input\_cloudsql\_private\_range\_name) | Name of the default IP range. | `string` | `"privateip-range"` | no |
| <a name="input_cloudsql_private_range_prefix_length"></a> [cloudsql\_private\_range\_prefix\_length](#input\_cloudsql\_private\_range\_prefix\_length) | Prefix length of the private IP range. | `string` | `"20"` | no |
| <a name="input_create_first_vc_router"></a> [create\_first\_vc\_router](#input\_create\_first\_vc\_router) | Select 'true' to create a separate router for this VLAN attachment, or 'false' to use the current router configuration. | `bool` | `false` | no |
| <a name="input_create_mysql_db"></a> [create\_mysql\_db](#input\_create\_mysql\_db) | Flag to check if an mysql db needs to be created | `bool` | `true` | no |
| <a name="input_create_nat"></a> [create\_nat](#input\_create\_nat) | Boolean variable to create the Cloud NAT for allowing the VM to connect to external Internet. | `bool` | `true` | no |
| <a name="input_create_network"></a> [create\_network](#input\_create\_network) | Variable to determine if a new network should be created or not. | `bool` | `true` | no |
| <a name="input_create_second_vc_router"></a> [create\_second\_vc\_router](#input\_create\_second\_vc\_router) | Select 'true' to create a separate router for this VLAN attachment, or 'false' to use the current router configuration. | `bool` | `false` | no |
| <a name="input_create_subnetwork"></a> [create\_subnetwork](#input\_create\_subnetwork) | Variable to determine if a new sub network should be created or not. | `bool` | `true` | no |
| <a name="input_deletion_protection"></a> [deletion\_protection](#input\_deletion\_protection) | Enable delete protection. | `bool` | `true` | no |
| <a name="input_enable_export_routes"></a> [enable\_export\_routes](#input\_enable\_export\_routes) | Enable export routes. | `bool` | `true` | no |
| <a name="input_enable_import_routes"></a> [enable\_import\_routes](#input\_enable\_import\_routes) | Enable import routes. | `bool` | `true` | no |
| <a name="input_enable_private_access"></a> [enable\_private\_access](#input\_enable\_private\_access) | Enable private access in the subnet. | `bool` | `true` | no |
| <a name="input_first_va_bandwidth"></a> [first\_va\_bandwidth](#input\_first\_va\_bandwidth) | Provisioned bandwidth capacity for the first interconnect attachment. | `string` | `"BPS_10G"` | no |
| <a name="input_first_va_description"></a> [first\_va\_description](#input\_first\_va\_description) | The description of the first interconnect attachment | `string` | `"interconnect-a vlan attachment 0"` | no |
| <a name="input_first_va_name"></a> [first\_va\_name](#input\_first\_va\_name) | The name of the first interconnect attachment | `string` | `"vlan-attachment-a"` | no |
| <a name="input_flow_sampling"></a> [flow\_sampling](#input\_flow\_sampling) | Flow sampling rate of logs. | `number` | `0.5` | no |
| <a name="input_gce_sa_name"></a> [gce\_sa\_name](#input\_gce\_sa\_name) | Name of the service account to be used by the compute instance. | `string` | `"gce-sa"` | no |
| <a name="input_gce_tags"></a> [gce\_tags](#input\_gce\_tags) | List of tags to be applied to gce instance. | `list(string)` | <pre>[<br>  "cloudsql"<br>]</pre> | no |
| <a name="input_ic_router_advertise_groups"></a> [ic\_router\_advertise\_groups](#input\_ic\_router\_advertise\_groups) | User-specified list of prefix groups to advertise in custom mode. This field can only be populated if advertiseMode is CUSTOM and is advertised to all peers of the router. | `list(string)` | <pre>[<br>  "ALL_SUBNETS"<br>]</pre> | no |
| <a name="input_ic_router_advertise_mode"></a> [ic\_router\_advertise\_mode](#input\_ic\_router\_advertise\_mode) | User-specified flag to indicate which mode to use for advertisement. Default value is DEFAULT. Possible values are: DEFAULT, CUSTOM | `string` | `"CUSTOM"` | no |
| <a name="input_ic_router_name"></a> [ic\_router\_name](#input\_ic\_router\_name) | Name of the interconnect router. | `string` | `"interconnect-router"` | no |
| <a name="input_metadata"></a> [metadata](#input\_metadata) | The metadata annotations that you want to include in the logs. | `string` | `"INCLUDE_ALL_METADATA"` | no |
| <a name="input_nat_name"></a> [nat\_name](#input\_nat\_name) | Name of the cloud nat to be created. | `string` | `"sqleasy-nat"` | no |
| <a name="input_network_routing_mode"></a> [network\_routing\_mode](#input\_network\_routing\_mode) | Network Routing Mode to be used, Could be REGIONAL or GLOBAL. | `string` | `"GLOBAL"` | no |
| <a name="input_network_tier"></a> [network\_tier](#input\_network\_tier) | Networking tier to be used. | `string` | `"STANDARD"` | no |
| <a name="input_router_name"></a> [router\_name](#input\_router\_name) | Name of the router to be used by the NAT. | `string` | `"sqleasynatrouter"` | no |
| <a name="input_second_va_bandwidth"></a> [second\_va\_bandwidth](#input\_second\_va\_bandwidth) | Provisioned bandwidth capacity for the second interconnect attachment. | `string` | `"BPS_10G"` | no |
| <a name="input_second_va_description"></a> [second\_va\_description](#input\_second\_va\_description) | The description of the second interconnect attachment | `string` | `"interconnect-b vlan attachment 1"` | no |
| <a name="input_second_va_name"></a> [second\_va\_name](#input\_second\_va\_name) | The name of the Second interconnect attachment. | `string` | `"vlan-attachment-b"` | no |
| <a name="input_service_project_id"></a> [service\_project\_id](#input\_service\_project\_id) | Project Id of the Service GCP Project attached to the Host GCP project. | `string` | `null` | no |
| <a name="input_source_image"></a> [source\_image](#input\_source\_image) | Source disk image. If neither source\_image nor source\_image\_family is specified, defaults to the latest public image. | `string` | `""` | no |
| <a name="input_source_image_family"></a> [source\_image\_family](#input\_source\_image\_family) | Source image family. If neither source\_image nor source\_image\_family is specified, defaults to the latest public image. | `string` | `"ubuntu-2204-lts"` | no |
| <a name="input_source_image_project"></a> [source\_image\_project](#input\_source\_image\_project) | Project where the source image comes from. The default project contains images. | `string` | `"ubuntu-os-cloud"` | no |
| <a name="input_startup_script"></a> [startup\_script](#input\_startup\_script) | User startup script to run when instances spin up | `string` | `null` | no |
| <a name="input_terraform_sa_name"></a> [terraform\_sa\_name](#input\_terraform\_sa\_name) | Name of the service account to be used by the terraform for automation & execution in build jobs. | `string` | `"terraform-sa"` | no |
| <a name="input_test_dbname"></a> [test\_dbname](#input\_test\_dbname) | Database Name to be created from startup script. | `string` | `"test_db"` | no |
| <a name="input_user_specified_ip_range"></a> [user\_specified\_ip\_range](#input\_user\_specified\_ip\_range) | User-specified list of individual IP ranges to advertise in custom mode. This range specifies google private api address. | `list(string)` | <pre>[<br>  "199.36.154.8/30"<br>]</pre> | no |

## Outputs

| Name | Description |
|------|-------------|
| <a name="output_cloudsql_generated_password"></a> [cloudsql\_generated\_password](#output\_cloudsql\_generated\_password) | CloudSQL instance generated password. |
| <a name="output_cloudsql_instance_name"></a> [cloudsql\_instance\_name](#output\_cloudsql\_instance\_name) | Name of the cloud sql instance created in the project. |
| <a name="output_cloudsql_private_ip"></a> [cloudsql\_private\_ip](#output\_cloudsql\_private\_ip) | CloudSQL instance private ip address. |
| <a name="output_host_network_id"></a> [host\_network\_id](#output\_host\_network\_id) | Network ID for the host VPC network created in the host project. |
| <a name="output_host_psa_ranges"></a> [host\_psa\_ranges](#output\_host\_psa\_ranges) | PSA range allocated for the private service connection in the host vpc. |
| <a name="output_host_subnetwork_id"></a> [host\_subnetwork\_id](#output\_host\_subnetwork\_id) | Subnetwork ID created inside the host VPC network created in the host project. |
| <a name="output_host_vpc_name"></a> [host\_vpc\_name](#output\_host\_vpc\_name) | Name of the host VPC created in the host project. |
<!-- END_TF_DOCS -->
