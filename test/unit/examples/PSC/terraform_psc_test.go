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

package sqlpsctest

import (
	"log"
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

// Setting required variables for the unit tests.

const terraformDirectoryPath = "../../../../examples/3.PSC"

var consumerProjectID = os.Getenv("TF_VAR_consumer_project_id")
var producerProjectID = os.Getenv("TF_VAR_producer_project_id")
var cloudSQLInstanceName = "cloudsql-easy-networking"
var consumerNetworkName = "cloudsql-easy"
var subnetworkIPCidr = "10.0.0.0/16"
var consumerSubnetworkName = "cloudsql-easy-subnet"
var region = "us-central1"
var zone = "us-central1-a"
var testDbname = "test"
var databaseVersion = "MYSQL_8_0"
var deletionProtection = false
var reservedIps = []string{"10.0.0.5"}
var tfVars = map[string]any{
	"consumer_project_id":      consumerProjectID,
	"producer_project_id":      producerProjectID,
	"database_version":         databaseVersion,
	"cloudsql_instance_name":   cloudSQLInstanceName,
	"region":                   region,
	"zone":                     zone,
	"consumer_network_name":    consumerNetworkName,
	"consumer_subnetwork_name": consumerSubnetworkName,
	"consumer_cidr":            subnetworkIPCidr,
	"deletion_protection":      deletionProtection,
	"test_dbname":              testDbname,
	"reserved_ips":             reservedIps,
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
		PlanFilePath: `./plan`,
		Reconfigure:  true,
		Lock:         true,
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
		PlanFilePath: `./plan`,
		NoColor:      true,
	})
	planExitCode, err := terraform.InitAndPlanWithExitCodeE(t, terraformOptions)
	if err != nil {
		log.Println("==Error==")
		log.Print(err.Error())
	}
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
		PlanFilePath: `./plan`,
		NoColor:      true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)

	resourceCount := terraform.GetResourceCount(t, planStruct)
	assert.Equal(t, 48, resourceCount.Add)
	assert.Equal(t, 0, resourceCount.Change)
	assert.Equal(t, 0, resourceCount.Destroy)
}

// TestTerraformModuleResourceAddressListMatch compares the modules expected and being created by the test run

func TestTerraformModuleResourceAddressListMatch(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	expectedModulesAddress := []string{"module.sql_db", "module.sql_db.module.mysql[0]", "module.compute_address", "module.consumer_nat[0]", "module.google_compute_instance", "module.gce_sa", "module.consumer_vpc", "module.consumer_project.module.project_services", "module.producer_project.module.project_services", "module.firewall_rules.module.firewall_rules", "module.terraform_service_accounts"}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		PlanFilePath: `./plan`,
		Reconfigure:  true,
		Lock:         true,
		NoColor:      true,
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
