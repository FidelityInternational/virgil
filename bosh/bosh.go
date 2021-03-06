package bosh

import (
	"github.com/FidelityInternational/virgil/utility"
	"github.com/cloudfoundry-community/gogobosh"
	"regexp"
)

// FindDeployment - takes deployments and a regex to return the first matching deployment name
func FindDeployment(deployments []gogobosh.Deployment, regex string) string {
	for _, deployment := range deployments {
		matched, _ := regexp.MatchString(regex, deployment.Name)
		if matched {
			return deployment.Name
		}
	}
	return ""
}

// FindVMs - takes an array of VMs and a regex to filter on, returning a new array of all matching vms
func FindVMs(deploymentVMs []gogobosh.VM, regex string) []gogobosh.VM {
	var matchedVMs []gogobosh.VM
	for _, deploymentVM := range deploymentVMs {
		matched, _ := regexp.MatchString(regex, deploymentVM.JobName)
		if matched {
			matchedVMs = append(matchedVMs, deploymentVM)
		}
	}
	return matchedVMs
}

// GetAllIPs - Returns an array unique IP addresses for the Deployment VMs
func GetAllIPs(deploymentVms []gogobosh.VM) []string {
	var ips []string
	for _, deployment := range deploymentVms {
		ips = append(ips, deployment.IPs...)
	}
	utility.RemoveDuplicates(&ips)
	return ips
}
