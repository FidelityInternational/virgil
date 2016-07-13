package bosh_test

import (
	"github.com/FidelityInternational/virgil/bosh"
	. "github.com/cloudfoundry-community/gogobosh"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("#FindDeployment", func() {
	var deployments []Deployment

	BeforeEach(func() {
		deployments = []Deployment{
			{
				Name:        "cf-warden-12345",
				CloudConfig: "none",
				Releases: []Resource{
					Resource{
						Name:    "cf",
						Version: "223",
					},
				},
				Stemcells: []Resource{
					Resource{
						Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
						Version: "3126",
					},
				},
			},
			{
				Name:        "cf-garden-12345",
				CloudConfig: "none",
				Releases: []Resource{
					Resource{
						Name:    "cf",
						Version: "223",
					},
				},
				Stemcells: []Resource{
					Resource{
						Name:    "bosh-warden-boshlite-ubuntu-trusty-go_agent",
						Version: "3126",
					},
				},
			},
		}
	})

	Context("when a deployment can be found", func() {
		It("finds the first matching deployment name based on a regex", func() {
			Ω(bosh.FindDeployment(deployments, "cf-garden*")).Should(Equal("cf-garden-12345"))
		})
	})

	Context("when a deployment cannot be found", func() {
		It("returns an empty string", func() {
			Ω(bosh.FindDeployment(deployments, "bosh*")).Should(BeEmpty())
		})
	})
})

var _ = Describe("#FindVMs", func() {
	It("Returns an array of all VMs matching the given regex", func() {
		vms := []VM{
			{
				IPs:     []string{"1.1.1.1"},
				JobName: "etcd_server-12344",
			},
			{
				IPs:     []string{"4.4.4.4"},
				JobName: "consul_server-567887",
			},
			{
				IPs:     []string{"3.3.3.3"},
				JobName: "etcd_server-98764",
			},
			{
				IPs:     []string{"4.4.4.4"},
				JobName: "consul_server-12344",
			},
			{
				IPs:     []string{"5.5.5.5"},
				JobName: "etcd_server-567887",
			},
		}
		matchedVMs := bosh.FindVMs(vms, "^etcd_server.+$")
		Ω(matchedVMs).Should(HaveLen(3))
		Ω(matchedVMs).Should(ContainElement(VM{
			IPs:     []string{"1.1.1.1"},
			JobName: "etcd_server-12344",
		}))
		Ω(matchedVMs).Should(ContainElement(VM{
			IPs:     []string{"3.3.3.3"},
			JobName: "etcd_server-98764",
		}))
		Ω(matchedVMs).Should(ContainElement(VM{
			IPs:     []string{"5.5.5.5"},
			JobName: "etcd_server-567887",
		}))
	})
})

var _ = Describe("#GetAllIPs", func() {
	It("return IPs for the provided VMs", func() {
		var deploymentVMs = []VM{
			{
				JobName: "dea-partition-d284104a9345228c01e2",
				Index:   0,
				VMCID:   "11",
				AgentID: "11",
				IPs:     []string{"11.11.11.11"},
			},
			{
				JobName: "dea-partition-d284104a9345228c01e2",
				Index:   1,
				VMCID:   "2",
				AgentID: "2",
				IPs:     []string{"2.2.2.2"},
			},
			{
				JobName: "dea-partition-d284104a9345228c01e2",
				Index:   2,
				VMCID:   "6",
				AgentID: "6",
				IPs:     []string{"6.6.6.6"},
			},
			{
				JobName: "dea-partition-d284104a9345228c01e2",
				Index:   3,
				VMCID:   "7",
				AgentID: "7",
				IPs:     []string{"7.7.7.7"},
			},
			{
				JobName: "dea-partition-d284104a9345228c01e2",
				Index:   4,
				VMCID:   "8",
				AgentID: "8",
				IPs:     []string{"8.8.8.8"},
			},
			{
				JobName: "dea-partition-d284104a9345228c01e2",
				Index:   5,
				VMCID:   "9",
				AgentID: "9",
				IPs:     []string{"9.9.9.9"},
			},
			{
				JobName: "dea-partition-d284104a9345228c01e2",
				Index:   6,
				VMCID:   "10",
				AgentID: "10",
				IPs:     []string{"10.10.10.10"},
			},
			{
				JobName: "diego_cell-partition-d284104a9345228c01e2",
				Index:   0,
				VMCID:   "4",
				AgentID: "4",
				IPs:     []string{"4.4.4.4"},
			},
			{
				JobName: "diego_cell-partition-d284104a9345228c01e2",
				Index:   1,
				VMCID:   "5",
				AgentID: "5",
				IPs:     []string{"5.5.5.5"},
			},
		}
		vmIPs := bosh.GetAllIPs(deploymentVMs)
		Expect(vmIPs).To(HaveLen(9))
		Expect(vmIPs).To(ContainElement("11.11.11.11"))
		Expect(vmIPs).To(ContainElement("2.2.2.2"))
		Expect(vmIPs).To(ContainElement("4.4.4.4"))
		Expect(vmIPs).To(ContainElement("5.5.5.5"))
		Expect(vmIPs).To(ContainElement("6.6.6.6"))
		Expect(vmIPs).To(ContainElement("7.7.7.7"))
		Expect(vmIPs).To(ContainElement("8.8.8.8"))
		Expect(vmIPs).To(ContainElement("9.9.9.9"))
		Expect(vmIPs).To(ContainElement("10.10.10.10"))
	})
})
