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

package hostservicetest

import (
	"os"
	"fmt"
	"time"
	"testing"
	"github.com/tidwall/gjson"
	"github.com/stretchr/testify/assert"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/shell"
)

const terraformDirectoryPath  = "../../../../examples/1.Host-Service-Project";
var hostProjectID             = os.Getenv("TF_VAR_host_project_id");
var serviceProjectID          = os.Getenv("TF_VAR_service_project_id");
var cloudSQLInstanceName      = "cn-sqlinstance10-test";
var subnetworkName            = "cloudsql-easy-subnet";
var region                    = "us-central1";
var zone 											= "us-central1-a";
var testDbname                = "test_db"
var databaseVersion 					= "MYSQL_8_0"
var deletionProtection        = false

/*
This test creates all the resources including the vpc network and subnetwork along with other
resources like cloudsql, computeinstance, service networking etc.

It then validates if
1. cloud sql only contains the privateIP
2. compute instance in service project is able to connect and perform operation on cloudsql instance
3. validates the existence of the network and subnetwork
*/
func TestMySqlPrivateModule(t *testing.T) {
	//wait for 60 seconds to allow resources to be available
	time.Sleep(60 * time.Second)
	var iteration int;

	networkName             := "cloudsql-easy";
	subnetworkName          = "cloudsql-easy-subnet";
	subnetworkIPCidr        := "10.2.0.0/16";

	tfVars := map[string]interface{}{
		"host_project_id"            : hostProjectID,
		"service_project_id"         : serviceProjectID,
		"database_version"           : databaseVersion,
		"cloudsql_instance_name"     : cloudSQLInstanceName,
		"region"                     : region,
		"zone"                       : zone,
		"create_network"             : true,
		"create_subnetwork"          : true,
		"network_name"               : networkName,
		"subnetwork_name"            : subnetworkName,
		"subnetwork_ip_cidr"         : subnetworkIPCidr,
		"deletion_protection" 			 : deletionProtection,
		"test_dbname"                : testDbname,
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars : tfVars,
		TerraformDir  : terraformDirectoryPath,
		Reconfigure : true,
		Lock: true,
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
	assert.Equal(t, networkName, output)
	fmt.Println(" ========= Verify Subnetwork Id ========= ")
	output = terraform.Output(t, terraformOptions, "host_subnetwork_id")
	cloudSQLInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	subnetworkID := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s",hostProjectID,region,subnetworkName)
	assert.Equal(t,subnetworkID , output)

	text := "sql"
	cmd := shell.Command{
		Command : "gcloud",
		Args : []string{text,"instances","describe",cloudSQLInstanceName,"--project="+serviceProjectID,"--format=json","--verbosity=none"},
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
			Args : []string{"sql","databases","describe",testDbname,"--instance="+cloudSQLInstanceName,"--project="+serviceProjectID,"--format=json","--verbosity=none"},
		}
		op,err = shell.RunCommandAndGetOutputE(t, cmd)
		if err == nil || iteration > 3 {
			break
		} else {
			fmt.Printf("Database with Database Name %s not found in cloud sql instance %s in project %s, will reattempt in few sec", testDbname, cloudSQLInstanceName, serviceProjectID)
		}
		time.Sleep(60 * time.Second)
		iteration++;
	}

	if err != nil {
		t.Fatalf("Expected Database Name : %s at Cloudsql Instance :%s does not exists in Project : %s ", testDbname, cloudSQLInstanceName, serviceProjectID)
	}
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result = gjson.Parse(op)
	if err != nil {
		fmt.Sprintf("=== Error %s Encountered while executing %s", err ,text)
	}
	assert.Equal(t, testDbname, gjson.Get(result.String(),"name").String())
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
	//wait for 60 seconds to allow resources to be available
	time.Sleep(60 * time.Second)
	var iteration int;

	networkName             := "host-cloudsql-easy";
	subnetworkName          := "host-cloudsql-easy-subnet";

	tfVars := map[string]interface{}{
		"host_project_id"            : hostProjectID,
		"service_project_id"         : serviceProjectID,
		"database_version"           : databaseVersion,
		"cloudsql_instance_name"     : cloudSQLInstanceName,
		"region"                     : region,
		"zone"                       : zone,
		"create_network"             : false,
		"create_subnetwork"          : false,
		"network_name"               : networkName,
		"subnetwork_name"            : subnetworkName,
		"deletion_protection" 			 : deletionProtection,
		"test_dbname"                : testDbname,
	}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars : tfVars,
		TerraformDir: terraformDirectoryPath,
		Reconfigure : true,
		Lock: true,
		//PlanFilePath: "./plan",
		NoColor: true,
		SetVarsAfterVarFiles: true,
		//VarFiles: [] string {"dev.tfvars" }, //an additional tfvars containing configuration files can be passed here
	})

	//validate if the VPC already exists in host project
	text := "compute"
	cmd := shell.Command{
		Command : "gcloud",
		Args : []string{text,"networks","describe",networkName,"--project="+hostProjectID,"--format=json","--verbosity=none"},
	}
	op,err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Fatalf("Expected Network : %s does not exists in Project : %s ", networkName, hostProjectID)
	}

	//validate if the subnet already exists in host project
	cmd = shell.Command{
		Command : "gcloud",
		Args : []string{text,"networks","subnets","describe",subnetworkName,"--project="+hostProjectID,"--region="+region,"--format=json","--verbosity=none"},
	}
	op,err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Fatalf("Expected Sub network : %s does not exists in Project : %s ", subnetworkName, hostProjectID)
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
	assert.Equal(t, networkName, output)

	fmt.Println(" ========= Verify Subnetwork Id ========= ")
	output = terraform.Output(t, terraformOptions, "host_subnetwork_id")
	cloudSQLInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	subnetworkID := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s",hostProjectID,region,subnetworkName)
	assert.Equal(t,subnetworkID , output)
	text = "sql"
	cmd = shell.Command{
		Command : "gcloud",
		Args : []string{text,"instances","describe",cloudSQLInstanceName,"--project="+serviceProjectID,"--format=json","--verbosity=none"},
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
			Args : []string{"sql","databases","describe",testDbname,"--instance="+cloudSQLInstanceName,"--project="+serviceProjectID,"--format=json","--verbosity=none"},
		}
		op,err = shell.RunCommandAndGetOutputE(t, cmd)
		if err == nil || iteration > 3 {
			break
		} else {
			fmt.Printf("Database with Database Name %s not found in cloud sql instance %s in project %s, will reattempt in few sec", testDbname, cloudSQLInstanceName, serviceProjectID)
		}
		time.Sleep(60 * time.Second)
		iteration++;
	}

	if err != nil {
		t.Fatalf("Expected Database Name : %s at Cloudsql Instance :%s does not exists in Project : %s ", testDbname, cloudSQLInstanceName, serviceProjectID)
	}
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result = gjson.Parse(op)
	if err != nil {
		fmt.Sprintf("=== Error %s Encountered while executing %s", err ,text)
	}
	assert.Equal(t, testDbname, gjson.Get(result.String(),"name").String())
}
