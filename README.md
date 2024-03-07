# CloudSQL Easy Networking


Streamline your Google Cloud SQL instance deployment and management with this comprehensive Terraform module repository. This repository simplifies Google Cloud networking configuration for Cloud SQL instances. It bundles Terraform modules to make it easier to create and manage Cloud SQL instances and all dependent resources, such as cloud networking resources, IAM policies and service accounts.

This repository contains easy-to-use Terraform modules that would help you set up all the prerequisites components required to consume Cloud SQL with private IP in a GCP project. The modules make it easy to manage Cloud SQL and all relevant GCP resources. This solution will help different users like database administrators and application engineers who want to quickly configure Cloud SQL with cloud networking.

Here are some specific benefits of using this repository:

- **Simplified configuration:** The Terraform modules in this repository abstract away the complexity of configuring Google Cloud networking for Cloud SQL. This makes it easier to get started, even for users with less experience with Google Cloud.

- **Consistency and repeatability:** Terraform ensures that your Cloud SQL configuration is consistent and repeatable across multiple environments. This can help to reduce errors and improve operational efficiency.

- **Flexibility and scalability:** The Terraform modules in this repository are flexible and scalable enough to meet the needs of a wide range of deployments, from small development projects to large enterprise applications.
If you are looking for a simple and efficient way to configure Cloud SQL with private IP in Google Cloud, this repository is the perfect place to start.

## Example Use Case

Imagine you are a database administrator responsible for setting up a new Cloud SQL instance for a production application. You need to configure the instance to use a private IP so that it is only accessible from within your VPC network.

Using the Terraform modules in this repository, you can easily configure the Cloud SQL instance and all the necessary network resources with just a few lines of code. This would typically involve the following steps:

1. Create a new Terraform configuration file and define the variables for your Cloud SQL instance and network configuration.
2. Initialize Terraform and download the necessary providers and modules.
3. Plan and apply the Terraform configuration.

Once the Terraform configuration has been applied, your Cloud SQL instance will be up and running with a private IP. You will then be able to connect to the instance from within your VPC network using the private IP address.

## Supported Usage

Many examples are included in the [examples](./examples/) folder which describes the complete end-to-end examples covering different scenarios along with its implementation guide and architecture design.

1. [Host Service Project Scenario](./examples/1.Host-Service-Project) : This solution guides a user through the steps to establish a host and a service project, create a Cloud SQL instance and a VM instance in the service project, and connect the VM instance to the Cloud SQL instance using the VM's private IP address. The host project contains the VPC, subnets, and firewall rules.


2. This section outlines two user journeys for Service Networking: using **HA-VPN** or **Dedicated Interconnect**.

    2.1 [VPC across VPN Tunnel Scenario](./examples/2.VPC-Across-VPN) : This solution guides a user to create a highly available (HA) VPN connection between a user project and a host project with a service project attached. The solution then establishes a Cloud SQL connection using the private IP address of the Cloud SQL instance created in the service project and a VM instance created in the user project.

    2.2. [Service Networking Across Interconnect](./examples/ServiceNetworkingAcrossInterconnect) : This solution helps the user to create a dedicated interconnect connection with redundancy & 99.9% availability between an on-premise environment and a GCP project. The solution establishes a Cloud SQL connection using the private IP address of a Cloud SQL instance created inside a GCP project and a VM instance present in the user on-premise environment.

3. [PSC across VPN Scenario](./examples/3.PSC) : This solution guides a user to create a PSC enabled Cloud SQL instance with a consumer and producer project setup having a compute VM instance created in the consumer project connecting to the Cloud SQL instance through PSC service endpoint. The consumer project contains the consumer VPC, service endpoint & firewall rules to connect to the SQL instance in the producer project.

4. [PSC across VPN Scenario](./examples/4.PSC-Across-VPN) : This solution guides a user to create a HA VPN connection between user and consumer project to connect to a PSC enabled Cloud SQL instance in a producer project from a compute VM instance through PSC service endpoint. The compute instance in the consumer project connects to the PSC service endpoint via the VPN connection. The PSC service endpoint connects to the Cloud SQL instance.

All [GA offerings](https://cloud.google.com/sql/docs/sqlserver/private-ip) for Cloud SQL (MySQL, PostgreSQL, MSSQL) are supported through our solutions.

## Variables

To control module's behavior, change variables' values. Behavior of each of these variables has been documented in the readme file for the respective examples. Some of the common variables has been described below but the respective readme file should act as reference point as it describe about them in more details:

- `cloudsql_instance_name` - set this variable to a string value and the cloud sql instance will be named after this variable.
- `database_version` - set this variable to a specific database version for the cloud sql instance. This [link](https://cloud.google.com/sql/docs/mysql/db-versions) talks more about the database version.
- `host_project_id` - set this variable to the GCP project id which will be acting as the host project for the setup and the necessary VPC network and subnetwork will be created here(if they do not exists already)
- `region` - set this variable to the GCP region and this will be used to create resources like Cloud SQL instance in GCP. This [link](https://cloud.google.com/compute/docs/regions-zones) talks more about the region.
- `create_network` - this is a boolean variable and set this to `true` if you want to create a new network else set to `false` if you plan to use existing network.
- `create_subnetwork` - this is a boolean variable and set this to `true` if you want to create a new subnetwork else set to `false` if you plan to use existing network.
- `network_name` - set this to the vpc name and either a new network will be created or an existing network will be used depending on the value of `create_network`.
- `subnetwork_name` - set this to the subnetwork name and either a new sub network will be created or an existing sub network will be used depending on the value of `create_subnetwork`.
- `deletion_protection` - set this flag to `true` if you want to enable deletion protection which protects from accidental protections else set to `false`.
- `access_config` - set this to `null` if you do not want the compute instance to have an exernal IP
- `create_nat` - set this to `true` if you want to have a [Cloud NAT](https://cloud.google.com/nat/docs/overview) configured else set this to `false`

**Note** - Having `access_config = null` and `create_nat = false` will disable the compute instance ability to download from public internet and you will not be able to download any client etc however the connections to other VM and Cloud sql instance using private IP will continue to work.

## State Files

Terraform must store state about your managed infrastructure and configuration. It is recommended to use the GCS bucket or an equivalent storage place where the terraform stores its state file. More details at [link](https://developer.hashicorp.com/terraform/language/state).

### provider.tf file

Each providers.tf file carries information like service account to be impersonated, bucket to be used for storing the GCS bucket and prefix that will be used in this bucket for storing the terraform state.

Following are sample commands that can be used to update the existing provider.tf.template file to create an provider.tf file.

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

For running the unit test cases and integration test cases, please refer to the [test](./test/README.md) documentation.


## Installation

### Terraform
Be sure you have the correct Terraform version (1.5.x), you can choose the binary here:
- https://releases.hashicorp.com/terraform/

### gcloud
Ensure you have gcloud already installed [gcloud cli installation steps](https://cloud.google.com/sdk/docs/install)

## Contributions

If you have any specific suggestions or ideas then feel free to share that by creating an issue against this repo.
