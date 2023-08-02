package host_service_project_test

import (
	"fmt"
	"time"
	"testing"
	"github.com/tidwall/gjson"
	"github.com/stretchr/testify/assert"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/shell"
)

const terraformDirectoryPath  = "../../../../../cloudsql-easy-networking/examples/1.Host-Service-Project";
var host_project_id           = "pm-singleproject-20";
var service_project_id        = "pm-test-10-e90f";
var cloudsql_instance_name    = "cn-sqlinstance10-test";
var subnetwork_name           = "cloudsql-easy-subnet";
var region                    = "us-central1";
var test_dbname               = "test_db"
var network_name = ""
var subnetwork_ip_cidr = ""
var database_version = ""

// name the function as Test*
func TestMySqlPrivateModule(t *testing.T) {
	host_project_id          = "pm-singleproject-20";
  service_project_id       = "pm-test-10-e90f";
	network_name             = "cloudsql-easy";
	subnetwork_name          = "cloudsql-easy-subnet";
	region                   = "us-central1";
	subnetwork_ip_cidr       = "10.2.0.0/16";

	tfVars := map[string]interface{}{
		"host_project_id"            : host_project_id,
    "service_project_id"         : service_project_id,
		"region"                     : region,
		"create_subnetwork"          : true,
		"create_network"             : true,
		"network_name"               : network_name,
		"subnetwork_name"            : subnetwork_name,
		"cloudsql_instance_name"     : cloudsql_instance_name,
		"database_version"           : database_version,
		"test_dbname"                : test_dbname,
		"subnetwork_ip_cidr"         : subnetwork_ip_cidr,
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars : tfVars,
		TerraformDir: terraformDirectoryPath,
		//PlanFilePath: "./plan",
		NoColor: true,
		SetVarsAfterVarFiles: true,
		VarFiles: [] string {"dev.tfvars" },
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	//wait for 90 seconds to let resource acheive stable state
	time.Sleep(90* time.Second)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "host_vpc_name")

	fmt.Println(" ========= Verify Subnet Name ========= ")
	assert.Equal(t, network_name, output)
	fmt.Println(" ========= Verify Subnetwork Id ========= ")
	output = terraform.Output(t, terraformOptions, "host_subnetwork_id")
	cloudSqlInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	subnetworkId := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s",host_project_id,region,subnetwork_name)
	assert.Equal(t,subnetworkId , output)
	//gcloud sql instances describe cn-sqlinstance10-u9s --project pm-test-10-e90f "sql","instances","describe","cn-sqlinstance10-u9s","--project","pm-test-10-e90f"
	text := "sql"
	cmd := shell.Command{
		Command : "gcloud",
		Args : []string{text,"instances","describe",cloudSqlInstanceName,"--project="+service_project_id,"--format=json"},
	}
	op,err := shell.RunCommandAndGetOutputE(t, cmd)
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result := gjson.Parse(op)
	if err != nil {
		fmt.Sprintf("===Error %s Encountered while executing %s", err ,text)
	}
	fmt.Println(" ========= Verify if public IP is disabled ========= ")
	assert.Equal(t, false, gjson.Get(result.String(),"settings.ipConfiguration.ipv4Enabled").Bool())
	fmt.Println(" ========= Verify SQL RUNNING Instance State ========= ")
	assert.Equal(t, "RUNNABLE", gjson.Get(result.String(),"state").String())


	//Iterate through list of database to ensure a new db was created
	fmt.Println(" ====================================================== ")
	fmt.Println(" ========= Verify DB Creation ========= ")
	cmd = shell.Command{
		Command : "gcloud",
		Args : []string{"sql","databases","describe",test_dbname,"--instance="+cloudSqlInstanceName,"--project="+service_project_id,"--format=json"},
	}
	op,err = shell.RunCommandAndGetOutputE(t, cmd)
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result = gjson.Parse(op)
	if err != nil {
		fmt.Sprintf("=== Error %s Encountered while executing %s", err ,text)
	}
	assert.Equal(t, test_dbname, gjson.Get(result.String(),"name").String())
}

