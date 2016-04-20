package utility_test

import (
	"github.com/FidelityInternational/virgil/Godeps/_workspace/src/github.com/cloudfoundry-community/go-cfclient"
	. "github.com/FidelityInternational/virgil/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/FidelityInternational/virgil/Godeps/_workspace/src/github.com/onsi/gomega"
	"github.com/FidelityInternational/virgil/utility"
	"sort"
)

var _ = Describe("#PortExpand", func() {
	Context("When the port string is valid", func() {
		It("returns an array of strings", func() {
			const portString1 = "3,5-7"
			const portString2 = "1,2,3"
			const portString3 = "1, 2, 3"
			const portString4 = "3-4"
			const portString5 = "3 - 4"
			Expect(utility.PortExpand(portString1)).To(Equal([]string{"3", "5", "6", "7"}))
			Expect(utility.PortExpand(portString2)).To(Equal([]string{"1", "2", "3"}))
			Expect(utility.PortExpand(portString3)).To(Equal([]string{"1", "2", "3"}))
			Expect(utility.PortExpand(portString4)).To(Equal([]string{"3", "4"}))
			Expect(utility.PortExpand(portString5)).To(Equal([]string{"3", "4"}))
		})
	})

	Context("When the port string is not valid", func() {
		Context("because it contains invalid ranges", func() {
			It("returns an error", func() {
				const portString1 = "3,6-6"
				const portString2 = "3,9-6"
				const portString3 = "3,3-6-8"
				const portString4 = "-4"
				const portString5 = "4-"
				const portString6 = "-"
				ports, err := utility.PortExpand(portString1)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port range 6-6 was invalid"))
				ports, err = utility.PortExpand(portString2)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port range 9-6 was invalid"))
				ports, err = utility.PortExpand(portString3)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port range 3-6-8 was invalid"))
				ports, err = utility.PortExpand(portString4)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port  was invalid as part of range -4"))
				ports, err = utility.PortExpand(portString5)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port  was invalid as part of range 4-"))
				ports, err = utility.PortExpand(portString6)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port  was invalid as part of range -"))
			})
		})

		Context("because it contains invalid characters", func() {
			It("returns an error", func() {
				const portString1 = "3,5-6,*"
				const portString2 = "#-7"
				const portString3 = "7-d"
				ports, err := utility.PortExpand(portString1)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port * was invalid"))
				ports, err = utility.PortExpand(portString2)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port # was invalid as part of range #-7"))
				ports, err = utility.PortExpand(portString3)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port d was invalid as part of range 7-d"))
			})
		})

		Context("because it contains values that are not valid ports", func() {
			It("returns an error", func() {
				const portString1 = "0"
				const portString2 = "3456789"
				const portString3 = "3,4-8,678800"
				const portString4 = "3,4-8,12-99999999"
				const portString5 = "3,4-8,1256789-98"
				const portString6 = "65536"
				ports, err := utility.PortExpand(portString1)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port 0 was invalid"))
				ports, err = utility.PortExpand(portString2)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port 3456789 was invalid"))
				ports, err = utility.PortExpand(portString3)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port 678800 was invalid"))
				ports, err = utility.PortExpand(portString4)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port 99999999 was invalid as part of range 12-99999999"))
				ports, err = utility.PortExpand(portString5)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port 1256789 was invalid as part of range 1256789-98"))
				ports, err = utility.PortExpand(portString6)
				Expect(ports).To(HaveLen(0))
				Expect(err).To(MatchError("Port 65536 was invalid"))
			})
		})
	})
})

