package bosh_test

import (
	. "github.com/FidelityInternational/virgil/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/FidelityInternational/virgil/Godeps/_workspace/src/github.com/onsi/gomega"
	"testing"
)

func TestBosh(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Bosh test suite")
}
