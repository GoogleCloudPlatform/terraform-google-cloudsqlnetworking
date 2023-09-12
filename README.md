# CloudSQL Easy Networking

This repository makes it easier to create  & manage the Google Cloud SQL instance and all the other dependent resources like cloud networking resources, IAM, service accounts etc.

This repository contains the easy to use terraform modules that helps to setup all the pre-requsities required to consume the cloud sql with private ip in a GCP project.

The modules makes it easy to manage to the cloud sql and all the relevant gcp resources.

## Usage

Many examples are included in the [examples](./examples/) folder which describes the complete end-to-end examples covering different scenarios along with its implementation guide and architecture design.

1. [Host Service Project Scenario](./examples/1.Host-Service-Project)
2. [VPC across VPN Tunnel Scenario](./examples/2.VPC-Across-VPN)

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

1. For example1 - Host- Service project scenario

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

1. For example2 - VPC Across VPN scenario

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

## Testing


### Unit Test


#### Running all the unit test

1. cd REPO_NAME
2. go mod init test
3. go mod tidy
4. **go test -v -json ./... | ./test-summary**

    **Note :** test-summary is used to provide summary of the test results. More details [here](https://pkg.go.dev/gocloud.dev/internal/testing/test-summary)


#### Running Example specific test
1. cd /tests/unit/examples/SCENARIO_NAME
2. go mod init test
3. go mod tidy
4. **go test -timeout 12m -v**

    **e.g.** Here is an example demonstrating how a unit test for scenario 1 can be executed
    ```
    cd /tests/unit/examples/1.Host-Service-Project
    go mod init test
    go mod tidy
    go test -timeout 12m -v
    ```

### Integration Test


#### Running all the integration test

1. cd REPO_NAME
2. go mod init test
3. go mod tidy
4. **go test -v -json -timeout 60m ./... | ./test-summary**

    **Note :** test-summary is used to provide summary of the test results. More details [here](https://pkg.go.dev/gocloud.dev/internal/testing/test-summary)


#### Running Example specific test
1. cd /tests/unit/examples/SCENARIO_NAME
2. go mod init test
3. go mod tidy
4. **go test -timeout 60m -v**

    **e.g.** Here is an example demonstrating how a unit test for scenario 1 can be executed
    ```
    cd /tests/integration/examples/1.Host-Service-Project
    go mod init test
    go mod tidy
    go test -timeout 60m -v
    ```


## Installation
### Terraform
Be sure you have the correct Terraform version (1.5.x), you can choose the binary here:
- https://releases.hashicorp.com/terraform/

### gcloud
Ensure you have gcloud already installed [gcloud cli installation steps](https://cloud.google.com/sdk/docs/install)
