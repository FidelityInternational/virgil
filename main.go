package main

import (
	"fmt"
	"github.com/FidelityInternational/virgil/bosh"
	"github.com/FidelityInternational/virgil/utility"
	"github.com/cloudfoundry-community/go-cfclient"
	"github.com/codegangsta/cli"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func main() {
	var systemDomain, cfUser, cfPassword, boshUser, boshPassword, boshURI, boshPort string
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
			Destination: &boshPort,
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
		if boshUser == "" || boshPassword == "" || boshURI == "" || boshPort == "" {
			fmt.Println("BOSH user, passowrd, URI and Port must all be set")
			os.Exit(1)
		}
		boshConfig := &bosh.Config{
			Username:          boshUser,
			Password:          boshPassword,
			BoshURI:           boshURI,
			Port:              boshPort,
			SkipSSLValidation: true,
		}
		boshClient := bosh.NewClient(boshConfig)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		allSecGroups, err := client.ListSecGroups()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		deployment, err := boshClient.SearchDeployment("^cf-.+")
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		boshVMs, err := boshClient.GetRuntimeVMs(deployment)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		sources := boshVMs.GetAllIPs()
		secGroups := utility.GetUsedSecGroups(allSecGroups)
		firewallRules := utility.GetFirewallRules(sources, secGroups)
		yml, _ := yaml.Marshal(&firewallRules)
		ioutil.WriteFile(c.Args()[0], []byte(fmt.Sprintf("---\n%v", string(yml))), os.FileMode(0444))
	}
	app.Run(os.Args)
}
