package main

import (
	"fmt"
	"github.com/FidelityInternational/virgil/utility"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func main() {
	var systemDomain, cfUser, cfPassword string
	app := cli.NewApp()
	app.Name = "virgil"
	app.Usage = "A CLI App to return a list of firewall rules based on Cloud Foundry Security Groups"
	app.UsageText = "virgil [options] output_file"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "cf-sys-domain, csd",
			Usage:       "Cloud Foundry System Domain",
			Destination: &systemDomain,
		},
		cli.StringFlag{
			Name:        "cf-user, cu",
			Usage:       "Cloud Foundry Admin User",
			Destination: &cfUser,
		},
		cli.StringFlag{
			Name:        "cf-password, cp",
			Usage:       "Cloud Foundry Admin Password",
			Destination: &cfPassword,
		},
	}
	app.Action = func(c *cli.Context) {
		if systemDomain == "" || cfUser == "" || cfPassword == "" || c.NArg() == 0 {
			fmt.Println("cf-sys-domain, cf-user, cf-password and output_file must all be set")
			os.Exit(1)
		}
		fmt.Println("args", c.Args().First())
		config := &cfclient.Config{
			ApiAddress: fmt.Sprintf("https://api.%s", systemDomain),
			Username:   cfUser,
			Password:   cfPassword,
		}
		client, err := cfclient.NewClient(config)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		allSecGroups, err := client.ListSecGroups()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		secGroups := utility.GetUsedSecGroups(allSecGroups)
		firewallRules := utility.GetFirewallRules(secGroups)
		yml, _ := yaml.Marshal(&firewallRules)
		ioutil.WriteFile(c.Args()[0], []byte(fmt.Sprintf("---\nschema_version: \"1\"\n%v", string(yml))), os.FileMode(0444))
	}
	app.Run(os.Args)
}
