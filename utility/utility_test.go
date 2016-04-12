package utility_test

import (
	"github.com/FidelityInternational/virgil/utility"
	"github.com/cloudfoundry-community/go-cfclient"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

			rules, err := utility.ProcessRule(securityGroupRule1, []utility.FirewallRule{})
			Expect(err).To(BeNil())
			Expect(rules).To(HaveLen(7))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "12", Protocol: "tcp", Destination: []string{"1.1.1.1"}}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "15", Protocol: "tcp", Destination: []string{"1.1.1.1"}}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "16", Protocol: "tcp", Destination: []string{"1.1.1.1"}}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "17", Protocol: "tcp", Destination: []string{"1.1.1.1"}}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "18", Protocol: "tcp", Destination: []string{"1.1.1.1"}}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "19", Protocol: "tcp", Destination: []string{"1.1.1.1"}}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "20", Protocol: "tcp", Destination: []string{"1.1.1.1"}}))
			Expect(rules).ToNot(ContainElement(utility.FirewallRule{Port: "21", Protocol: "tcp", Destination: []string{"2.2.2.2"}}))
			rules, err = utility.ProcessRule(securityGroupRule2, rules)
			Expect(err).To(BeNil())
			Expect(rules).To(HaveLen(8))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "12", Protocol: "tcp", Destination: []string{"1.1.1.1", "2.2.2.2"}}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "15", Protocol: "tcp", Destination: []string{"1.1.1.1"}}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "16", Protocol: "tcp", Destination: []string{"1.1.1.1"}}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "17", Protocol: "tcp", Destination: []string{"1.1.1.1"}}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "18", Protocol: "tcp", Destination: []string{"1.1.1.1", "2.2.2.2"}}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "19", Protocol: "tcp", Destination: []string{"1.1.1.1", "2.2.2.2"}}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "20", Protocol: "tcp", Destination: []string{"1.1.1.1", "2.2.2.2"}}))
			Expect(rules).To(ContainElement(utility.FirewallRule{Port: "21", Protocol: "tcp", Destination: []string{"2.2.2.2"}}))
		})
	})

	Context("when ports cannot be expanded", func() {
		It("returns an error", func() {
			var securityGroupRule1 = cfclient.SecGroupRule{
				Ports:       "12,21-20",
				Protocol:    "tcp",
				Destination: "1.1.1.1",
			}
			rules, err := utility.ProcessRule(securityGroupRule1, []utility.FirewallRule{})
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
