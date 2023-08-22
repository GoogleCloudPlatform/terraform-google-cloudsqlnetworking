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

const terraformDirectoryPath  = "../../../../examples/1.Host-Service-Project";
var host_project_id           = "pm-singleproject-20";
var service_project_id        = "pm-test-10-e90f";
var cloudsql_instance_name    = "cn-sqlinstance10-test";
var subnetwork_name           = "cloudsql-easy-subnet";
var region                    = "us-central1";
var zone 											= "us-central1-a";
var test_dbname               = "test_db"
var database_version 					= "MYSQL_8_0"
var deletion_protection       = false;
// var backendConfig  						=  map[string]interface{}{
// 	"impersonate_service_account" : "iac-sa-test@pm-singleproject-20.iam.gserviceaccount.com",
// 	"bucket" 											: "pm-cncs-cloudsql-easy-networking",
// 	"prefix" 											: "test/example1",
//  }

/*
This test creates all the resources including the vpc network and subnetwork along with other
resources like cloudsql, computeinstance, service networking etc.

It then validates if
1. cloud sql only contains the privateIP
2. compute instance in service project is able to connect and perform operation on cloudsql instance
3. validates the existence of the network and subnetwork
*/
func TestMySqlPrivateModule(t *testing.T) {
	var iteration int;
	host_project_id          = "pm-singleproject-20";
  service_project_id       = "pm-test-10-e90f";
	network_name             := "cloudsql-easy";
	subnetwork_name          = "cloudsql-easy-subnet";
	subnetwork_ip_cidr       := "10.2.0.0/16";

	tfVars := map[string]interface{}{
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

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars : tfVars,
		//BackendConfig : backendConfig,
		TerraformDir  : terraformDirectoryPath,
		//PlanFilePath : "./plan",
		Reconfigure : true,
		NoColor : true,
		SetVarsAfterVarFiles: true,
		//VarFiles: [] string {"dev.tfvars" }, //an additional tfvars containing configuration files can be passed here
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	//wait for 60 seconds to let resource acheive stable state
	time.Sleep(60* time.Second)

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
	iteration = 0;

	// performs iterations for 3 times to check if the database gets created or not
	for {
		cmd = shell.Command{
			Command : "gcloud",
			Args : []string{"sql","databases","describe",test_dbname,"--instance="+cloudSqlInstanceName,"--project="+service_project_id,"--format=json"},
		}
		op,err = shell.RunCommandAndGetOutputE(t, cmd)
		if err == nil || iteration > 3 {
			break
		} else {
			fmt.Printf("Database with Database Name %s not found in cloud sql instance %s in project %s, will reattempt in few sec", test_dbname, cloudSqlInstanceName, service_project_id)
		}
		time.Sleep(60 * time.Second)
		iteration++;
	}

	if err != nil {
		t.Fatalf("Expected Database Name : %s at Cloudsql Instance :%s does not exists in Project : %s ", test_dbname, cloudSqlInstanceName, service_project_id)
	}
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result = gjson.Parse(op)
	if err != nil {
		fmt.Sprintf("=== Error %s Encountered while executing %s", err ,text)
	}
	assert.Equal(t, test_dbname, gjson.Get(result.String(),"name").String())
}

/*
This test creates consumes existing VPC and subnetwork but include all the other
resources like cloudsql, computeinstance, service networking etc.

It then validates if
1. cloud sql only contains the privateIP
2. compute instance in service project is able to connect and perform operation on cloudsql instance
3. validates the existence of the network and subnetwork
*/
func TestUsingExistingNetworkMySqlPrivateModule(t *testing.T) {
	var iteration int;

	host_project_id          = "pm-singleproject-20";
  service_project_id       = "pm-test-10-e90f";
	network_name             := "host-cloudsql-easy";
	subnetwork_name          := "host-cloudsql-easy-subnet";

	tfVars := map[string]interface{}{
		"host_project_id"            : host_project_id,
    "service_project_id"         : service_project_id,
		"database_version"           : database_version,
		"cloudsql_instance_name"     : cloudsql_instance_name,
		"region"                     : region,
		"zone"                       : zone,
		"create_network"             : false,
		"create_subnetwork"          : false,
		"network_name"               : network_name,
		"subnetwork_name"            : subnetwork_name,
		"deletion_protection" 			 : deletion_protection,
		"test_dbname"                : test_dbname,
	}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars : tfVars,
		TerraformDir: terraformDirectoryPath,
		//BackendConfig : backendConfig,
		Reconfigure : true,
		//PlanFilePath: "./plan",
		NoColor: true,
		SetVarsAfterVarFiles: true,
		//VarFiles: [] string {"dev.tfvars" }, //an additional tfvars containing configuration files can be passed here
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

	//wait for 60 seconds to let resource acheive stable state
	time.Sleep(60* time.Second)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "host_vpc_name")

	fmt.Println(" ========= Verify Subnet Name ========= ")
	assert.Equal(t, network_name, output)

	fmt.Println(" ========= Verify Subnetwork Id ========= ")
	output = terraform.Output(t, terraformOptions, "host_subnetwork_id")
	cloudSqlInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	subnetworkId := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s",host_project_id,region,subnetwork_name)
	assert.Equal(t,subnetworkId , output)
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
	iteration = 0;
	// performs iterations for 3 times to check if the database gets created or not
	for {
		cmd = shell.Command{
			Command : "gcloud",
			Args : []string{"sql","databases","describe",test_dbname,"--instance="+cloudSqlInstanceName,"--project="+service_project_id,"--format=json"},
		}
		op,err = shell.RunCommandAndGetOutputE(t, cmd)
		if err == nil || iteration > 3 {
			break
		} else {
			fmt.Printf("Database with Database Name %s not found in cloud sql instance %s in project %s, will reattempt in few sec", test_dbname, cloudSqlInstanceName, service_project_id)
		}
		time.Sleep(60 * time.Second)
		iteration++;
	}

	if err != nil {
		t.Fatalf("Expected Database Name : %s at Cloudsql Instance :%s does not exists in Project : %s ", test_dbname, cloudSqlInstanceName, service_project_id)
	}
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result = gjson.Parse(op)
	if err != nil {
		fmt.Sprintf("=== Error %s Encountered while executing %s", err ,text)
	}
	assert.Equal(t, test_dbname, gjson.Get(result.String(),"name").String())
}
