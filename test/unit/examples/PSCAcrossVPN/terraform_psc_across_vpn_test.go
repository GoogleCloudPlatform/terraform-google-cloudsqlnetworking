// Copyright 2023-2024 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vpnacrosspsctest

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

// Setting required variables for the unit tests.

const terraformDirectoryPath = "../../../../examples/4.PSC-Across-VPN"

var consumerProjectID = os.Getenv("TF_VAR_consumer_project_id")
var producerProjectID = os.Getenv("TF_VAR_producer_project_id")
var userProjectID = os.Getenv("TF_VAR_user_project_id")
var cloudSQLInstanceName = "cn-sqlinstance10-test"
var databaseVersion = "MYSQL_8_0"
var consumerNetworkName = "cloudsql-easy"
var consumerSubnetworkName = "cloudsql-easy-subnet"
var region = "us-central1"
var zone = "us-central1-a"
var testDbname = "test_db"
var natName = "cloudsql-nat"
var routerName = "cloudsql-router"
var endpointIP = []string{"10.0.0.5"}
var uservpcNetworkName = "cloudsql-user"
var uservpcSubnetworkName = "cloudsql-user-subnet"
var consumerSubnetworkIPCIDR = "192.168.0.0/24"
var uservpcSubnetworkIPCidr = "10.0.0.0/24"
var user_region = "us-central1"
var user_zone = "us-central1-a"
var deletionProtection = false
var tfVars = map[string]any{
	"consumer_project_id":        consumerProjectID,
	"producer_project_id":        producerProjectID,
	"database_version":           databaseVersion,
	"cloudsql_instance_name":     cloudSQLInstanceName,
	"region":                     region,
	"zone":                       zone,
	"create_network":             true,
	"create_subnetwork":          true,
	"user_network_name":          uservpcNetworkName,
	"user_subnetwork_name":       uservpcSubnetworkName,
	"consumer_network_name":      consumerNetworkName,
	"consumer_subnetwork_name":   consumerSubnetworkName,
	"consumer_cidr":              consumerSubnetworkIPCIDR,
	"user_project_id":            userProjectID,
	"create_user_vpc_network":    true,
	"create_user_vpc_subnetwork": true,
	"user_cidr":                  uservpcSubnetworkIPCidr,
	"test_dbname":                testDbname,
	"deletion_protection":        deletionProtection,
	"nat_name":                   natName,
	"router_name":                routerName,
	"endpoint_ip":                endpointIP,
	"user_region":                user_region,
	"user_zone":                  user_zone,
}

// TestInitAndPlanRunWithTfVars runs init and plan with tfvars variable.
func TestInitAndPlanRunWithTfVars(t *testing.T) {
	/*
	 0 = Succeeded with empty diff (no changes)
	 1 = Error
	 2 = Succeeded with non-empty diff (changes present)
	*/
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})

	planExitCode := terraform.InitAndPlanWithExitCode(t, terraformOptions)
	assert.Equal(t, 2, planExitCode)
}

// TestInitAndPlanRunWithoutTfVarsExpectFailureScenario runs init and plan without tfvars variable.
// This plan should fail and return an error as expected.

func TestInitAndPlanRunWithoutTfVarsExpectFailureScenario(t *testing.T) {
	/*
	 0 = Succeeded with empty diff (no changes)
	 1 = Error
	 2 = Succeeded with non-empty diff (changes present)
	*/
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})
	planExitCode := terraform.InitAndPlanWithExitCode(t, terraformOptions)
	assert.Equal(t, 1, planExitCode)
}

// TestResourcesCount should count the number of resources getting created and match it with the expected values
func TestResourcesCount(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)

	resourceCount := terraform.GetResourceCount(t, planStruct)
	assert.Equal(t, 81, resourceCount.Add)
	assert.Equal(t, 0, resourceCount.Change)
	assert.Equal(t, 0, resourceCount.Destroy)
}

// TestTerraformModuleResourceAddressListMatch compares the modules expected and being created by the test run
func TestTerraformModuleResourceAddressListMatch(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	expectedModulesAddress := []string{"module.consumer_project_vpn", "module.terraform_service_accounts", "module.user_project_vpn", "module.compute_address", "module.consumer_vpc", "module.user_project.module.project_services", "module.consumer_project.module.project_services", "module.user_nat[0]", "module.producer_project.module.project_services", "module.sql_db.module.mysql[0]", "module.firewall_rules.module.firewall_rules", "module.user_gce_sa", "module.user_firewall_rules.module.firewall_rules", "module.sql_db", "module.user_vpc", "module.user_project_instance"}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
		NoColor:      true,
		//VarFiles: [] string {"dev.tfvars" },
	})

	planStruct := terraform.InitAndPlanAndShow(t, terraformOptions)
	content, err := terraform.ParsePlanJSON(planStruct)
	actualModuleAddress := make([]string, 0)
	for _, element := range content.ResourceChangesMap {
		if !slices.Contains(actualModuleAddress, element.ModuleAddress) && len(element.ModuleAddress) > 0 {
			actualModuleAddress = append(actualModuleAddress, element.ModuleAddress)
		}
	}
	if err != nil {
		t.Errorf(err.Error())
	}
	assert.ElementsMatch(t, expectedModulesAddress, actualModuleAddress)
}
