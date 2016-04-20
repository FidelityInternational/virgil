# virgil

[![codecov.io](https://codecov.io/github/FidelityInternational/virgil/coverage.svg?branch=master)](https://codecov.io/github/FidelityInternational/virgil?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/FidelityInternational/virgil)](https://goreportcard.com/report/github.com/FidelityInternational/virgil) [![Build Status](https://travis-ci.org/FidelityInternational/virgil.svg?branch=master)](https://travis-ci.org/FidelityInternational/virgil)

### Overview

`virgil` generates Firewall Policies based on Cloud Foundry security groups.

The BOSH API is used to generate the source section of the firewall rules based on jobs names as "dea-partition*" or "diego_cell-partition".

The Cloud Foundry security groups are used for Destination, Port and Protocol and the rule set is then compressed to remove duplicates.

### Usage

#### As a CLI

`virgil` can be used as a CLI to produce YAML files of the generated polcies.

```
go get github.com/FidelityInternational/virgil
virgil \
--cf-system-domain='domain.example.com' \
--cf-user='cf_admin_user' \
--cf-password='cf_admin_password' \
--bosh-user='bosh_username' \
--bosh-password='bosh_password' \
--bosh-uri='bosh.example.com' \
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
  github.com/cloudfoundry-community/go-cfclient
  gopkg.in/FidelityInternational/virgil.v1/bosh
  gopkg.in/FidelityInternational/virgil.v1/utility
)
cfConfig := &cfclient.Config{
  ApiAddress:        "https://api.domain.example.com",
  Username:          "cf_admin_user",
  Password:          "cf_admin_password",
  SkipSslValidation: false,
}
boshConfig := &bosh.Config{
  Username:          "bosh_username",
  Password:          "bosh_password",
  BoshURI:           "bosh.example.com",
  Port:              "25555",
  SkipSSLValidation: false,
}
cfClient, _ := cfclient.NewClient(cfConfig)
boshClient := bosh.NewClient(boshConfig)
allSecGroups, _ := cfClient.ListSecGroups()
deployment, _ := boshClient.SearchDeployment("^cf-.+")
boshVMs, _ := boshClient.GetRuntimeVMs(deployment)
sources := boshVMs.GetAllIPs()
secGroups := utility.GetUsedSecGroups(allSecGroups)
firewallRules := utility.GetFirewallRules(sources, secGroups)
```
