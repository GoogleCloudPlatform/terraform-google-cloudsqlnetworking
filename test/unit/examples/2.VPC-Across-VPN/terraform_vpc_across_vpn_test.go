package vpc_across_vpn_test

import (
	"testing"
	"golang.org/x/exp/slices"
	"github.com/stretchr/testify/assert"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

const terraformDirectoryPath = "../../../../examples/2.VPC-Across-VPN";

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
		PlanFilePath: "./plan",
		NoColor: true,
		VarFiles: [] string {"dev.tfvars" },
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
		PlanFilePath: "./plan",
		NoColor: true,
		VarFiles: [] string {"dev.tfvars" },
	})

	//plan *PlanStruct
	planStruct := terraform.InitAndPlan(t, terraformOptions)

	resourceCount := terraform.GetResourceCount(t, planStruct)
	assert.Equal(t,91,resourceCount.Add)
	assert.Equal(t,0,resourceCount.Change)
	assert.Equal(t,0,resourceCount.Destroy)
}

func TestTerraformModuleResourceAddressListMatch(t *testing.T) {
	// Construct the terraform options with default retryable errors to handle the most common
	// retryable errors in terraform testing.
	expectedModulesAddress := [] string {"module.google_compute_instance.module.compute_instance","module.sql-db.module.mysql[0]","module.gce_sa","module.host_project_vpn","module.user_project_vpn","module.host-vpc","module.terraform_service_accounts","module.user-vpc","module.project_services.module.project_services","module.firewall_rules.module.firewall_rules","module.user_project_services.module.project_services","module.user_google_compute_instance.module.compute_instance","module.user_gce_sa","module.host_project.module.project_services","module.sql-db","module.user_firewall_rules.module.firewall_rules","module.user_google_compute_instance.module.instance_template","module.google_compute_instance.module.instance_template"}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		PlanFilePath: "./plan",
		NoColor: true,
		VarFiles: [] string {"dev.tfvars" },
	})

	//plan *PlanStruct
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
