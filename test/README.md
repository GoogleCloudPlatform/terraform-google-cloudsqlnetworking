
## Testing

Following sections describes how the examples can be tested in a GCP environment.


### Running locally
While Running these locally(or in your development machines) make sure you have declared following as the environment variables.

- For HostServiceProject and VPCAcrossVPN :

```
export TF_VAR_host_project_id=<HOST_PROJECT_ID>
export TF_VAR_service_project_id=<SERVICE_PROJECT_ID>
export TF_VAR_user_project_id=<USER_PROJECT_ID>
```

- For PSC and PSC-Across-VPN :

```
export TF_VAR_consumer_project_id=<CONSUMER_PROJECT_ID>
export TF_VAR_producer_project_id=<PRODUCER_PROJECT_ID>
export TF_VAR_user_project_id=<USER_PROJECT_ID>
```

- For ServiceNetworkingAcrossInterconnect :

```
export TF_VAR_interconnect_project_id=<INTERCONNECT_PROJECT_ID>
export TF_VAR_host_project_id=<HOST_PROJECT_ID>
export deployed_interconnect_name=<DEPLOYED_INTERCONNECT_NAME>
```

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