var _ = Describe("#ProcessRule", func() {
	var source = []string{"1.2.3.4", "2.3.4.5"}
	Context("when ports can be expanded", func() {
		It("returns valid firewall rules", func() {
			var securityGroupRule1 = cfclient.SecGroupRule{
				Ports:       "12,15-20",
				Protocol:    "tcp",
				Destination: "1.1.1.1",
			}
			var securityGroupRule2 = cfclient.SecGroupRule{
				Ports:       "12,18-21",
				Protocol:    "tcp",
				Destination: "2.2.2.2",
			}

			rules, err := utility.ProcessRule(securityGroupRule1, []utility.FirewallRule{}, source)
			Expect(err).To(BeNil())
			Expect(rules).To(HaveLen(7))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "12", Protocol: "tcp", Destination: []string{"1.1.1.1"}, Source: source}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "15", Protocol: "tcp", Destination: []string{"1.1.1.1"}, Source: source}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "16", Protocol: "tcp", Destination: []string{"1.1.1.1"}, Source: source}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "17", Protocol: "tcp", Destination: []string{"1.1.1.1"}, Source: source}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "18", Protocol: "tcp", Destination: []string{"1.1.1.1"}, Source: source}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "19", Protocol: "tcp", Destination: []string{"1.1.1.1"}, Source: source}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "20", Protocol: "tcp", Destination: []string{"1.1.1.1"}, Source: source}))
			Expect(rules).ToNot(ContainElement(utility.FirewallRule{Port: "21", Protocol: "tcp", Destination: []string{"2.2.2.2"}, Source: source}))
			rules, err = utility.ProcessRule(securityGroupRule2, rules, source)
			Expect(err).To(BeNil())
			Expect(rules).To(HaveLen(8))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "12", Protocol: "tcp", Destination: []string{"1.1.1.1", "2.2.2.2"}, Source: source}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "15", Protocol: "tcp", Destination: []string{"1.1.1.1"}, Source: source}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "16", Protocol: "tcp", Destination: []string{"1.1.1.1"}, Source: source}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "17", Protocol: "tcp", Destination: []string{"1.1.1.1"}, Source: source}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "18", Protocol: "tcp", Destination: []string{"1.1.1.1", "2.2.2.2"}, Source: source}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "19", Protocol: "tcp", Destination: []string{"1.1.1.1", "2.2.2.2"}, Source: source}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "20", Protocol: "tcp", Destination: []string{"1.1.1.1", "2.2.2.2"}, Source: source}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "21", Protocol: "tcp", Destination: []string{"2.2.2.2"}, Source: source}))
		})
	})

	Context("when ports cannot be expanded", func() {
		It("returns an error", func() {
			var securityGroupRule1 = cfclient.SecGroupRule{
				Ports:       "12,21-20",
				Protocol:    "tcp",
				Destination: "1.1.1.1",
			}
			rules, err := utility.ProcessRule(securityGroupRule1, []utility.FirewallRule{}, source)
			Expect(rules).To(HaveLen(0))
			Expect(err).To(MatchError("Port range 21-20 was invalid"))
		})
	})
})

var _ = Describe("#GetUsedSecGroups", func() {
	It("returns an array of in use security groups", func() {
		var securityGroups = []cfclient.SecGroup{
			cfclient.SecGroup{
				Guid:       "1",
				Name:       "test-sec-group1",
				Running:    false,
				Staging:    false,
				SpacesData: []cfclient.SpaceResource{},
			},
			cfclient.SecGroup{
				Guid:       "2",
				Name:       "test-sec-group2",
				Running:    true,
				Staging:    false,
				SpacesData: []cfclient.SpaceResource{},
			},
			cfclient.SecGroup{
				Guid:       "3",
				Name:       "test-sec-group3",
				Running:    false,
				Staging:    true,
				SpacesData: []cfclient.SpaceResource{},
			},
			cfclient.SecGroup{
				Guid:       "4",
				Name:       "test-sec-group4",
				Running:    true,
				Staging:    true,
				SpacesData: []cfclient.SpaceResource{},
			},
			cfclient.SecGroup{
				Guid:    "5",
				Name:    "test-sec-group5",
				Running: false,
				Staging: false,
				SpacesData: []cfclient.SpaceResource{
					cfclient.SpaceResource{
						Meta:   cfclient.Meta{Guid: "1"},
						Entity: cfclient.Space{Guid: "1", Name: "test-space1"},
					},
				},
			},
		}
		usedSecGroups := utility.GetUsedSecGroups(securityGroups)
		Expect(usedSecGroups).To(HaveLen(4))
		Expect(usedSecGroups).ToNot(ContainElement(cfclient.SecGroup{
			Guid:       "1",
			Name:       "test-sec-group1",
			Running:    false,
			Staging:    false,
			SpacesData: []cfclient.SpaceResource{},
		}))
		Expect(usedSecGroups).To(ContainElement(cfclient.SecGroup{
			Guid:       "2",
			Name:       "test-sec-group2",
			Running:    true,
			Staging:    false,
			SpacesData: []cfclient.SpaceResource{},
		}))
		Expect(usedSecGroups).To(ContainElement(cfclient.SecGroup{
			Guid:       "3",
			Name:       "test-sec-group3",
			Running:    false,
			Staging:    true,
			SpacesData: []cfclient.SpaceResource{},
		}))
		Expect(usedSecGroups).To(ContainElement(cfclient.SecGroup{
			Guid:       "4",
			Name:       "test-sec-group4",
			Running:    true,
			Staging:    true,
			SpacesData: []cfclient.SpaceResource{},
		}))
		Expect(usedSecGroups).To(ContainElement(cfclient.SecGroup{
			Guid:    "5",
			Name:    "test-sec-group5",
			Running: false,
			Staging: false,
			SpacesData: []cfclient.SpaceResource{
				cfclient.SpaceResource{
					Meta:   cfclient.Meta{Guid: "1"},
					Entity: cfclient.Space{Guid: "1", Name: "test-space1"},
				},
			},
		}))
	})
})

