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

package vpc_across_vpn_test

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
var host_project_id          = os.Getenv("TF_VAR_host_project_id");
var service_project_id       = os.Getenv("TF_VAR_service_project_id");
var user_project_id          = os.Getenv("TF_VAR_user_project_id");
var cloudsql_instance_name   = "cn-sqlinstance10-test";
var network_name             = "cloudsql-easy";
var subnetwork_name          = "cloudsql-easy-subnet";
var region                   = "us-central1";
var zone 										 = "us-central1-a";
var test_dbname              = "test_db";
var user_region              = "us-west1";
var user_zone                = "us-west1-a";
var uservpc_network_name     = "user-cloudsql-easy";
var uservpc_subnetwork_name  = "user-cloudsql-easy-subnet";
var database_version 					= "MYSQL_8_0"
var deletion_protection       = false;

// name the function as Test*
func TestMySqlPrivateAndVPNModule(t *testing.T) {
	//wait for 60 seconds to allow resources to be available
	time.Sleep(60 * time.Second)
	var iteration int;
	region                   = "us-central1";
	cloudsql_instance_name   = "cn-sqlinstance10-test";
	network_name             = "cloudsql-easy";
	subnetwork_name          = "cloudsql-easy-subnet";
	subnetwork_ip_cidr       := "10.2.0.0/16"
	uservpc_network_name       = "cloudsql-user"
	uservpc_subnetwork_name    = "cloudsql-user-subnet"
	uservpc_subnetwork_ip_cidr := "10.10.30.0/24"
	user_region                = "us-west1"
	user_zone                  = "us-west1-a"

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
		"subnetwork_name"            : subnetwork_name, // this subnetwork will be created
		"subnetwork_ip_cidr"         : subnetwork_ip_cidr,
    "user_project_id"            : user_project_id,
		"user_region"                : user_region,
		"user_zone"                  : user_zone,
		"create_user_vpc_network"    : true,
		"create_user_vpc_subnetwork" : true,
		"uservpc_network_name"       : uservpc_network_name,
		"uservpc_subnetwork_name"    : uservpc_subnetwork_name,
		"uservpc_subnetwork_ip_cidr" : uservpc_subnetwork_ip_cidr,
		"test_dbname"                : test_dbname,
		"deletion_protection" 	     : deletion_protection,

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
	assert.Equal(t, network_name, output)

	fmt.Println(" ========= Verify Subnetwork Id ========= ")
	output = terraform.Output(t, terraformOptions, "host_subnetwork_id")
	cloudSqlInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	subnetworkId := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s",host_project_id,region,subnetwork_name)
	assert.Equal(t,subnetworkId , output)

	// Validate if SQL instance wih private IP is up and running
	text := "sql"
	cmd := shell.Command{
		Command : "gcloud",
		Args : []string{text,"instances","describe",cloudSqlInstanceName,"--project="+service_project_id,"--format=json","--verbosity=none"},
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
	var projectId = ""
	for _, v := range vpnTunnelName {
		if v == "ha-vpn-1-remote-0" || v=="ha-vpn-1-remote-1" {
			projectId = host_project_id;
		} else {
			projectId = user_project_id;
		}
		cmd = shell.Command{
			Command : "gcloud",
			Args : []string{"compute","vpn-tunnels","describe",v,"--project",projectId,"--region",region,"--format=json","--verbosity=none"},
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
			Args : []string{"sql","databases","describe",test_dbname,"--instance="+cloudSqlInstanceName,"--project="+service_project_id,"--format=json","--verbosity=none"},
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
		fmt.Sprintf("======= Error %s Encountered while executing %s", err ,text)
	}
	assert.Equal(t, test_dbname, gjson.Get(result.String(),"name").String())
}

func TestUsingExistingNetworkMySqlPrivateAndVPNModule(t *testing.T) {
	//wait for 60 seconds to allow resources to be available
	time.Sleep(60 * time.Second)
	var iteration int;
	cloudsql_instance_name   = "cn-sqlinstance10-test";
	network_name             = "host-cloudsql-easy";
	subnetwork_name          = "host-cloudsql-easy-subnet";
	uservpc_network_name     = "user-cloudsql-easy";
	uservpc_subnetwork_name  = "user-cloudsql-easy-subnet";

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
		"subnetwork_name"            : subnetwork_name, // this subnetwork will be created
		"user_project_id"            : user_project_id,
		"user_region"                : user_region,
		"user_zone"                  : user_zone,
		"create_user_vpc_network"    : false,
		"create_user_vpc_subnetwork" : false,
		"uservpc_network_name"       : uservpc_network_name,
		"uservpc_subnetwork_name"    : uservpc_subnetwork_name,
		"test_dbname"                : test_dbname,
		"deletion_protection" 	     : deletion_protection,
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
		Args : []string{text,"networks","describe",network_name,"--project="+host_project_id,"--format=json","--verbosity=none"},
	}
	op,err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Fatalf("Expected Network : %s does not exists in Project : %s ", network_name, host_project_id)
	}

	//validate if the subnet already exists in host project
	cmd = shell.Command{
		Command : "gcloud",
		Args : []string{text,"networks","subnets","describe",subnetwork_name,"--project="+host_project_id,"--region="+region,"--format=json","--verbosity=none"},
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
	time.Sleep(60 * time.Second)


	// Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "host_vpc_name")

	fmt.Println(" ========= Verify Subnet Name ========= ")
	assert.Equal(t, network_name, output)

	fmt.Println(" ========= Verify Subnetwork Id ========= ")
	output = terraform.Output(t, terraformOptions, "host_subnetwork_id")
	cloudSqlInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	subnetworkId := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s",host_project_id,region,subnetwork_name)
	assert.Equal(t,subnetworkId , output)

	// Validate if SQL instance wih private IP is up and running
	text = "sql"
	cmd = shell.Command{
		Command : "gcloud",
		Args : []string{text,"instances","describe",cloudSqlInstanceName,"--project="+service_project_id,"--format=json","--verbosity=none"},
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
	var projectId = ""
	for _, v := range vpnTunnelName {
		if v == "ha-vpn-1-remote-0" || v=="ha-vpn-1-remote-1" {
			projectId = host_project_id;
		} else {
			projectId = user_project_id;
		}
		cmd = shell.Command{
			Command : "gcloud",
			Args : []string{"compute","vpn-tunnels","describe",v,"--project",projectId,"--region",region,"--format=json","--verbosity=none"},
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
			Args : []string{"sql","databases","describe",test_dbname,"--instance="+cloudSqlInstanceName,"--project="+service_project_id,"--format=json","--verbosity=none"},
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
		fmt.Sprintf("======= Error %s Encountered while executing %s", err ,text)
	}
	assert.Equal(t, test_dbname, gjson.Get(result.String(),"name").String())
}

