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

package vpnacrosspsctest

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

const terraformDirectoryPath = "../../../../examples/4.PSC-Across-VPN"

// Setting required variables for the integration tests.

var consumerProjectID = os.Getenv("TF_VAR_consumer_project_id")
var producerProjectID = os.Getenv("TF_VAR_producer_project_id")
var userProjectID = os.Getenv("TF_VAR_user_project_id")
var cloudSQLInstanceName = "cn-sqlinstance-test"
var consumerNetworkName = "cloudsql-easy-int-test"
var consumerSubnetworkName = "cloudsql-easy-subnet-int-test"
var region = "us-central1"
var zone = "us-central1-a"
var user_region = "us-central1"
var user_zone = "us-central1-a"
var testDbname = "test_db"
var uservpcNetworkName = "user-cloudsql-easy-int-test"
var uservpcSubnetworkName = "user-cloudsql-easy-subnet-int-test"
var databaseVersion = "MYSQL_8_0"
var deletionProtection = false
var natName = "cloudsql-nat"
var routerName = "cloudsql-router"
var endpointIP = []string{"10.0.0.5"}

/*
TestMySQLPrivateAndVPNModule creates all the resources including the VPC network and subnetwork along with other
resources like Cloud SQL, Compute VM instance, Private Service Connect networking, VPN etc.

It then validates if,

1. Cloud SQL public IP is disabled
2. Compute VM instance in consumer project is able to connect and perform operations on Cloud SQL instance in producer project
3. Validates the existence of the network, subnetwork & VPN connectivity
*/
func TestMySQLPrivateAndVPNModule(t *testing.T) {
	// Wait for 60 seconds to allow resources to be available
	time.Sleep(60 * time.Second)
	var iteration int
	region = "us-central1"
	user_region = "us-central1"
	zone = "us-central1-a"
	user_zone = "us-central1-a"
	cloudSQLInstanceName = "cn-sqlinstance-vpn-test"
	consumerNetworkName = "cloudsql-easy"
	consumerSubnetworkName = "cloudsql-easy-subnet"
	uservpcNetworkName = "cloudsql-user"
	uservpcSubnetworkName = "cloudsql-user-subnet"
	uservpcSubnetworkIPCidr := "192.168.0.0/24"
	consumerSubnetworkIPCIDR := "10.0.0.0/16"

	tfVars := map[string]any{

		"consumer_project_id":        consumerProjectID,
		"producer_project_id":        producerProjectID,
		"user_project_id":            userProjectID,
		"database_version":           databaseVersion,
		"cloudsql_instance_name":     cloudSQLInstanceName,
		"region":                     region,
		"zone":                       zone,
		"user_region":                user_region,
		"user_zone":                  user_zone,
		"create_network":             true,
		"create_subnetwork":          true,
		"consumer_network_name":      consumerNetworkName,
		"consumer_subnetwork_name":   consumerSubnetworkName,
		"user_network_name":          uservpcNetworkName,
		"user_subnetwork_name":       uservpcSubnetworkName,
		"consumer_cidr":              consumerSubnetworkIPCIDR,
		"user_cidr":                  uservpcSubnetworkIPCidr,
		"create_user_vpc_network":    true,
		"create_user_vpc_subnetwork": true,
		"test_dbname":                testDbname,
		"deletion_protection":        deletionProtection,
		"nat_name":                   natName,
		"router_name":                routerName,
		"endpoint_ip":                endpointIP,
	}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars:                 tfVars,
		TerraformDir:         terraformDirectoryPath,
		Reconfigure:          true,
		NoColor:              true,
		Lock:                 true,
		SetVarsAfterVarFiles: true,
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Wait for 60 seconds to let resource achieve stable state
	time.Sleep(60 * time.Second)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "consumer_network_name")

	log.Println(" ========= Verify Subnet Name ========= ")
	assert.Equal(t, consumerNetworkName, output)

	log.Println(" ========= Verify Subnetwork Id ========= ")
	output = terraform.Output(t, terraformOptions, "consumer_subnetwork_id")
	cloudSQLInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	subnetworkID := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", consumerProjectID, region, consumerSubnetworkName)
	assert.Equal(t, subnetworkID, output)

	// Validate if SQL instance wih private IP is up and running
	text := "sql"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "instances", "describe", cloudSQLInstanceName, "--project=" + producerProjectID, "--format=json", "--verbosity=none"},
	}
	op, err := shell.RunCommandAndGetOutputE(t, cmd)
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result := gjson.Parse(op)
	if err != nil {
		log.Printf("===Error %s Encountered while executing %s", err, text)
	}
	log.Println(" ========= Verify if public IP is disabled ========= ")
	assert.Equal(t, false, gjson.Get(result.String(), "settings.ipConfiguration.ipv4Enabled").Bool())
	log.Println(" ========= Verify if PSC is enabled ========= ")
	assert.Equal(t, true, gjson.Get(result.String(), "settings.ipConfiguration.pscConfig.pscEnabled").Bool())
	log.Println(" ========= Verify if correct PSC project is allowlisted ========= ")
	assert.Equal(t, consumerProjectID, gjson.Get(result.String(), "settings.ipConfiguration.pscConfig.allowedConsumerProjects.0").String())

	log.Println(" ========= Verify SQL RUNNING Instance State ========= ")
	assert.Equal(t, "RUNNABLE", gjson.Get(result.String(), "state").String())

	// Validate if VPN tunnels are up & running with Established Connection
	log.Println(" ====================================================== ")
	log.Println(" ========= Verify VPN Tunnel ========= ")

	var vpnTunnelName = []string{"gcp-vpc-gateway1-remote-0", "gcp-vpc-gateway1-remote-1", "gcp-vpc-gateway2-remote-0", "gcp-vpc-gateway2-remote-1"}
	var projectID = ""
	for _, v := range vpnTunnelName {
		if v == "gcp-vpc-gateway1-remote-0" || v == "gcp-vpc-gateway1-remote-1" {
			projectID = consumerProjectID
			cmd = shell.Command{
				Command: "gcloud",
				Args:    []string{"compute", "vpn-tunnels", "describe", v, "--project", projectID, "--region", region, "--format=json", "--verbosity=none"},
			}
		} else {
			projectID = userProjectID
			cmd = shell.Command{
				Command: "gcloud",
				Args:    []string{"compute", "vpn-tunnels", "describe", v, "--project", projectID, "--region", region, "--format=json", "--verbosity=none"},
			}
		}
		op, err = shell.RunCommandAndGetOutputE(t, cmd)
		if !gjson.Valid(op) {
			t.Fatalf("Error parsing output, invalid json: %s", op)
		}
		result = gjson.Parse(op)
		if err != nil {
			log.Printf("===Error %s Encountered while executing %s", err, text)
		}
		log.Printf(" \n========= validating tunnel %s ============\n", v)
		log.Println(" ========= check if tunnel is up & running ========= ")
		assert.Equal(t, "Tunnel is up and running.", gjson.Get(result.String(), "detailedStatus").String())
		log.Println(" ========= check if connection is established ========= ")
		assert.Equal(t, "ESTABLISHED", gjson.Get(result.String(), "status").String())
	}

	//Iterate through list of database to ensure a new db was created
	log.Println(" ====================================================== ")
	log.Println(" =========== Verify DB Creation =========== ")
	iteration = 0
	// performs iterations for 3 times to check if the database gets created or not
	for {
		cmd = shell.Command{
			Command: "gcloud",
			Args:    []string{"sql", "databases", "describe", testDbname, "--instance=" + cloudSQLInstanceName, "--project=" + producerProjectID, "--format=json", "--verbosity=none"},
		}
		op, err = shell.RunCommandAndGetOutputE(t, cmd)
		if err == nil || iteration > 3 {
			break
		} else {
			log.Printf("Database with Database Name %s not found in cloud sql instance %s in project %s, will reattempt in few sec", testDbname, cloudSQLInstanceName, producerProjectID)
		}
		time.Sleep(60 * time.Second)
		iteration++
	}
	if err != nil {
		t.Fatalf("Expected Database Name : %s at Cloudsql Instance :%s does not exists in Project : %s ", testDbname, cloudSQLInstanceName, producerProjectID)
	}
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result = gjson.Parse(op)
	if err != nil {
		log.Printf("======= Error %s Encountered while executing %s", err, text)
	}
	assert.Equal(t, testDbname, gjson.Get(result.String(), "name").String())
}