var _ = Describe("#GetFirewallRules", func() {
	var source = []string{"1.2.3.4", "2.3.4.5"}

	It("Returns an array of Firewall Rules", func() {
		var securityGroups = []cfclient.SecGroup{
			cfclient.SecGroup{
				Guid:    "2",
				Name:    "test-sec-group2",
				Running: true,
				Staging: false,
				Rules: []cfclient.SecGroupRule{
					cfclient.SecGroupRule{
						Protocol:    "tcp",
						Ports:       "1,2-4",
						Destination: "2.2.2.2",
					},
					cfclient.SecGroupRule{
						Protocol:    "tcp",
						Ports:       "8",
						Destination: "5.5.5.5",
					},
				},
				SpacesData: []cfclient.SpaceResource{},
			},
			cfclient.SecGroup{
				Guid:    "3",
				Name:    "test-sec-group3",
				Running: false,
				Staging: true,
				Rules: []cfclient.SecGroupRule{
					cfclient.SecGroupRule{
						Protocol:    "tcp",
						Ports:       "2,3-4",
						Destination: "3.3.3.3",
					},
				},
				SpacesData: []cfclient.SpaceResource{},
			},
			cfclient.SecGroup{
				Guid:    "6",
				Name:    "test-sec-group6",
				Running: false,
				Staging: true,
				Rules: []cfclient.SecGroupRule{
					cfclient.SecGroupRule{
						Protocol:    "all",
						Destination: "9.9.9.9",
					},
				},
				SpacesData: []cfclient.SpaceResource{},
			},
			cfclient.SecGroup{
				Guid:    "4",
				Name:    "test-sec-group4",
				Running: true,
				Staging: true,
				Rules: []cfclient.SecGroupRule{
					cfclient.SecGroupRule{
						Protocol:    "tcp",
						Ports:       "1,4-7",
						Destination: "4.4.4.4",
					},
				},
				SpacesData: []cfclient.SpaceResource{},
			},
			cfclient.SecGroup{
				Guid:    "5",
				Name:    "test-sec-group5",
				Running: false,
				Staging: false,
				Rules: []cfclient.SecGroupRule{
					cfclient.SecGroupRule{
						Protocol:    "udp",
						Ports:       "2",
						Destination: "1.1.1.1",
					},
				},
				SpacesData: []cfclient.SpaceResource{
					cfclient.SpaceResource{
						Meta:   cfclient.Meta{Guid: "1"},
						Entity: cfclient.Space{Guid: "1", Name: "test-space1"},
					},
				},
			},
			cfclient.SecGroup{
				Guid:    "7",
				Name:    "test-sec-group7",
				Running: false,
				Staging: false,
				Rules: []cfclient.SecGroupRule{
					cfclient.SecGroupRule{
						Protocol:    "tcp",
						Ports:       "99",
						Destination: "9.9.9.9",
					},
				},
				SpacesData: []cfclient.SpaceResource{
					cfclient.SpaceResource{
						Meta:   cfclient.Meta{Guid: "1"},
						Entity: cfclient.Space{Guid: "1", Name: "test-space1"},
					},
				},
			},
			cfclient.SecGroup{
				Guid:    "8",
				Name:    "test-sec-group8",
				Running: false,
				Staging: false,
				Rules: []cfclient.SecGroupRule{
					cfclient.SecGroupRule{
						Protocol:    "tcp",
						Ports:       "100",
						Destination: "9.9.9.9",
					},
				},
				SpacesData: []cfclient.SpaceResource{
					cfclient.SpaceResource{
						Meta:   cfclient.Meta{Guid: "1"},
						Entity: cfclient.Space{Guid: "1", Name: "test-space1"},
					},
				},
			},
			cfclient.SecGroup{
				Guid:    "9",
				Name:    "test-sec-group9",
				Running: false,
				Staging: false,
				Rules: []cfclient.SecGroupRule{
					cfclient.SecGroupRule{
						Protocol:    "tcp",
						Ports:       "110-120",
						Destination: "9.9.9.9",
					},
				},
				SpacesData: []cfclient.SpaceResource{
					cfclient.SpaceResource{
						Meta:   cfclient.Meta{Guid: "1"},
						Entity: cfclient.Space{Guid: "1", Name: "test-space1"},
					},
				},
			},
			cfclient.SecGroup{
				Guid:    "10",
				Name:    "test-sec-group10",
				Running: false,
				Staging: false,
				Rules: []cfclient.SecGroupRule{
					cfclient.SecGroupRule{
						Protocol:    "tcp",
						Ports:       "116",
						Destination: "11.1.1.1",
					},
				},
				SpacesData: []cfclient.SpaceResource{
					cfclient.SpaceResource{
						Meta:   cfclient.Meta{Guid: "1"},
						Entity: cfclient.Space{Guid: "1", Name: "test-space1"},
					},
				},
			},
		}
		policy := utility.GetFirewallRules(source, securityGroups)
		rules := policy.FirewallRules
		Expect(policy.SchemaVersion).To(Equal("1"))
		Expect(rules).To(HaveLen(11))
		Expect(rules).To(ContainElement(utility.FirewallRule{Port: "1", Protocol: "tcp", Destination: []string{"2.2.2.2", "4.4.4.4"}, Source: source}))
		Expect(rules).To(ContainElement(utility.FirewallRule{Port: "2-3", Protocol: "tcp", Destination: []string{"2.2.2.2", "3.3.3.3"}, Source: source}))
		Expect(rules).To(ContainElement(utility.FirewallRule{Port: "4", Protocol: "tcp", Destination: []string{"2.2.2.2", "3.3.3.3", "4.4.4.4"}, Source: source}))
		Expect(rules).To(ContainElement(utility.FirewallRule{Port: "5-7", Protocol: "tcp", Destination: []string{"4.4.4.4"}, Source: source}))
		Expect(rules).To(ContainElement(utility.FirewallRule{Port: "8", Protocol: "tcp", Destination: []string{"5.5.5.5"}, Source: source}))
		Expect(rules).To(ContainElement(utility.FirewallRule{Port: "2", Protocol: "udp", Destination: []string{"1.1.1.1"}, Source: source}))
		Expect(rules).To(ContainElement(utility.FirewallRule{Port: "", Protocol: "all", Destination: []string{"9.9.9.9"}, Source: source}))
		Expect(rules).To(ContainElement(utility.FirewallRule{Port: "99-100", Protocol: "tcp", Destination: []string{"9.9.9.9"}, Source: source}))
		Expect(rules).To(ContainElement(utility.FirewallRule{Port: "110-115", Protocol: "tcp", Destination: []string{"9.9.9.9"}, Source: source}))
		Expect(rules).To(ContainElement(utility.FirewallRule{Port: "117-120", Protocol: "tcp", Destination: []string{"9.9.9.9"}, Source: source}))
		Expect(rules).To(ContainElement(utility.FirewallRule{Port: "116", Protocol: "tcp", Destination: []string{"9.9.9.9", "11.1.1.1"}, Source: source}))
	})
})

