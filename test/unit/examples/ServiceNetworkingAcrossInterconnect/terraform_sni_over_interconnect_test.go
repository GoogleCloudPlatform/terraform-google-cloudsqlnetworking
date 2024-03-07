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

package sniinterconnecttest

import (
	"os"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

const terraformDirectoryPath = "../../../../examples/ServiceNetworkingAcrossInterconnect"

var hostProjectID = os.Getenv("TF_VAR_host_project_id")
var serviceProjectID = os.Getenv("TF_VAR_service_project_id")
var databaseVersion = "MYSQL_8_0"
var region = "us-west2"
var zone = "us-west2-a"
var cloudSQLInstanceName = "cn-sqlinstance10-test"
var networkName = "cloudsql-easy"
var subnetworkName = "cloudsql-easy-subnet"
var subnetworkIPCidr = "10.2.0.0/16"
var deletionProtection = false

// variables for Interconnect configuration
var interconnectProjectID = os.Getenv("TF_VAR_interonnect_project_id")
var firstInterconnectName = "cso-lab-interconnect-1"
var secondInterconnectName = "cso-lab-interconnect-2"
var icRouterBgpAsn = 65004

// first vlan attachment configuration values
var firstVaAsn = "65418"
var firstVaBandwidth = "BPS_1G"
var firstVaBgpRange = "169.254.61.0/29"
var firstVlanTag = 601

// second vlan attachment configuration values
var secondVaAsn = "65418"
var secondVaBandwidth = "BPS_1G"
var secondVaBgpRange = "169.254.61.8/29"
var secondVlanTag = 601

var tfVars = map[string]any{
	"host_project_id":          hostProjectID,
	"service_project_id":       serviceProjectID,
	"database_version":         databaseVersion,
	"cloudsql_instance_name":   cloudSQLInstanceName,
	"region":                   region,
	"zone":                     zone,
	"create_network":           true,
	"create_subnetwork":        true,
	"network_name":             networkName,
	"subnetwork_name":          subnetworkName, // this subnetwork will be created
	"subnetwork_ip_cidr":       subnetworkIPCidr,
	"deletion_protection":      deletionProtection,
	"interconnect_project_id":  interconnectProjectID,
	"first_interconnect_name":  firstInterconnectName,
	"second_interconnect_name": secondInterconnectName,
	"ic_router_bgp_asn":        icRouterBgpAsn,
	"first_va_asn":             firstVaAsn,
	"first_va_bandwidth":       firstVaBandwidth,
	"first_va_bgp_range":       firstVaBgpRange,
	"first_vlan_tag":           firstVlanTag,
	"second_va_asn":            secondVaAsn,
	"second_va_bandwidth":      secondVaBandwidth,
	"second_va_bgp_range":      secondVaBgpRange,
	"second_vlan_tag":          secondVlanTag,
}

/* TestInitAndPlanRunWithTfVars performs sanity check to ensure the terraform init
&& terraform plan is executed successfully and returns a valid Succeeded run code. */
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

/* TestInitAndPlanRunWithoutTfVarsExpectFailureScenario performs test runs without tfvars file
 to ensure the terraform init && terraform plan is executed unsuccessfully and returns an expected error run code. */
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

/* TestResourcesCount performs validation to verify number of  resources created, deleted and
updated. */
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
	assert.Equal(t, 58, resourceCount.Add)
	assert.Equal(t, 0, resourceCount.Change)
	assert.Equal(t, 0, resourceCount.Destroy)
}

/* TestTerraformModuleResourceAddressListMatch compares and verifies the list of resources, modules
created by the terraform solution. */
func TestTerraformModuleResourceAddressListMatch(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	expectedModulesAddress := []string{"module.google_compute_instance", "module.nat[0]", "module.sql_db.module.mysql[0]", "module.gce_sa", "module.host_vpc", "module.firewall_rules.module.firewall_rules", "module.host_project.module.project_services", "module.sql_db", "module.tf_host_prjct_service_accounts", "module.service_project[0].module.project_services", "module.tf_svc_prjct_service_accounts[0]", "module.vlan_attachment_b", "module.vlan_attachment_a"}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		Vars:         tfVars,
		Reconfigure:  true,
		Lock:         true,
		PlanFilePath: "./plan",
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
		print(err.Error())
	}
	assert.ElementsMatch(t, expectedModulesAddress, actualModuleAddress)
}