/*
TestUsingExistingNetworkMySQLPrivateAndVPNModule consumes existing VPC, subnetwork & VPN connection from the producer project but includes all the other
resources like Cloud SQL, Compute VM Instance, Private Service Connect networking etc.

It then validates if,

1. Cloud SQL does not have a public IP address
2. Compute VM instance in consumer project is able to connect and perform operations on Cloud SQL instance in customer's producer project
3. Validates the existence of the network and subnetwork
*/
func TestUsingExistingNetworkMySQLPrivateAndVPNModule(t *testing.T) {
	// Wait for 60 seconds to allow resources to be available
	time.Sleep(60 * time.Second)
	var iteration int
	cloudSQLInstanceName = "cn-sqlinstance-vpn-test-existing"
	consumerNetworkName = "cloudsql-easy-existing"
	consumerSubnetworkName = "cloudsql-easy-subnet-existing"
	uservpcNetworkName = "user-cloudsql-easy-existing"
	uservpcSubnetworkName = "user-cloudsql-easy-subnet-existing"
	uservpcSubnetworkIPCidr := "192.168.0.0/24"
	consumerSubnetworkIPCIDR := "10.0.0.0/16"

	tfVars := map[string]any{
		"consumer_project_id":        consumerProjectID,
		"producer_project_id":        producerProjectID,
		"user_project_id":            userProjectID,
		"database_version":           databaseVersion,
		"cloudsql_instance_name":     cloudSQLInstanceName,
		"region":                     region,
		"zone":                       zone,
		"user_region":                user_region,
		"user_zone":                  user_zone,
		"create_network":             false,
		"create_subnetwork":          false,
		"consumer_network_name":      consumerNetworkName,
		"consumer_subnetwork_name":   consumerSubnetworkName,
		"create_user_vpc_network":    false,
		"create_user_vpc_subnetwork": false,
		"user_network_name":          uservpcNetworkName,
		"user_subnetwork_name":       uservpcSubnetworkName,
		"test_dbname":                testDbname,
		"deletion_protection":        deletionProtection,
		"nat_name":                   natName,
		"router_name":                routerName,
		"endpoint_ip":                endpointIP,
		"consumer_cidr":              consumerSubnetworkIPCIDR,
		"user_cidr":                  uservpcSubnetworkIPCidr,
	}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		Vars:                 tfVars,
		TerraformDir:         terraformDirectoryPath,
		Reconfigure:          true,
		Lock:                 true,
		NoColor:              true,
		SetVarsAfterVarFiles: true,
	})

	// Create VPC for Consumer

	text := "compute"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "create", consumerNetworkName, "--project=" + consumerProjectID, "--format=json", "--subnet-mode=custom", "--verbosity=none"},
	}
	op, err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		log.Printf("===Error %s Encountered while executing %s", err, text)
	}

	// Create subnet for Consumer

	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "subnets", "create", consumerSubnetworkName, "--project=" + consumerProjectID, "--format=json", "--network=" + consumerNetworkName, "--range=" + consumerSubnetworkIPCIDR, "--region=" + region, "--verbosity=none"},
	}
	op, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		log.Printf("===Error %s Encountered while executing %s", err, text)
	}

	// Create VPC for user

	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "create", uservpcNetworkName, "--project=" + userProjectID, "--format=json", "--subnet-mode=custom", "--verbosity=none"},
	}
	op, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		log.Printf("===Error %s Encountered while executing %s", err, text)
	}

	// Create subnets for user

	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "subnets", "create", uservpcSubnetworkName, "--project=" + userProjectID, "--format=json", "--network=" + uservpcNetworkName, "--range=" + uservpcSubnetworkIPCidr, "--region=" + region, "--verbosity=none"},
	}
	op, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		log.Printf("===Error %s Encountered while executing %s", err, text)
	}

	//validate if the VPC already exists in Consumer project
	text = "compute"
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "describe", consumerNetworkName, "--project=" + consumerProjectID, "--format=json", "--verbosity=none"},
	}
	op, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Fatalf("Expected Network : %s does not exists in Project : %s ", consumerNetworkName, consumerProjectID)
	}

	//validate if the subnet already exists in Consumer project
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "subnets", "describe", consumerSubnetworkName, "--project=" + consumerProjectID, "--region=" + region, "--format=json", "--verbosity=none"},
	}
	op, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Fatalf("Expected Sub network : %s does not exists in Project : %s ", consumerSubnetworkName, consumerProjectID)
	}

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Wait for 60 seconds to let resource achieve stable state
	time.Sleep(60 * time.Second)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "consumer_network_name")

	log.Println(" ========= Verify Network Name ========= ")
	assert.Equal(t, consumerNetworkName, output)

	log.Println(" ========= Verify Subnetwork Id ========= ")
	output = terraform.Output(t, terraformOptions, "consumer_subnetwork_id")
	cloudSQLInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	subnetworkID := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", consumerProjectID, region, consumerSubnetworkName)
	assert.Equal(t, subnetworkID, output)

	// Validate if SQL instance wih private IP is up and running
	text = "sql"
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "instances", "describe", cloudSQLInstanceName, "--project=" + producerProjectID, "--format=json", "--verbosity=none"},
	}
	op, err = shell.RunCommandAndGetOutputE(t, cmd)
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result := gjson.Parse(op)
	if err != nil {
		log.Printf("===Error %s Encountered while executing %s", err, text)
	}
	log.Println(" ========= Verify if public IP is disabled ========= ")
	assert.Equal(t, false, gjson.Get(result.String(), "settings.ipConfiguration.ipv4Enabled").Bool())
	log.Println(" ========= Verify if PSC is enabled ========= ")
	assert.Equal(t, true, gjson.Get(result.String(), "settings.ipConfiguration.pscConfig.pscEnabled").Bool())
	log.Println(" ========= Verify if correct PSC project is allowlisted ========= ")
	assert.Equal(t, consumerProjectID, gjson.Get(result.String(), "settings.ipConfiguration.pscConfig.allowedConsumerProjects.0").String())

	log.Println(" ========= Verify SQL RUNNING Instance State ========= ")
	assert.Equal(t, "RUNNABLE", gjson.Get(result.String(), "state").String())

	// Validate if VPN tunnels are up & running with Established Connection
	log.Println(" ====================================================== ")
	log.Println(" ========= Verify VPN Tunnel ========= ")

	var vpnTunnelName = []string{"gcp-vpc-gateway1-remote-0", "gcp-vpc-gateway1-remote-1", "gcp-vpc-gateway2-remote-0", "gcp-vpc-gateway2-remote-1"}
	var projectID = ""
	for _, v := range vpnTunnelName {
		if v == "gcp-vpc-gateway1-remote-0" || v == "gcp-vpc-gateway1-remote-1" {
			projectID = consumerProjectID
			text = "compute"
			cmd = shell.Command{
				Command: "gcloud",
				Args:    []string{text, "vpn-tunnels", "describe", v, "--project", projectID, "--region", region, "--format=json"},
			}
		} else {
			projectID = userProjectID
			text = "compute"
			cmd = shell.Command{
				Command: "gcloud",
				Args:    []string{text, "vpn-tunnels", "describe", v, "--project", projectID, "--region", user_region, "--format=json"},
			}
		}
		op, err = shell.RunCommandAndGetOutputE(t, cmd)
		if !gjson.Valid(op) {
			t.Fatalf("Error parsing output, invalid json: %s", op)
		}
		result = gjson.Parse(op)
		if err != nil {
			log.Printf("=== Error %s Encountered while executing %s", err, text)
		}
		log.Printf(" \n========= validating tunnel %s ============\n", v)
		log.Println(" ========= check if tunnel is up & running ========= ")
		assert.Equal(t, "Tunnel is up and running.", gjson.Get(result.String(), "detailedStatus").String())
		log.Println(" ========= check if connection is established ========= ")
		assert.Equal(t, "ESTABLISHED", gjson.Get(result.String(), "status").String())
	}

	//Iterate through list of database to ensure a new db was created
	log.Println(" ====================================================== ")
	log.Println(" =========== Verify DB Creation =========== ")
	iteration = 0
	//performs iterations for 3 times to check if the database gets created or not
	for {
		cmd = shell.Command{
			Command: "gcloud",
			Args:    []string{"sql", "databases", "describe", testDbname, "--instance=" + cloudSQLInstanceName, "--project=" + producerProjectID, "--format=json", "--verbosity=none"},
		}
		op, err = shell.RunCommandAndGetOutputE(t, cmd)
		if err == nil || iteration > 3 {
			break
		} else {
			log.Printf("Database with Database Name %s not found in cloud sql instance %s in project %s, will reattempt in few sec", testDbname, cloudSQLInstanceName, producerProjectID)
		}
		time.Sleep(60 * time.Second)
		iteration++
	}
	if err != nil {
		t.Fatalf("Expected Database Name : %s at Cloudsql Instance :%s does not exists in Project : %s ", testDbname, cloudSQLInstanceName, producerProjectID)
	}
	if !gjson.Valid(op) {
		t.Fatalf("Error parsing output, invalid json: %s", op)
	}
	result = gjson.Parse(op)
	if err != nil {
		log.Printf("======= Error %s Encountered while executing %s", err, text)
	}
	assert.Equal(t, testDbname, gjson.Get(result.String(), "name").String())
}