var _ = Describe("#ByPort", func() {
	It("sorts FirewallRules in order by port", func() {
		var firewallRules = []utility.FirewallRule{
			utility.FirewallRule{
				Port:     "5",
				Protocol: "udp",
			},
			utility.FirewallRule{
				Port:     "1",
				Protocol: "tcp",
			},
			utility.FirewallRule{
				Port:     "5",
				Protocol: "tcp",
			},
			utility.FirewallRule{
				Port:     "1",
				Protocol: "udp",
			},
			utility.FirewallRule{
				Port:     "3",
				Protocol: "tcp",
			},
			utility.FirewallRule{
				Port:     "9",
				Protocol: "udp",
			},
			utility.FirewallRule{
				Port:     "9",
				Protocol: "tcp",
			},
			utility.FirewallRule{
				Port:     "2",
				Protocol: "tcp",
			},
			utility.FirewallRule{
				Port:     "2",
				Protocol: "udp",
			},
			utility.FirewallRule{
				Port:     "3",
				Protocol: "udp",
			},
		}
		Expect(firewallRules).To(HaveLen(10))
		Expect(firewallRules[0].Port).To(Equal("5"))
		Expect(firewallRules[0].Protocol).To(Equal("udp"))
		Expect(firewallRules[1].Port).To(Equal("1"))
		Expect(firewallRules[1].Protocol).To(Equal("tcp"))
		Expect(firewallRules[2].Port).To(Equal("5"))
		Expect(firewallRules[2].Protocol).To(Equal("tcp"))
		Expect(firewallRules[3].Port).To(Equal("1"))
		Expect(firewallRules[3].Protocol).To(Equal("udp"))
		Expect(firewallRules[4].Port).To(Equal("3"))
		Expect(firewallRules[4].Protocol).To(Equal("tcp"))
		Expect(firewallRules[5].Port).To(Equal("9"))
		Expect(firewallRules[5].Protocol).To(Equal("udp"))
		Expect(firewallRules[6].Port).To(Equal("9"))
		Expect(firewallRules[6].Protocol).To(Equal("tcp"))
		Expect(firewallRules[7].Port).To(Equal("2"))
		Expect(firewallRules[7].Protocol).To(Equal("tcp"))
		Expect(firewallRules[8].Port).To(Equal("2"))
		Expect(firewallRules[8].Protocol).To(Equal("udp"))
		Expect(firewallRules[9].Port).To(Equal("3"))
		Expect(firewallRules[9].Protocol).To(Equal("udp"))
		sort.Sort(utility.ByPort(firewallRules))
		Expect(firewallRules).To(HaveLen(10))
		Expect(firewallRules[0].Port).To(Equal("1"))
		Expect(firewallRules[0].Protocol).To(Equal("tcp"))
		Expect(firewallRules[1].Port).To(Equal("2"))
		Expect(firewallRules[1].Protocol).To(Equal("tcp"))
		Expect(firewallRules[2].Port).To(Equal("3"))
		Expect(firewallRules[2].Protocol).To(Equal("tcp"))
		Expect(firewallRules[3].Port).To(Equal("5"))
		Expect(firewallRules[3].Protocol).To(Equal("tcp"))
		Expect(firewallRules[4].Port).To(Equal("9"))
		Expect(firewallRules[4].Protocol).To(Equal("tcp"))
		Expect(firewallRules[5].Port).To(Equal("1"))
		Expect(firewallRules[5].Protocol).To(Equal("udp"))
		Expect(firewallRules[6].Port).To(Equal("2"))
		Expect(firewallRules[6].Protocol).To(Equal("udp"))
		Expect(firewallRules[7].Port).To(Equal("3"))
		Expect(firewallRules[7].Protocol).To(Equal("udp"))
		Expect(firewallRules[8].Port).To(Equal("5"))
		Expect(firewallRules[8].Protocol).To(Equal("udp"))
		Expect(firewallRules[9].Port).To(Equal("9"))
		Expect(firewallRules[9].Protocol).To(Equal("udp"))
	})
})

var _ = Describe("RemoveDuplicates", func() {
	It("removes duplicates from an array of strings", func() {
		var strings = []string{"a", "a", "b", "c", "a", "b"}
		Expect(strings).To(HaveLen(6))
		utility.RemoveDuplicates(&strings)
		Expect(strings).To(HaveLen(3))
		Expect(strings).To(ContainElement("a"))
		Expect(strings).To(ContainElement("b"))
		Expect(strings).To(ContainElement("c"))
	})

})
