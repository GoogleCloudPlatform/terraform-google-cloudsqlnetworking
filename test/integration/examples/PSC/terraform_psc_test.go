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

package sqlpsctest

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

const terraformDirectoryPath = "../../../../examples/3.PSC"

// Setting required variables for the integration tests.

var consumerProjectID = os.Getenv("TF_VAR_consumer_project_id")
var producerProjectID = os.Getenv("TF_VAR_producer_project_id")
var cloudSQLInstanceName = "cloudsql-easy-networking"
var consumerNetworkName = "cloudsql-easy"
var subnetworkIPCidr = "10.0.0.0/16"
var consumerSubnetworkName = "cloudsql-easy-subnet"
var region = "us-central1"
var zone = "us-central1-a"
var testDbname = "test_db"
var databaseVersion = "MYSQL_8_0"
var deletionProtection = false
var reservedIps = []string{"10.0.0.5"}
var tfVars = map[string]any{
	"consumer_project_id":      consumerProjectID,
	"producer_project_id":      producerProjectID,
	"database_version":         databaseVersion,
	"cloudsql_instance_name":   cloudSQLInstanceName,
	"region":                   region,
	"zone":                     zone,
	"consumer_network_name":    consumerNetworkName,
	"consumer_subnetwork_name": consumerSubnetworkName,
	"consumer_cidr":            subnetworkIPCidr,
	"deletion_protection":      deletionProtection,
	"test_dbname":              testDbname,
	"reserved_ips":             reservedIps,
}

/*
TestMySQLPrivateModule creates all the resources including the vpc network and subnetwork along with other
resources like Cloud SQL, Compute VM instance, Private Service Connect networking etc.

It then validates if,

1. Cloud SQL public IP is disabled
2. Compute VM instance in consumer project is able to connect and perform operations on Cloud SQL instance in producer project
3. Validates the existence of the network and subnetwork
*/
func TestMySQLPrivateModule(t *testing.T) {
	// Wait for 60 seconds to allow resources to be available
	time.Sleep(60 * time.Second)
	var iteration int

	consumerNetworkName = "cloudsql-easy"
	subnetworkIPCidr = "10.0.0.0/16"
	consumerSubnetworkName = "cloudsql-easy-subnet"

	tfVars = map[string]any{
		"consumer_project_id":      consumerProjectID,
		"producer_project_id":      producerProjectID,
		"database_version":         databaseVersion,
		"cloudsql_instance_name":   cloudSQLInstanceName,
		"region":                   region,
		"zone":                     zone,
		"consumer_network_name":    consumerNetworkName,
		"consumer_subnetwork_name": consumerSubnetworkName,
		"consumer_cidr":            subnetworkIPCidr,
		"deletion_protection":      deletionProtection,
		"test_dbname":              testDbname,
		"reserved_ips":             reservedIps,
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

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Wait for 120 seconds to let resource acheive stable state
	time.Sleep(120 * time.Second)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "consumer_network_name")

	log.Println(" ========= Verify Network Name ========= ")
	assert.Equal(t, consumerNetworkName, output)
	log.Println(" ========= Verify Subnetwork ID ========= ")
	output = terraform.Output(t, terraformOptions, "consumer_subnetwork_id")
	cloudSQLInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	consumerSubnetworkID := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", consumerProjectID, region, consumerSubnetworkName)
	assert.Equal(t, consumerSubnetworkID, output)

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
	log.Println(" ========= Verify SQL RUNNING Instance State ========= ")
	assert.Equal(t, "RUNNABLE", gjson.Get(result.String(), "state").String())

	// Iterate through list of database to ensure a new db was created
	log.Println(" ====================================================== ")
	log.Println(" ========= Verify DB Creation ========= ")
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
		time.Sleep(180 * time.Second)
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
		log.Printf("=== Error %s Encountered while executing %s", err, text)
	}
	assert.Equal(t, testDbname, gjson.Get(result.String(), "name").String())
}

/*
TestUsingExistingNetworkMySQLPrivateModule consumes existing VPC and subnetwork from the producer project but includes all the other
resources like Cloud SQL, Compute VM Instance, Private Service Connect networking etc.

It then validates if,

1. Cloud SQL does not have a public IP address
2. Compute VM instance in consumer project is able to connect and perform operations on Cloud SQL instance in customer's producer project
3. Validates the existence of the network and subnetwork
*/
func TestUsingExistingNetworkMySQLPrivateModule(t *testing.T) {
	// Wait for 120 seconds to allow resources to be available
	time.Sleep(120 * time.Second)
	var iteration int

	consumerNetworkName := "cloudsql-easy-existing"
	subnetworkIPCidr := "10.0.0.0/16"
	consumerSubnetworkName := "cloudsql-easy-subnet-existing"

	tfVars = map[string]any{
		"consumer_project_id":      consumerProjectID,
		"producer_project_id":      producerProjectID,
		"database_version":         databaseVersion,
		"cloudsql_instance_name":   cloudSQLInstanceName,
		"region":                   region,
		"zone":                     zone,
		"consumer_network_name":    consumerNetworkName,
		"consumer_subnetwork_name": consumerSubnetworkName,
		"create_network":           false,
		"create_subnetwork":        false,
		"consumer_cidr":            subnetworkIPCidr,
		"deletion_protection":      deletionProtection,
		"test_dbname":              testDbname,
		"reserved_ips":             reservedIps,
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
		Args:    []string{text, "networks", "subnets", "create", consumerSubnetworkName, "--project=" + consumerProjectID, "--format=json", "--network=" + consumerNetworkName, "--range=" + subnetworkIPCidr, "--region=" + region, "--verbosity=none"},
	}
	op, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		log.Printf("===Error %s Encountered while executing %s", err, text)
	}

	// validate if the VPC already exists in Consumer project
	text = "compute"
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "describe", consumerNetworkName, "--project=" + consumerProjectID, "--format=json", "--verbosity=none"},
	}
	op, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Fatalf("Expected Network : %s does not exists in Project : %s ", consumerNetworkName, consumerProjectID)
	}

	// Validate if the subnet already exists in Consumer project
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

	// Wait for 120 seconds to let resource acheive stable state
	time.Sleep(120 * time.Second)

	// Run `terraform output` to get the values of output variables and check they have the expected values.
	output := terraform.Output(t, terraformOptions, "consumer_network_name")

	log.Println(" ========= Verify Subnet Name ========= ")
	assert.Equal(t, consumerNetworkName, output)

	log.Println(" ========= Verify Subnetwork Name ========= ")
	output = terraform.Output(t, terraformOptions, "consumer_subnetwork_id")

	cloudSQLInstanceName := terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	consumerSubnetworkID := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", consumerProjectID, region, consumerSubnetworkName)
	assert.Equal(t, consumerSubnetworkID, output)

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
	log.Println(" ========= Verify SQL RUNNING Instance State ========= ")
	assert.Equal(t, "RUNNABLE", gjson.Get(result.String(), "state").String())

	// Iterate through list of database to ensure a new db was created
	log.Println(" ====================================================== ")
	log.Println(" ========= Verify DB Creation ========= ")
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
		time.Sleep(120 * time.Second)
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
		log.Printf("=== Error %s Encountered while executing %s", err, text)
	}
	assert.Equal(t, testDbname, gjson.Get(result.String(), "name").String())
}
