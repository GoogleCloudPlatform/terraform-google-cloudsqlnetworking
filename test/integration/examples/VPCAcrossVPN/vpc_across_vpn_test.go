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
	"fmt"
	"time"
	"testing"
	"github.com/tidwall/gjson"
	"github.com/stretchr/testify/assert"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/shell"
)

const terraformDirectoryPath   = "../../../../examples/2.VPC-Across-VPN";
var hostProjectID              = os.Getenv("TF_VAR_host_project_id");
var serviceProjectID           = os.Getenv("TF_VAR_service_project_id");
var userProjectID              = os.Getenv("TF_VAR_user_project_id");
var cloudSQLInstanceName       = "cn-sqlinstance10-test";
var networkName                = "cloudsql-easy";
var subnetworkName             = "cloudsql-easy-subnet";
var region                     = "us-central1";
var zone 										   = "us-central1-a";
var testDbname 							   = "test_db";
var userRegion                 = "us-west1";
var userZone                   = "us-west1-a";
var uservpcNetworkName         = "user-cloudsql-easy";
var uservpcSubnetworkName      = "user-cloudsql-easy-subnet";
var databaseVersion 				   = "MYSQL_8_0"
var deletionProtection         = false;

// name the function as Test*
func TestMySqlPrivateAndVPNModule(t *testing.T) {
	//wait for 60 seconds to allow resources to be available
	time.Sleep(60 * time.Second)
	var iteration int;
	region                    = "us-central1";
	cloudSQLInstanceName      = "cn-sqlinstance10-test";
	networkName               = "cloudsql-easy";
	subnetworkName            = "cloudsql-easy-subnet";
	subnetworkIPCidr         := "10.2.0.0/16"
	uservpcNetworkName        = "cloudsql-user"
	uservpcSubnetworkName     = "cloudsql-user-subnet"
	uservpcSubnetworkIPCidr  := "10.10.30.0/24"
	userRegion                = "us-west1"
	userZone                  = "us-west1-a"

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
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars : tfVars,
		TerraformDir: terraformDirectoryPath,
		Reconfigure : true,
		NoColor: true,
		Lock: true,
		SetVarsAfterVarFiles: true,
		//VarFiles: [] string {"dev.tfvars" },
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	//wait for 60 seconds to let resource acheive stable state
	time.Sleep(60 * time.Second)


	// Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "host_vpc_name")

	fmt.Println(" ========= Verify Subnet Name ========= ")
	assert.Equal(t, networkName, output)

	fmt.Println(" ========= Verify Subnetwork Id ========= ")
	output = terraform.Output(t, terraformOptions, "host_subnetwork_id")
	cloudSQLInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	subnetworkID := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s",hostProjectID,region,subnetworkName)
	assert.Equal(t,subnetworkID , output)

	// Validate if SQL instance wih private IP is up and running
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

	// Validate if VPN tunnels are up & running with Established Connection
	fmt.Println(" ====================================================== ")
	fmt.Println(" ========= Verify VPN Tunnel ========= ")

	var vpnTunnelName = []string { "ha-vpn-1-remote-0","ha-vpn-1-remote-1","ha-vpn-2-remote-0","ha-vpn-2-remote-1"}
	var projectID = ""
	for _, v := range vpnTunnelName {
		if v == "ha-vpn-1-remote-0" || v=="ha-vpn-1-remote-1" {
			projectID = hostProjectID;
		} else {
			projectID = userProjectID;
		}
		cmd = shell.Command{
			Command : "gcloud",
			Args : []string{"compute","vpn-tunnels","describe",v,"--project",projectID,"--region",region,"--format=json","--verbosity=none"},
		}
		op,err = shell.RunCommandAndGetOutputE(t, cmd)
		if !gjson.Valid(op) {
			t.Fatalf("Error parsing output, invalid json: %s", op)
		}
		result = gjson.Parse(op)
		if err != nil {
			fmt.Sprintf("===Error %s Encountered while executing %s", err ,text)
		}
		fmt.Printf(" \n========= validating tunnel %s ============\n",v);
		fmt.Println(" ========= check if tunnel is up & running ========= ",)
		assert.Equal(t, "Tunnel is up and running.", gjson.Get(result.String(),"detailedStatus").String())
		fmt.Println(" ========= check if connection is established ========= ")
		assert.Equal(t, "ESTABLISHED", gjson.Get(result.String(),"status").String())
	}

	//Iterate through list of database to ensure a new db was created
	fmt.Println(" ====================================================== ")
	fmt.Println(" =========== Verify DB Creation =========== ")
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
		fmt.Sprintf("======= Error %s Encountered while executing %s", err ,text)
	}
	assert.Equal(t, testDbname, gjson.Get(result.String(),"name").String())
}

