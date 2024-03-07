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
	"fmt"
	"strconv"
	"log"
	"os"
	"testing"
	"time"
	"github.com/gruntwork-io/terratest/modules/shell"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/google/go-cmp/cmp"
	"github.com/tidwall/gjson"
)

const terraformDirectoryPath = "../../../../examples/ServiceNetworkingAcrossInterconnect"
const terraformPreReqDirPath = "./prereq"

var hostProjectID = os.Getenv("TF_VAR_host_project_id")
// Name of the deployed dedicated interconnect received after deploying the resource in the test lab
// e.g. dedicated-ix-vpn-client-0
var deployedInterconnectName = os.Getenv("deployed_interconnect_name")

// Variables for Interconnect configuration.
var interconnectProjectID = os.Getenv("TF_VAR_interconnect_project_id")

var deploymentNumber = 0
var databaseVersion = "MYSQL_8_0"
var region = "us-west2"
var zone = "us-west2-a"
var cloudSQLInstanceName = "cn-sqlinstance10-test"
var networkName = "cloudsql-easy"
var subnetworkName = "cloudsql-easy-subnet"
var subnetworkIPCidr = "10.0.0.0/24"
var deletionProtection = false
var testDBName = "test_db"

// Compute instance specific customisation.
var startupScript = "setupsql-interconnect.sh"
var gceTags = []string{"internet-egress"}

// Variables for Interconnect configuration.
var firstInterconnectName = "cso-lab-interconnect-1"
var secondInterconnectName = "cso-lab-interconnect-2"
var userSpecifiedIPRange = []string{"0.0.0.0/0", "199.36.154.8/30"}

// First vlan attachment configuration values.
var firstVaAsn = "65418"
var firstVlanAttachmentName = "vlan-attachment-a"
var firstVaBandwidth = "BPS_1G"

// Second vlan attachment configuration values.
var secondVaAsn = "65418"
var secondVlanAttachmentName = "vlan-attachment-b"
var secondVaBandwidth = "BPS_1G"

/* TestMySQLPrivateAndInterconnectVPN tests the creation of ServiceNetworkingAcrossInterconnect
example by creating a new vpc and a new subnet. */
func TestMySQLPrivateAndInterconnectVPN(t *testing.T) {

	deploymentNumber, err := strconv.Atoi(deployedInterconnectName[len(deployedInterconnectName)-1:])
	if err != nil {
		t.Errorf("Deployment number is not an int, using default value for deployment number.")
		deploymentNumber = 0;
	}
	var icRouterBgpAsn = 65000 + deploymentNumber;
	var firstVaBgpRange = fmt.Sprintf("169.254.6%d.0/29",deploymentNumber)
	var firstVlanTag = 600+deploymentNumber
	var secondVaBgpRange = fmt.Sprintf("169.254.6%d.8/29",deploymentNumber)
	var secondVlanTag = 600+deploymentNumber
	var tfVars = map[string]any{
		"host_project_id":          hostProjectID,
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
		"user_specified_ip_range":  userSpecifiedIPRange,
		"interconnect_project_id":  interconnectProjectID,
		"first_interconnect_name":  firstInterconnectName,
		"second_interconnect_name": secondInterconnectName,
		"first_va_name":            firstVlanAttachmentName,
		"ic_router_bgp_asn":        icRouterBgpAsn,
		"first_va_asn":             firstVaAsn,
		"first_va_bandwidth":       firstVaBandwidth,
		"first_va_bgp_range":       firstVaBgpRange,
		"first_vlan_tag":           firstVlanTag,
		"second_va_asn":            secondVaAsn,
		"second_va_name":           secondVlanAttachmentName,
		"second_va_bandwidth":      secondVaBandwidth,
		"second_va_bgp_range":      secondVaBgpRange,
		"second_vlan_tag":          secondVlanTag,
		"startup_script":           startupScript,
		"gce_tags":                 gceTags,
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir:         terraformDirectoryPath,
		Vars:                 tfVars,
		Reconfigure:          true,
		Lock:                 true,
		NoColor:              true,
		SetVarsAfterVarFiles: true,
	})
	initiateTestForNetworkResource(t,terraformOptions,firstVlanTag)
}

