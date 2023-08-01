package utility

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/shell"
)

func CheckNetworkExists(t testing.TestingT, networkName string,projectName string) bool {
	text := "compute"
	cmd := shell.Command{
		Command : "gcloud",
		Args : []string{text,"networks","describe",networkName,"--project="+projectName,"--format=json"},
	}
	_ ,err := shell.RunCommandAndGetOutputE(t, cmd)
	if err != nil {
		t.Printf("Expected Network : %s does not exists in Project : %s ", networkName, projectName)
		return false
	}
	return true
}

