# CloudSQL Simplified Networking

**Note** - This repository is no longer actively maintained. Please refer to [cloudnetworking-config-solutions](https://github.com/GoogleCloudPlatform/cloudnetworking-config-solutions/tree/main) for the latest version.

Streamline your Google Cloud SQL deployment and management with this comprehensive Terraform repository. This repository's solutions simplify networking configuration for Cloud SQL instances. It bundles Terraform modules to make it easier to create and manage Cloud SQL instances and all its dependencies such as cloud networking resources, IAM policies and associated service accounts.

This repository contains easy-to-use Terraform based solutions that would help a user set up all the required components to consume Cloud SQL with private IPs. This aims to help different users like database administrators and application engineers quickly configure Cloud SQL with cloud networking.

Here are some benefits of using this repository:

- **Simplified configuration:** The Terraform modules in this repository abstract away the complexity of configuring Google Cloud networking for Cloud SQL, making it easier to get started, especially for users without a comprehensive Google Cloud experience.

- **Consistency and repeatability:** Terraform ensures that your Cloud SQL configuration is consistent and repeatable across multiple environments. This helps reduce errors and improve operational efficiency.

- **Flexibility and scalability:** The Terraform modules in this repository are flexible and scalable enough to meet the needs of a wide range of deployments from small development projects to large scale enterprise applications.

If you are looking for a simple and efficient way to configure Cloud SQL with private IP in Google Cloud, this repository is a good place to start.

## Example use-case

A database administrator is responsible for setting up Cloud SQL instance(s) for a production application and is asked to configure the Cloud SQL instance to use a private IP using Service Networking or Private Service Connect.

Using the Terraform modules in this repository, the database administrator can easily configure the Cloud SQL instance and all the necessary network resources with just a few lines of code. This would typically involve the following steps:

1. Creating a new Terraform configuration file and defining the variables for the Cloud SQL instance and network configuration.
2. Initializing Terraform and downloading the necessary providers and modules.
3. Running Terraform plan and apply for the configuration.

Once the Terraform configuration has been applied, your Cloud SQL instance will be up and running with a private IP address to connect to the instance.

## Supported Usage

Many examples are included in the [examples](./examples/) folder which describes the complete end-to-end examples covering different scenarios along with its implementation guide and architecture design.

1. [Host Service Project Scenario](./examples/1.Host-Service-Project) : This solution guides a user through the steps to establish a host and a service project, create a Cloud SQL instance and a VM instance in the service project, and connect the VM instance to the Cloud SQL instance using the VM's private IP address. The host project contains the VPC, subnets, and firewall rules.


2. This section outlines two user journeys for Service Networking: using **HA-VPN** or **Dedicated Interconnect**.

    2.1 [VPC across VPN Tunnel Scenario](./examples/2.VPC-Across-VPN) : This solution guides a user to create a highly available (HA) VPN connection between a user project and a host project with a service project attached. The solution then establishes a Cloud SQL connection using the private IP address of the Cloud SQL instance created in the service project and a VM instance created in the user project.

    2.2. [Service Networking Across Interconnect](./examples/ServiceNetworkingAcrossInterconnect) : This solution helps the user to create a dedicated interconnect connection with redundancy & 99.9% availability between an on-premise environment and a GCP project. The solution establishes a Cloud SQL connection using the private IP address of a Cloud SQL instance created inside a GCP project and a VM instance present in the user on-premise environment.

3. [PSC](./examples/3.PSC) : This solution guides a user to create a PSC enabled Cloud SQL instance with a consumer and producer project setup having a compute VM instance created in the consumer project connecting to the Cloud SQL instance through PSC service endpoint. The consumer project contains the consumer VPC, service endpoint & firewall rules to connect to the SQL instance in the producer project.

4. [PSC across VPN Scenario](./examples/4.PSC-Across-VPN) : This solution guides a user to create a HA VPN connection between user and consumer project to connect to a PSC enabled Cloud SQL instance in a producer project from a compute VM instance through PSC service endpoint. The compute instance in the consumer project connects to the PSC service endpoint via the VPN connection. The PSC service endpoint connects to the Cloud SQL instance.

All [GA offerings](https://cloud.google.com/sql/docs/sqlserver/private-ip) for Cloud SQL (MySQL, PostgreSQL, MSSQL) are supported through our solutions.

## Variables

To control a module's behavior, change variable values. Behavior of each of these variables has been documented in the README files for the respective examples. Some of the common variables have been described below but the respective README file should act as reference point as it is more descriptive:

- `cloudsql_instance_name` - set this variable to a string value as the name of the Cloud SQL instance.
- `database_version` - set this variable to a specific [database version](https://cloud.google.com/sql/docs/mysql/db-versions) for the Cloud SQL instance.
- `host_project_id` - set this variable to the Google Cloud project ID, which will be acting as the host project for the setup. The necessary VPC network and subnetwork will be created here (if they do not exist already).
- `region` - set this variable to the Google Cloud [region](https://cloud.google.com/compute/docs/regions-zones) and this will be used to create resources like Cloud SQL instance in Google Cloud.
- `create_network` - this is a boolean variable. Set this to `true` if you want to create a new network, else set it to `false` if you plan to use an existing network.
- `create_subnetwork` - this is a boolean variable. Set this to `true` if you want to create a new subnetwork, else set to `false` if you plan to use existing network.
- `network_name` - set this to the VPC name; either a new VPC network will be created or an existing network will be used, depending on the value of `create_network`.
- `subnetwork_name` - set this to the subnetwork name; either a new subnetwork will be created or an existing subnetwork will be used, depending on the value of `create_subnetwork`.
- `deletion_protection` - set this flag to `true` if you want to enable deletion protection, which protects from accidental deletions. Otherwise, set it to `false`.
- `access_config` - set this to `null` if you do not want the VM instance to have an exernal IP
- `create_nat` - set this to `true` if you want to have a [Cloud NAT](https://cloud.google.com/nat/docs/overview) configured, else set this to `false`

**Note** - If `access_config = null` and `create_nat = false`, the VM instance will be disabled from downloading any clients. However, connections to other VMs and the Cloud SQL instance using private IP will continue to work.

## State files

Terraform must store state about your managed infrastructure and configuration. It is recommended to use the Google Cloud Storage bucket or an equivalent storage place where the terraform stores its [state file](https://developer.hashicorp.com/terraform/language/state).

### provider.tf file

Each provider.tf file carries information like service account to be impersonated, bucket to be used for storing the terraform state file and prefix that will be used in this bucket for storing the terraform state.

Following are sample commands that can be used to update the existing provider.tf.template file to create a provider.tf file.

- For example : 1.Host-Service-Project

```
export _TF_SERVICE_ACCOUNT="<ENTER THE SERVICE ACCOUNT HERE>"
export _TF_BUCKET_NAME="<ENTER THE GCS BUCKET NAME HERE>"
export _TF_EXAMPLE1_PREFIX="<ENTER THE GCS PREFIX NAME HERE>"

sed \
-e "s|ENTER_TF_SERVICE_ACCOUNT|$_TF_SERVICE_ACCOUNT|" \
-e "s|ENTER_TF_BUCKET_NAME|$_TF_BUCKET_NAME|" \
-e "s|ENTER_TF_EXAMPLE1_PREFIX|$_TF_EXAMPLE1_PREFIX|" \
examples/1.Host-Service-Project/provider.tf.template > examples/1.Host-Service-Project/provider.tf
```

- For example : 2.VPC-Across-VPN

```
export _TF_SERVICE_ACCOUNT="<ENTER THE SERVICE ACCOUNT HERE>"
export _TF_BUCKET_NAME="<ENTER THE GCS BUCKET NAME HERE>"
export _TF_EXAMPLE2_PREFIX="<ENTER THE GCS PREFIX NAME HERE>"

sed \
-e "s|ENTER_TF_SERVICE_ACCOUNT|$_TF_SERVICE_ACCOUNT|" \
-e "s|ENTER_TF_BUCKET_NAME|$_TF_BUCKET_NAME|" \
-e "s|ENTER_TF_EXAMPLE2_PREFIX|$_TF_EXAMPLE2_PREFIX|" \
examples/2.VPC-Across-VPN/provider.tf.template >examples/2.VPC-Across-VPN/provider.tf
```

- For example : 3.PSC

```
export _TF_SERVICE_ACCOUNT="<ENTER THE SERVICE ACCOUNT HERE>"
export _TF_BUCKET_NAME="<ENTER THE GCS BUCKET NAME HERE>"
export _TF_EXAMPLE3_PREFIX="<ENTER THE GCS PREFIX NAME HERE>"

sed \
-e "s|ENTER_TF_SERVICE_ACCOUNT|$_TF_SERVICE_ACCOUNT|" \
-e "s|ENTER_TF_BUCKET_NAME|$_TF_BUCKET_NAME|" \
-e "s|ENTER_TF_EXAMPLE3_PREFIX|$_TF_EXAMPLE3_PREFIX|" \
examples/3.PSC/provider.tf.template >examples/3.PSC/provider.tf
```

- For example : 4.PSC-Across-VPN

```
export _TF_SERVICE_ACCOUNT="<ENTER THE SERVICE ACCOUNT HERE>"
export _TF_BUCKET_NAME="<ENTER THE GCS BUCKET NAME HERE>"
export _TF_EXAMPLE4_PREFIX="<ENTER THE GCS PREFIX NAME HERE>"

sed \
-e "s|ENTER_TF_SERVICE_ACCOUNT|$_TF_SERVICE_ACCOUNT|" \
-e "s|ENTER_TF_BUCKET_NAME|$_TF_BUCKET_NAME|" \
-e "s|ENTER_TF_EXAMPLE4_PREFIX|$_TF_EXAMPLE4_PREFIX|" \
examples/4.PSC-Across-VPN/provider.tf.template >examples/4.PSC-Across-VPN/provider.tf
```


- For example : ServiceNetworkingAcrossInterconnect

```
export _TF_SERVICE_ACCOUNT="<ENTER THE SERVICE ACCOUNT HERE>"
export _TF_BUCKET_NAME="<ENTER THE GCS BUCKET NAME HERE>"
export _TF_EXAMPLE_SNAcrossIC_PREFIX="<ENTER THE GCS PREFIX NAME HERE>"

sed \
-e "s|ENTER_TF_SERVICE_ACCOUNT|$_TF_SERVICE_ACCOUNT|" \
-e "s|ENTER_TF_BUCKET_NAME|$_TF_BUCKET_NAME|" \
-e "s|ENTER_TF_EXAMPLE_SNAcrossIC_PREFIX|$_TF_EXAMPLE_SNAcrossIC_PREFIX|" \
examples/ServiceNetworkingAcrossInterconnect/provider.tf.template >examples/ServiceNetworkingAcrossInterconnect/provider.tf
```

## Testing

For running the unit and integration tests, please refer to the [testing](./test/README.md) documentation.

## Installation

### Terraform

Be sure you have the correct Terraform version (1.5.x), you can choose the [binary](https://releases.hashicorp.com/terraform/).

### gcloud SDK

Ensure you have gcloud already [installed](https://cloud.google.com/sdk/docs/install).

## Contributions

If you have any specific suggestions or ideas then feel free to share that by creating an issue against this repo.
