package main

import (
	"fmt"
	"github.com/FidelityInternational/virgil/Godeps/_workspace/src/github.com/cloudfoundry-community/go-cfclient"
	"github.com/FidelityInternational/virgil/Godeps/_workspace/src/github.com/codegangsta/cli"
	"github.com/FidelityInternational/virgil/Godeps/_workspace/src/gopkg.in/yaml.v2"
	"github.com/FidelityInternational/virgil/bosh"
	"github.com/FidelityInternational/virgil/utility"
	"io/ioutil"
	"os"
)

func main() {
	var (
		systemDomain, cfUser, cfPassword, boshUser, boshPassword, boshURI, boshPort string
		skipSSLValidation                                                           = false
	)

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
		cli.StringFlag{
			Name:        "bosh-user, bu",
			Usage:       "BOSH User",
			Destination: &boshUser,
		},
		cli.StringFlag{
			Name:        "bosh-password, bp",
			Usage:       "BOSH Password",
			Destination: &boshPassword,
		},
		cli.StringFlag{
			Name:        "bosh-uri, buri",
			Usage:       "BOSH URI",
			Destination: &boshURI,
		},
		cli.StringFlag{
			Name:        "bosh-port, bport",
			Usage:       "BOSH Port",
			Value:       "25555",
			Destination: &boshPort,
		},
		cli.BoolFlag{
			Name:        "skip-ssl-validation, skip-ssl",
			Usage:       "Skip SSL Validation",
			Destination: &skipSSLValidation,
		},
	}
	app.Action = func(c *cli.Context) {
		if systemDomain == "" || cfUser == "" || cfPassword == "" || c.NArg() == 0 {
			fmt.Println("cf-sys-domain, cf-user, cf-password and output_file must all be set")
			os.Exit(1)
		}
		config := &cfclient.Config{
			ApiAddress:        fmt.Sprintf("https://api.%s", systemDomain),
			Username:          cfUser,
			Password:          cfPassword,
			SkipSslValidation: skipSSLValidation,
		}
		client, err := cfclient.NewClient(config)
		if boshUser == "" || boshPassword == "" || boshURI == "" || boshPort == "" {
			fmt.Println("BOSH user, password, URI and Port must all be set")
			os.Exit(1)
		}
		boshConfig := &bosh.Config{
			Username:          boshUser,
			Password:          boshPassword,
			BoshURI:           boshURI,
			Port:              boshPort,
			SkipSSLValidation: skipSSLValidation,
		}
		boshClient := bosh.NewClient(boshConfig)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("CF\t- Fetching Security Groups...")
		allSecGroups, err := client.ListSecGroups()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("BOSH\t- Finding CF deployment...")
		deployment, err := boshClient.SearchDeployment("^cf-.+")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("BOSH\t- Fetching DEA/Diego Cell VM details...")
		boshVMs, err := boshClient.GetRuntimeVMs(deployment)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("BOSH\t- Fetching DEA/Diego Cell VM IPs...")
		sources := boshVMs.GetAllIPs()
		fmt.Println("Virgil\t- Filtering for 'used' Security Groups...")
		secGroups := utility.GetUsedSecGroups(allSecGroups)
		fmt.Println("Virgil\t- Generating Firewall Rules...")
		firewallRules := utility.GetFirewallRules(sources, secGroups)
		fmt.Println("Virgil\t- Marshalling Firewall Rules to YAML...")
		yml, err := yaml.Marshal(&firewallRules)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		ioutil.WriteFile(c.Args()[0], []byte(fmt.Sprintf("---\n%v", string(yml))), os.FileMode(0644))
		fmt.Println("Firewall Policy written to file: ", c.Args()[0])
	}
	app.Run(os.Args)
}