/* TestMySQLPrivateAndICVPNWithoutVPCCreation tests the creation of
ServiceNetworkingAcrossInterconnect example by creating a new vpc and a new subnet. */
func TestMySQLPrivateAndICVPNWithoutVPCCreation(t *testing.T) {
	deploymentNumber, err := strconv.Atoi(deployedInterconnectName[len(deployedInterconnectName)-1:])
	if err != nil {
		t.Errorf("Deployment number is not an int, using default value for deployment number.")
		deploymentNumber = 0;
	}
	var icRouterBgpAsn = 65000 + deploymentNumber;
	var firstVaBgpRange = fmt.Sprintf("169.254.6%d.0/29",deploymentNumber)
	var firstVlanTag = 600+deploymentNumber
	var secondVaBgpRange = fmt.Sprintf("169.254.6%d.8/29",deploymentNumber)
	var secondVlanTag = 600+deploymentNumber
	networkName = "hostcloudsql-easy"
	subnetworkName = "hostcloudsql-easy-subnet"
	var tfVars = map[string]any{
		"host_project_id":          hostProjectID,
		"database_version":         databaseVersion,
		"cloudsql_instance_name":   cloudSQLInstanceName,
		"region":                   region,
		"zone":                     zone,
		"create_network":           false,
		"create_subnetwork":        false,
		"network_name":             networkName,
		"subnetwork_name":          subnetworkName,
		"subnetwork_ip_cidr":       subnetworkIPCidr,
		"deletion_protection":      deletionProtection,
		"user_specified_ip_range":  userSpecifiedIPRange,
		"interconnect_project_id":  interconnectProjectID,
		"first_interconnect_name":  firstInterconnectName,
		"second_interconnect_name": secondInterconnectName,
		"first_va_name":            firstVlanAttachmentName,
		"ic_router_bgp_asn":        icRouterBgpAsn,
		"first_va_asn":             firstVaAsn,
		"first_va_bandwidth":       firstVaBandwidth,
		"first_va_bgp_range":       firstVaBgpRange,
		"first_vlan_tag":           firstVlanTag,
		"second_va_asn":            secondVaAsn,
		"second_va_name":           secondVlanAttachmentName,
		"second_va_bandwidth":      secondVaBandwidth,
		"second_va_bgp_range":      secondVaBgpRange,
		"second_vlan_tag":          secondVlanTag,
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir:         terraformDirectoryPath,
		Vars:                 tfVars,
		Reconfigure:          true,
		Lock:                 true,
		NoColor:              true,
		SetVarsAfterVarFiles: true,
	})
	// Create VPC and subnet outside of the terraform module

	text := "compute"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "create", networkName, "--project=" + hostProjectID, "--format=json", "--bgp-routing-mode=global", "--subnet-mode=custom","--verbosity=none"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		log.Printf("===Error %s Encountered while executing %s", err, text)
	}
	cmd = shell.Command{
		Command: "gcloud",
		Args:    []string{text, "networks", "subnets", "create", subnetworkName, "--network="+networkName ,"--project=" + hostProjectID,"--range=10.0.0.0/24","--region="+region ,"--format=json", "--enable-private-ip-google-access", "--enable-flow-logs","--verbosity=none"},
	}
	_, err = shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		log.Printf("===Error %s Encountered while executing %s", err, text)
	}
	initiateTestForNetworkResource(t,terraformOptions, firstVlanTag)
}

