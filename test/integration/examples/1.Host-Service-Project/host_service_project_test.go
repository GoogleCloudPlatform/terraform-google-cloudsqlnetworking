package host_service_project_test

import (
	"fmt"
	"testing"
	//"slices"
	"github.com/stretchr/testify/assert"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

const terraformDirectoryPath = "../../../../../cloudsql-easy-networking/examples/1.Host-Service-Project";



// name the function as Test*
func TestMySqlPrivateModule(t *testing.T) {

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir: terraformDirectoryPath,
		//PlanFilePath: "./plan",
		NoColor: true,
		VarFiles: [] string {"dev.tfvars" },
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	//defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "host_vpc_name")
	fmt.Print("================")
	fmt.Print(output)
	fmt.Print("================")
	assert.Equal(t, "cloudsql-easy", output)
	output = terraform.Output(t, terraformOptions, "host_subnetwork_id")
	assert.Equal(t, "projects/pm-singleproject-20/regions/us-central1/subnetworks/cloudsql-easy-subnet", output)
}

