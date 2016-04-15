package utility

import (
	"fmt"
	"github.com/cloudfoundry-community/go-cfclient"
	"strconv"
	"strings"
)

// FirewallRules - A collection of Firewall Rules with version
type FirewallRules struct {
	SchemaVersion string         `yaml:"schema_version"`
	FirewallRules []FirewallRule `yaml:"firewall_rules"`
}

// FirewallRule struct
type FirewallRule struct {
	Port        string
	Destination []string
	Protocol    string
	Source      []string
}

// PortExpand - serperates port string into array, for example 2,5-7 becomes {2 5 6 7}
func PortExpand(portString string) ([]string, error) {
	var ports []string
	portsBefore := strings.Split(portString, ",")
	for _, port := range portsBefore {
		port = strings.TrimSpace(port)
		if strings.Contains(port, "-") {
			startFinish := strings.Split(port, "-")
			startString := strings.TrimSpace(startFinish[0])
			start, err := strconv.Atoi(startString)
			if err != nil || start <= 0 || start >= 65536 {
				return []string{}, fmt.Errorf("Port %s was invalid as part of range %s", startString, port)
			}
			endString := strings.TrimSpace(startFinish[1])
			end, err := strconv.Atoi(endString)
			if err != nil || end <= 0 || end >= 65536 {
				return []string{}, fmt.Errorf("Port %s was invalid as part of range %s", endString, port)
			}
			if len(startFinish) != 2 || start >= end {
				return []string{}, fmt.Errorf("Port range %s was invalid", port)
			}
			for i := start; i <= end; i++ {
				ports = append(ports, strconv.Itoa(i))
			}
		} else {
			portInt, err := strconv.Atoi(port)
			if err != nil || portInt <= 0 || portInt >= 65536 {
				return []string{}, fmt.Errorf("Port %s was invalid", port)
			}
			ports = append(ports, port)
		}
	}
	return ports, nil
}

// ProcessRule - returns a concise list of firewall rules for one security group rule
func ProcessRule(secGroupRule cfclient.SecGroupRule, firewallRules []FirewallRule, source []string) ([]FirewallRule, error) {
	ports, err := PortExpand(secGroupRule.Ports)
	if err != nil {
		return []FirewallRule{}, err
	}
	for _, port := range ports {
		var newRule = true
		for i, rule := range firewallRules {
			if rule.Port == port && rule.Protocol == secGroupRule.Protocol {
				rule.Destination = append(rule.Destination, secGroupRule.Destination)
				RemoveDuplicates(&rule.Destination)
				firewallRules[i] = rule
				newRule = false
			}
		}
		if newRule {
			newRules := FirewallRule{
				Port:        port,
				Protocol:    secGroupRule.Protocol,
				Destination: []string{secGroupRule.Destination},
				Source:      source,
			}
			firewallRules = append(firewallRules, newRules)
		}
	}
	return firewallRules, nil
}

// RemoveDuplicates - removes duplicated from array of strings
func RemoveDuplicates(xs *[]string) {
	found := make(map[string]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}

// GetUsedSecGroups - Trims out any security-groups that cannot be used. I.E not running, staging or bound
func GetUsedSecGroups(allSecGroups []cfclient.SecGroup) []cfclient.SecGroup {
	var secGroups []cfclient.SecGroup
	for _, secGroup := range allSecGroups {
		if secGroup.Running || secGroup.Staging || len(secGroup.SpacesData) != 0 {
			secGroups = append(secGroups, secGroup)
		}
	}
	return secGroups
}

// GetFirewallRules - Returns a concise list of firewall rules for all security groups
func GetFirewallRules(source []string, secGroups []cfclient.SecGroup) FirewallRules {
	var firewallRules FirewallRules
	firewallRules.SchemaVersion = "1"
	for _, secGroup := range secGroups {
		for _, secGroupRule := range secGroup.Rules {
			firewallRules.FirewallRules, _ = ProcessRule(secGroupRule, firewallRules.FirewallRules, source)
		}
	}
	return firewallRules
}
