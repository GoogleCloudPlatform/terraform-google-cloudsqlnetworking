# CloudSQL Easy Networking

This repository contains the easy to use terraform modules that helps to setup all the pre-requsities required to consume the cloud sql with private ip in a GCP project.

Complete end-to-end examples covering different scenarios along with its implementation guide has been documented in the `examples/` folder .

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


### Provider.tf file

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
