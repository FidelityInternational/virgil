package main

import (
	"fmt"
	"github.com/FidelityInternational/virgil/utility"
	"github.com/cloudfoundry-community/go-cfclient"
	"os"
)

func main() {
	var (
		systemDomain  = os.Getenv("CF_SYSTEM_DOMAIN")
		firewallRules []utility.FirewallRule
	)
	c := &cfclient.Config{
		ApiAddress:   fmt.Sprintf("https://api.%s", systemDomain),
		LoginAddress: fmt.Sprintf("https://login.%s", systemDomain),
		Username:     os.Getenv("CF_USERNAME"),
		Password:     os.Getenv("CF_PASSWORD"),
	}
	client := cfclient.NewClient(c)
	allSecGroups := client.ListSecGroups()
	secGroups := utility.GetUsedSecGroups(allSecGroups)
	for _, secGroup := range secGroups {
		for _, secGroupRule := range secGroup.Rules {
			firewallRules, _ = utility.ProcessRule(secGroupRule, firewallRules)
		}
	}
	for _, rule := range firewallRules {
		fmt.Println(rule.Port)
		fmt.Println(rule.Protocol)
		fmt.Println(rule.Destination)
	}
}
