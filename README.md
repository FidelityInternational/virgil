# virgil

[![codecov.io](https://codecov.io/github/FidelityInternational/virgil/coverage.svg?branch=master)](https://codecov.io/github/FidelityInternational/virgil?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/FidelityInternational/virgil)](https://goreportcard.com/report/github.com/FidelityInternational/virgil) [![Build Status](https://travis-ci.org/FidelityInternational/virgil.svg?branch=master)](https://travis-ci.org/FidelityInternational/virgil)

### Overview

`virgil` generates Firewall Policies based on Cloud Foundry security groups.

The BOSH API is used to generate the source section of the firewall rules based on jobs names as "dea-partition*" or "diego_cell-partition".

The Cloud Foundry security groups are used for Destination, Port and Protocol and the rule set is then compressed to remove duplicates.

### Usage

#### As a CLI

`virgil` can be used as a CLI to produce YAML files of the generated policies.

```
go get github.com/FidelityInternational/virgil
virgil \
--cf-system-domain='domain.example.com' \
--cf-user='cf_admin_user' \
--cf-password='cf_admin_password' \
--bosh-user='bosh_username' \
--bosh-password='bosh_password' \
--bosh-uri='https://bosh.example.com:25555' \
output_file_name.yml
```

Additional parameters available are `--bosh-port` and `--skip-ssl-validation`.

To get additional help with the CLI use:

```
virgil --help
```

#### As a library

`virgil` can also be used as a library to plug in to other tools to act directly on the generated objects.

The library option is currently just a stripped back version of the CLI code, the aim is to rework this in future versions to make it easier to use.

```
import (
  "github.com/cloudfoundry-community/go-cfclient"
  "gopkg.in/FidelityInternational/virgil.v2/bosh"
  "gopkg.in/FidelityInternational/virgil.v2/utility"
  "github.com/cloudfoundry-community/gogobosh"
)
cfConfig := &cfclient.Config{
  ApiAddress:        "https://api.domain.example.com",
  Username:          "cf_admin_user",
  Password:          "cf_admin_password",
  SkipSslValidation: false,
}
boshConfig := &gogobosh.Config{
  Username:          "bosh_username",
  Password:          "bosh_password",
  BOSHAddress:       "bosh.example.com",
  SkipSslValidation: false,
}
cfClient, _ := cfclient.NewClient(cfConfig)
boshClient := gogobosh.NewClient(boshConfig)
allSecGroups, _ := cfClient.ListSecGroups()
deployments, _ := boshClient.GetDeployments()
deployment := bosh.FindDeployment(deployments, "^cf-.+")
boshVMs, _ := boshClient.GetDeploymentVMs(deployment)
runtimeVMs := bosh.FindVMs(boshVMs, "^(dea|diego_cell)-partition.+")
sources := bosh.GetAllIPs(runtimeVMs)
secGroups := utility.GetUsedSecGroups(allSecGroups)
firewallRules := utility.GetFirewallRules(sources, secGroups)
```

### Contributing

The `master` branch of this repository is under control of a [Concourse](https://concourse.ci/) only - No puny humans may ever commit anything into there.

To develop on this repo, you will need to...

- create your own development branch/ fork
- raise a PR against `staging`

Once you've PR'd and been merged into `staging`, Concourse will make sure everything passes all the tests and merge it to `master` for you.

This concourse pipeline is actively monitored but if you have a PR merged and do not see the changes appear in master within 30 minutes or so please feel free to raise an issue.