func TestUsingExistingNetworkMySqlPrivateAndVPNModule(t *testing.T) {
	//wait for 60 seconds to allow resources to be available
	time.Sleep(60 * time.Second)
	var iteration int;
	cloudSQLInstanceName    = "cn-sqlinstance10-test";
	networkName             = "host-cloudsql-easy";
	subnetworkName          = "host-cloudsql-easy-subnet";
	uservpcNetworkName      = "user-cloudsql-easy";
	uservpcSubnetworkName   = "user-cloudsql-easy-subnet";

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
		"subnetwork_name"            : subnetworkName, // this subnetwork will be created
		"user_project_id"            : userProjectID,
		"user_region"                : userRegion,
		"user_zone"                  : userZone,
		"create_user_vpc_network"    : false,
		"create_user_vpc_subnetwork" : false,
		"uservpc_network_name"       : uservpcNetworkName,
		"uservpc_subnetwork_name"    : uservpcSubnetworkName,
		"test_dbname"                : testDbname,
		"deletion_protection" 	     : deletionProtection,
	}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars : tfVars,
		TerraformDir: terraformDirectoryPath,
		Reconfigure : true,
		Lock: true,
		NoColor: true,
		SetVarsAfterVarFiles: true,
		//VarFiles: [] string {"dev.tfvars" },
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
	time.Sleep(60 * time.Second)


	// Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "host_vpc_name")

	fmt.Println(" ========= Verify Subnet Name ========= ")
	assert.Equal(t, networkName, output)

	fmt.Println(" ========= Verify Subnetwork Id ========= ")
	output = terraform.Output(t, terraformOptions, "host_subnetwork_id")
	cloudSQLInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	subnetworkID := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s",hostProjectID,region,subnetworkName)
	assert.Equal(t,subnetworkID , output)

	// Validate if SQL instance wih private IP is up and running
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

	// Validate if VPN tunnels are up & running with Established Connection
	fmt.Println(" ====================================================== ")
	fmt.Println(" ========= Verify VPN Tunnel ========= ")

	var vpnTunnelName = []string { "ha-vpn-1-remote-0","ha-vpn-1-remote-1","ha-vpn-2-remote-0","ha-vpn-2-remote-1"}
	var projectID = ""
	for _, v := range vpnTunnelName {
		if v == "ha-vpn-1-remote-0" || v=="ha-vpn-1-remote-1" {
			projectID = hostProjectID;
		} else {
			projectID = userProjectID;
		}
		cmd = shell.Command{
			Command : "gcloud",
			Args : []string{"compute","vpn-tunnels","describe",v,"--project",projectID,"--region",region,"--format=json","--verbosity=none"},
		}
		op,err = shell.RunCommandAndGetOutputE(t, cmd)
		if !gjson.Valid(op) {
			t.Fatalf("Error parsing output, invalid json: %s", op)
		}
		result = gjson.Parse(op)
		if err != nil {
			fmt.Sprintf("=== Error %s Encountered while executing %s", err ,text)
		}
		fmt.Printf(" \n========= validating tunnel %s ============\n",v);
		fmt.Println(" ========= check if tunnel is up & running ========= ",)
		assert.Equal(t, "Tunnel is up and running.", gjson.Get(result.String(),"detailedStatus").String())
		fmt.Println(" ========= check if connection is established ========= ")
		assert.Equal(t, "ESTABLISHED", gjson.Get(result.String(),"status").String())
	}

	//Iterate through list of database to ensure a new db was created
	fmt.Println(" ====================================================== ")
	fmt.Println(" =========== Verify DB Creation =========== ")
	iteration = 0;
	//performs iterations for 3 times to check if the database gets created or not
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
		fmt.Sprintf("======= Error %s Encountered while executing %s", err ,text)
	}
	assert.Equal(t, testDbname, gjson.Get(result.String(),"name").String())
}