func TestUsingExistingNetworkMySqlPrivateModule(t *testing.T) {
	host_project_id          = "pm-host-networking";
  service_project_id       = "pm-service1-networking";
	network_name             = "host-cloudsql-easy";
	subnetwork_name          = "host-cloudsql-easy-subnet";
	region                   = "us-central1";

	tfVars := map[string]interface{}{
		"host_project_id"            : host_project_id,
    "service_project_id"         : service_project_id,
		"region"                     : region,
		"create_subnetwork"          : false,
		"create_network"             : false,
		"network_name"               : network_name,
		"subnetwork_name"            : subnetwork_name,
		"cloudsql_instance_name"     : cloudsql_instance_name,
		"database_version"           : database_version,
		"test_dbname"                : test_dbname,
	}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars : tfVars,
		TerraformDir: terraformDirectoryPath,
		//PlanFilePath: "./plan",
		NoColor: true,
		SetVarsAfterVarFiles: true,
		VarFiles: [] string {"dev.tfvars" },
	})

	//validate if the VPC already exists in host project
	text := "compute"
	cmd := shell.Command{
		Command : "gcloud",
		Args : []string{text,"networks","describe",network_name,"--project="+host_project_id,"--format=json"},
	}
	op,err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Fatalf("Expected Network : %s does not exists in Project : %s ", network_name, host_project_id)
	}

	//validate if the subnet already exists in host project
	cmd = shell.Command{
		Command : "gcloud",
		Args : []string{text,"networks","subnets","describe",subnetwork_name,"--project="+host_project_id,"--region="+region,"--format=json"},
	}
	op,err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Fatalf("Expected Sub network : %s does not exists in Project : %s ", subnetwork_name, host_project_id)
	}


	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	//wait for 90 seconds to let resource acheive stable state
	time.Sleep(90* time.Second)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "host_vpc_name")

	fmt.Println(" ========= Verify Subnet Name ========= ")
	assert.Equal(t, network_name, output)

	fmt.Println(" ========= Verify Subnetwork Id ========= ")
	output = terraform.Output(t, terraformOptions, "host_subnetwork_id")
	cloudSqlInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	subnetworkId := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s",host_project_id,region,subnetwork_name)
	assert.Equal(t,subnetworkId , output)
	//gcloud sql instances describe cn-sqlinstance10-u9s --project pm-test-10-e90f "sql","instances","describe","cn-sqlinstance10-u9s","--project","pm-test-10-e90f"
	text = "sql"
	cmd = shell.Command{
		Command : "gcloud",
		Args : []string{text,"instances","describe",cloudSqlInstanceName,"--project="+service_project_id,"--format=json"},
	}
	op,err = shell.RunCommandAndGetOutputE(t, cmd)
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result := gjson.Parse(op)
	if err != nil {
		fmt.Sprintf("===Error %s Encountered while executing %s", err ,text)
	}
	fmt.Println(" ========= Verify if public IP is disabled ========= ")
	assert.Equal(t, false, gjson.Get(result.String(),"settings.ipConfiguration.ipv4Enabled").Bool())
	fmt.Println(" ========= Verify SQL RUNNING Instance State ========= ")
	assert.Equal(t, "RUNNABLE", gjson.Get(result.String(),"state").String())


	//Iterate through list of database to ensure a new db was created
	fmt.Println(" ====================================================== ")
	fmt.Println(" ========= Verify DB Creation ========= ")
	cmd = shell.Command{
		Command : "gcloud",
		Args : []string{"sql","databases","describe",test_dbname,"--instance="+cloudSqlInstanceName,"--project="+service_project_id,"--format=json"},
	}
	op,err = shell.RunCommandAndGetOutputE(t, cmd)
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result = gjson.Parse(op)
	if err != nil {
		fmt.Sprintf("=== Error %s Encountered while executing %s", err ,text)
	}
	assert.Equal(t, test_dbname, gjson.Get(result.String(),"name").String())
}
