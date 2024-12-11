package main

import (
	"fmt"
	"github.com/FidelityInternational/virgil/bosh"
	"github.com/FidelityInternational/virgil/utility"
	"github.com/cloudfoundry-community/gogobosh"
	"github.com/cloudfoundry/go-cfclient/v3/client"
	"github.com/cloudfoundry/go-cfclient/v3/config"
	"github.com/cloudfoundry/go-cfclient/v3/resource"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
	"os"
	"sort"
)

func main() {
	var (
		systemDomain, cfUser, cfPassword, boshUser, boshPassword, boshURI string
		skipSSLValidation                                                 = false
	)

	app := cli.NewApp()
	app.Name = "virgil"
	app.Usage = "A CLI App to return a list of firewall rules based on Cloud Foundry Security Groups"
	app.UsageText = "virgil [options] output_file"
	app.Version = "1.0.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "cf-system-domain, csd",
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
		cli.BoolFlag{
			Name:        "skip-ssl-validation, skip-ssl",
			Usage:       "Skip SSL Validation",
			Destination: &skipSSLValidation,
		},
	}
	app.Action = func(c *cli.Context) error {
		if systemDomain == "" || cfUser == "" || cfPassword == "" || c.NArg() == 0 || boshUser == "" || boshPassword == "" || boshURI == "" {
			fmt.Println("cf-system-domain, cf-user, cf-password, bosh-user, bosh-password, bosh-uri and output_file must all be set")
			os.Exit(1)
		}
		config, err := config.New(fmt.Sprintf("https://api.%s", systemDomain), config.UserPassword(cfUser, cfPassword))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		boshConfig := &gogobosh.Config{
			Username:          boshUser,
			Password:          boshPassword,
			BOSHAddress:       boshURI,
			SkipSslValidation: skipSSLValidation,
		}
		client, err := client.New(config)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		boshClient, err := gogobosh.NewClient(boshConfig)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("CF\t- Fetching Security Groups...")
		allSecGroups, err := client.SecurityGroups.ListAll(nil, nil)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("BOSH\t- Finding CF deployment...")
		deployments, err := boshClient.GetDeployments()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		deployment := bosh.FindDeployment(deployments, "^cf.*")
		fmt.Println("BOSH\t- Fetching DEA/Diego Cell VM details...")
		boshVMs, err := boshClient.GetDeploymentVMs(deployment)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		fmt.Println("BOSH\t- Fetching DEA/Diego Cell VM IPs...")
		runtimeVMs := bosh.FindVMs(boshVMs, "^(dea|diego_cell|diego-cell).*")
		sources := bosh.GetAllIPs(runtimeVMs)
		sort.Strings(sources)
		fmt.Println("Virgil\t- Filtering for 'used' Security Groups...")
		var secGroupsList []resource.SecurityGroup
		for _, sg := range allSecGroups {
			secGroupsList = append(secGroupsList, *sg)
		}
		secGroups := utility.GetUsedSecGroups(secGroupsList)
		fmt.Println("Virgil\t- Generating Firewall Rules...")
		firewallRules := utility.GetFirewallRules(sources, secGroups)
		fmt.Println("Virgil\t- Marshalling Firewall Rules to YAML...")
		yml, err := yaml.Marshal(&firewallRules)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		os.WriteFile(c.Args()[0], []byte(fmt.Sprintf("---\n%v", string(yml))), os.FileMode(0644))
		fmt.Println("Firewall Policy written to file: ", c.Args()[0])
		return nil
	}
	app.Run(os.Args)
}
