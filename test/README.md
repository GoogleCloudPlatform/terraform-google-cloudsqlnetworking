## Testing

The following sections describe how the examples can be tested in a Google Cloud Platform environment. There would be unit & integration tests for each example in this directory.

### Tests - Environment Variables

While running these tests locally (or in your development machines) make sure you have declared following as the environment variables.

- For [HostServiceProject](../examples/1.Host-Service-Project/) and [VPCAcrossVPN](../examples/2.VPC-Across-VPN/) :

```
export TF_VAR_host_project_id=<HOST_PROJECT_ID>
export TF_VAR_service_project_id=<SERVICE_PROJECT_ID>
export TF_VAR_user_project_id=<USER_PROJECT_ID>
```

- For [PSC](../examples/3.PSC/) and [PSC-Across-VPN](../examples/4.PSC-Across-VPN/) :

```
export TF_VAR_consumer_project_id=<CONSUMER_PROJECT_ID>
export TF_VAR_producer_project_id=<PRODUCER_PROJECT_ID>
export TF_VAR_user_project_id=<USER_PROJECT_ID>
```

- For [ServiceNetworkingAcrossInterconnect](../examples/ServiceNetworkingAcrossInterconnect/) :

```
export TF_VAR_interconnect_project_id=<INTERCONNECT_PROJECT_ID>
export TF_VAR_host_project_id=<HOST_PROJECT_ID>
export deployed_interconnect_name=<DEPLOYED_INTERCONNECT_NAME>
```

### Unit Testing

#### Running all the unit tests

To run all the tests/functions under the unit testing directory for the example, please follow these steps:

1. cd REPO_NAME
2. go mod init test
3. go mod tidy
4. **go test -v -json ./... | ./test-summary**

    **Note :** [test-summary](https://pkg.go.dev/gocloud.dev/internal/testing/test-summary) is used to provide summary of the test results.

#### Running specific tests

To run specific tests/functions under the unit testing directory for the example, please follow these steps:

1. cd /tests/unit/examples/SCENARIO_NAME
2. go mod init test
3. go mod tidy
4. **go test -timeout 15m -v**

    **e.g.** Here is an example demonstrating how a unit test for example 1 can be executed :
    ```
    cd /tests/unit/examples/1.Host-Service-Project
    go mod init test
    go mod tidy
    go test -timeout 12m -v
    ```

### Integration Testing

#### Running all the integration tests

To run all the tests/functions under the integration testing directory for the example, please follow these steps:

1. cd REPO_NAME
2. go mod init test
3. go mod tidy
4. **go test -v -json -timeout 60m ./... | ./test-summary**

    **Note :** [test-summary](https://pkg.go.dev/gocloud.dev/internal/testing/test-summary) is used to provide summary of the test results.


#### Running specific tests

To run specific tests/functions under the integration testing directory for the example, please follow these steps:

1. cd /tests/integration/examples/SCENARIO_NAME
2. go mod init test
3. go mod tidy
4. **go test -timeout 60m -v**

    **e.g.** Here is an example demonstrating how an integration test for scenario 1 can be executed:

    ```
    cd /tests/integration/examples/1.Host-Service-Project
    go mod init test
    go mod tidy
    go test -timeout 60m -v
    ```