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

package host_service_test

import (
	"os"
	"fmt"
	"testing"
	"golang.org/x/exp/slices"
	"github.com/stretchr/testify/assert"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

const terraformDirectoryPath = "../../../../examples/1.Host-Service-Project";
var host_project_id           = os.Getenv("TF_VAR_host_project_id");
var service_project_id        = os.Getenv("TF_VAR_service_project_id");
var cloudsql_instance_name    = "cn-sqlinstance10-test";
var network_name              = "cloudsql-easy"
var subnetwork_ip_cidr        = "10.2.0.0/16"
var subnetwork_name           = "cloudsql-easy-subnet";
var region                    = "us-central1";
var zone 											= "us-central1-a";
var test_dbname               = "test_db"
var database_version 					= "MYSQL_8_0"
var deletion_protection       = false

 var tfVars = map[string]interface{}{
	"host_project_id"            : host_project_id,
	"service_project_id"         : service_project_id,
	"database_version"           : database_version,
	"cloudsql_instance_name"     : cloudsql_instance_name,
	"region"                     : region,
	"zone"                       : zone,
	"create_network"             : true,
	"create_subnetwork"          : true,
	"network_name"               : network_name,
	"subnetwork_name"            : subnetwork_name,
	"subnetwork_ip_cidr"         : subnetwork_ip_cidr,
	"deletion_protection" 			 : deletion_protection,
	"test_dbname"                : test_dbname,
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
		PlanFilePath: "./plan",
		Reconfigure : true,
		Lock: true,
		NoColor: true,
		//VarFiles: [] string {"dev.tfvars" },
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
	planExitCode, err := terraform.InitAndPlanWithExitCodeE(t, terraformOptions)
	if err != nil  {
		fmt.Println("==Error==")
		fmt.Print(err.Error())
	}
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
		//VarFiles: [] string {"dev.tfvars" },
	})

	//plan *PlanStruct
	planStruct := terraform.InitAndPlan(t, terraformOptions)

	resourceCount := terraform.GetResourceCount(t, planStruct)
	assert.Equal(t,51,resourceCount.Add)
	assert.Equal(t,0,resourceCount.Change)
	assert.Equal(t,0,resourceCount.Destroy)
}

func TestTerraformModuleResourceAddressListMatch(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	expectedModulesAddress := [] string {"module.firewall_rules.module.firewall_rules","module.gce_sa","module.host-vpc","module.terraform_service_accounts","module.google_compute_instance","module.sql-db.module.mysql[0]","module.host_project.module.project_services","module.project_services.module.project_services","module.sql-db","module.nat[0]"}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars : tfVars,
		PlanFilePath: "./plan",
		Reconfigure : true,
		Lock: true,
		NoColor: true,
		//VarFiles: [] string {"dev.tfvars" },
	})

	//plan *PlanStruct
	planStruct := terraform.InitAndPlanAndShow(t, terraformOptions)
	content, err := terraform.ParsePlanJSON(planStruct)
	print("\n\n")
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