/* initiateTestForNetworkResource is a helper function that helps in verification
of the resources being created as part of test. */
func initiateTestForNetworkResource(t *testing.T, terraformOptions *terraform.Options, firstVlanTag int) {
	t.Helper()
	var password string
	var privateIPAddressSQL string
	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Wait for 60 seconds to let resource acheive stable state.
	time.Sleep(60 * time.Second)

	// Run `terraform output` to get the values of output variables and Check they have the expected values.
	output := terraform.Output(t, terraformOptions, "host_vpc_name")

	log.Println(" ========= Verify Subnet Name ========= ")
	want := networkName
	got := output
	if !cmp.Equal(got, want) {
		t.Errorf("Test Network Name = %v, want = %v", got, want)
	}

	log.Println(" ========= Verify Subnetwork Id ========= ")
	output = terraform.Output(t, terraformOptions, "host_subnetwork_id")
	cloudSQLInstanceName = terraform.Output(t, terraformOptions, "cloudsql_instance_name")
	subnetworkID := fmt.Sprintf("projects/%s/regions/%s/subnetworks/%s", hostProjectID, region, subnetworkName)
	want = subnetworkID
	got = output
	if !cmp.Equal(got, want) {
		t.Errorf("Test Sub Network Id = %v, want = %v", got, want)
	}

	// Validate if SQL instance wih private IP is up and running.
	text := "sql"
	cmd := shell.Command{
		Command: "gcloud",
		Args:    []string{text, "instances", "describe", cloudSQLInstanceName, "--project=" + hostProjectID, "--format=json", "--verbosity=none"},
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
	want = "false"
	got = gjson.Get(result.String(),"settings.ipConfiguration.ipv4Enabled").String()
	if !cmp.Equal(got, want) {
		t.Errorf("Test Public Ip configuration = %v, want = %v", got, want)
	}
	log.Println(" ========= Verify SQL RUNNING Instance State ========= ")
	want = "RUNNABLE"
	got = gjson.Get(result.String(),"state").String()
	if !cmp.Equal(got, want) {
			t.Errorf("Test SQL State = %v, want = %v", got, want)
	}

	// Validate if interconnects vlans attachments are up & running with Established Connection.
	log.Println(" ====================================================== ")
	log.Println(" ========= Verify Interconnect/VLAN Tunnels ========= ")

	var interconnectAttachmentNameList = []string{firstVlanAttachmentName, secondVlanAttachmentName}

	for _, vlanAttachmentName := range interconnectAttachmentNameList {
		cmd = shell.Command{
			Command: "gcloud",
			Args:    []string{"compute", "interconnects", "attachments", "describe", vlanAttachmentName, "--region", region, "--project", hostProjectID, "--format=json", "--verbosity=none"},
		}
		op, err = shell.RunCommandAndGetOutputE(t, cmd)
		if err != nil {
			log.Printf("===Error %s Encountered while executing %s", err, text)
		}
		if !gjson.Valid(op) {
			t.Fatalf("Error parsing output, invalid json: %s", op)
		}
		result = gjson.Parse(op)
		if err != nil {
			log.Printf("=== Error %s Encountered while executing %s", err, text)
		}
		log.Printf(" \n========= Validating attachment %s ============\n", vlanAttachmentName)
		log.Println(" ========= Check if attach Operation Status is active ========= ")
		want = "OS_ACTIVE"
		got = gjson.Get(result.String(), "operationalStatus").String()
		if !cmp.Equal(got, want) {
				t.Errorf("Test VLAN Operational State = %v, want = %v", got, want)
		}
		log.Println(" ========= Check if state is Active ========= ")
		want = "ACTIVE"
		got = gjson.Get(result.String(), "state").String()
		if !cmp.Equal(got, want) {
				t.Errorf("Test VLAN State = %v, want = %v", got, want)
		}
		log.Println(" ========= Check if type is Dedicated ========= ")
		want = "DEDICATED"
		got = gjson.Get(result.String(), "type").String()
		if !cmp.Equal(got, want) {
				t.Errorf("Test Interconnect type = %v, want = %v", got, want)
		}

		log.Println(" ========= Check if vlan tag is Same as Configured ========= ")
		want = strconv.Itoa(firstVlanTag)
		got = gjson.Get(result.String(), "vlanTag8021q").String()
		if !cmp.Equal(got, want) {
				t.Errorf("Test VLAN tag = %v, want = %v", got, want)
		}
		// Fetch the cloudsql instance password and private IP address to be used for created test database.

		password = terraform.Output(t, terraformOptions, "cloudsql_generated_password")
		privateIPAddressSQL = terraform.Output(t, terraformOptions, "cloudsql_private_ip")
	}
	initateTestForCloudSQLTestDB(t, cloudSQLInstanceName, password, privateIPAddressSQL)
}

/* initateTestForCloudSQLTestDB is a helper fucntion that helps in creation of database
from the onprem machine connected over interconnect to google vpc. */
func initateTestForCloudSQLTestDB(t *testing.T, cloudSQLInstanceName string, password string, privateIPAddressSQL string) {
	t.Helper()
	var tfVars = map[string]any{
		"host_project_id":  hostProjectID,
		"region":           region,
		"zone":             zone,
		"network_name":     networkName,
		"subnetwork_name":  subnetworkName,
		"test_dbname":      testDBName,
		"gce_tags":         gceTags,
		"default_password": password,
		"private_sql_ip":   privateIPAddressSQL,
	}

	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		// Set the path to the Terraform code that will be tested.
		TerraformDir:         terraformPreReqDirPath,
		Vars:                 tfVars,
		Reconfigure:          true,
		Lock:                 true,
		NoColor:              true,
		SetVarsAfterVarFiles: true,
	})

	// Clean up resources with "terraform destroy" at the end of the test.
	defer terraform.Destroy(t, terraformOptions)

	// Run "terraform init" and "terraform apply". Fail the test if there are any errors.
	terraform.InitAndApply(t, terraformOptions)

	// Wait for 60 seconds to let resource acheive stable state.
	time.Sleep(60 * time.Second)

	// Iterate through list of database to ensure a new db was created.
	log.Println(" ====================================================== ")
	log.Println(" ========= Verify DB Creation ========= ")
	iteration := 0
	// Performs iterations for 3 times to check if the database gets created or not.
	for {
		cmd := shell.Command{
			Command: "gcloud",
			Args:    []string{"sql", "databases", "describe", testDBName, "--instance=" + cloudSQLInstanceName, "--project=" + hostProjectID, "--format=json", "--verbosity=none"},
		}
		op, err := shell.RunCommandAndGetOutputE(t, cmd)
		if err == nil || iteration > 3 {
			if !gjson.Valid(op) {
				t.Fatalf("Error parsing output, invalid json: %s", op)
			}
			result := gjson.Parse(op)
			if err != nil {
				log.Printf("=== Error %s Encountered.", err)
			}
			want := testDBName
			got := gjson.Get(result.String(), "name").String()
			if !cmp.Equal(got, want) {
					t.Errorf("Test DB Name = %v, want = %v", got, want)
			}

			break
		} else {
			log.Printf("Database with Database Name %s not found in cloud sql instance %s in project %s, will reattempt in few sec", testDBName, cloudSQLInstanceName, hostProjectID)
		}
		time.Sleep(60 * time.Second)
		iteration++
	}
}
