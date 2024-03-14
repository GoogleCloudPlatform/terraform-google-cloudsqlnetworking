// Copyright 2023 Google LLC
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

package vpcacrossvpntest

import (
	"os"
	"testing"
	"golang.org/x/exp/slices"
	"github.com/stretchr/testify/assert"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

const terraformDirectoryPath   = "../../../../examples/2.VPC-Across-VPN";
var hostProjectID              = os.Getenv("TF_VAR_host_project_id");
var serviceProjectID           = os.Getenv("TF_VAR_service_project_id");
var databaseVersion 				   = "MYSQL_8_0"
var region                     = "us-central1";
var zone										   = "us-central1-a";
var userProjectID              = os.Getenv("TF_VAR_user_project_id");
var cloudSQLInstanceName       = "cn-sqlinstance10-test";
var networkName                = "cloudsql-easy";
var subnetworkName             = "cloudsql-easy-subnet";
var subnetworkIPCidr           = "10.2.0.0/16"
var uservpcNetworkName         = "cloudsql-user"
var uservpcSubnetworkName      = "cloudsql-user-subnet"
var uservpcSubnetworkIPCidr    = "10.10.30.0/24"
var testDbname 							   = "test_db"
var userRegion                 = "us-west1"
var userZone                   = "us-west1-a"
var deletionProtection 			   = false
var tfVars = map[string]any{
	"host_project_id"            : hostProjectID,
	"service_project_id"         : serviceProjectID,
	"database_version"           : databaseVersion,
	"cloudsql_instance_name"     : cloudSQLInstanceName,
	"region"                     : region,
	"zone"                       : zone,
	"create_network"             : true,
	"create_subnetwork"          : true,
	"network_name"               : networkName,
	"subnetwork_name"            : subnetworkName, // this subnetwork will be created
	"subnetwork_ip_cidr"         : subnetworkIPCidr,
	"user_project_id"            : userProjectID,
	"user_region"                : userRegion,
	"user_zone"                  : userZone,
	"create_user_vpc_network"    : true,
	"create_user_vpc_subnetwork" : true,
	"uservpc_network_name"       : uservpcNetworkName,
	"uservpc_subnetwork_name"    : uservpcSubnetworkName,
	"uservpc_subnetwork_ip_cidr" : uservpcSubnetworkIPCidr,
	"test_dbname"                : testDbname,
	"deletion_protection" 	     : deletionProtection,

}

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
		Vars : tfVars,
		Reconfigure : true,
		Lock: true,
		PlanFilePath: "./plan",
		NoColor: true,
	})

	planExitCode := terraform.InitAndPlanWithExitCode(t, terraformOptions)
	assert.Equal(t, 2, planExitCode)
}

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
		Reconfigure : true,
		Lock: true,
		PlanFilePath: "./plan",
		NoColor: true,
	})
	planExitCode := terraform.InitAndPlanWithExitCode(t, terraformOptions)
	assert.Equal(t, 1, planExitCode)
}

func TestResourcesCount(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars : tfVars,
		Reconfigure : true,
		Lock: true,
		PlanFilePath: "./plan",
		NoColor: true,
	})

	planStruct := terraform.InitAndPlan(t, terraformOptions)

	resourceCount := terraform.GetResourceCount(t, planStruct)
	assert.Equal(t,91,resourceCount.Add)
	assert.Equal(t,0,resourceCount.Change)
	assert.Equal(t,0,resourceCount.Destroy)
}

func TestTerraformModuleResourceAddressListMatch(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	expectedModulesAddress := [] string {"module.google_compute_instance","module.sql-db.module.mysql[0]","module.gce_sa","module.host_project_vpn","module.user_project_vpn","module.host-vpc","module.terraform_service_accounts","module.user-vpc","module.project_services.module.project_services","module.firewall_rules.module.firewall_rules","module.user_project_services.module.project_services","module.user_google_compute_instance","module.user_gce_sa","module.host_project.module.project_services","module.sql-db","module.user_firewall_rules.module.firewall_rules","module.user-nat[0]"}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars : tfVars,
		Reconfigure : true,
		Lock: true,
		PlanFilePath: "./plan",
		NoColor: true,
	})

	planStruct := terraform.InitAndPlanAndShow(t, terraformOptions)
	content, err := terraform.ParsePlanJSON(planStruct)
	actualModuleAddress := make([]string, 0)
	for _, element := range content.ResourceChangesMap {
		if !slices.Contains(actualModuleAddress, element.ModuleAddress) && len(element.ModuleAddress) > 0 {
			actualModuleAddress = append(actualModuleAddress,element.ModuleAddress)
		}
	}
	if err != nil {
		print(err.Error())
	}
	assert.ElementsMatch(t, expectedModulesAddress, actualModuleAddress);
}
