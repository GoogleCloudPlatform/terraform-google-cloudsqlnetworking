# CloudSQL Easy Networking

This repository contains the easy to use terraform modules that helps to setup all the pre-requsities required to consume the cloud sql with private ip in a GCP project.

Complete end-to-end examples covering different scenarios along with its implementation guide has been documented in the `examples/` folder .

## Testing

### Unit Test

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



