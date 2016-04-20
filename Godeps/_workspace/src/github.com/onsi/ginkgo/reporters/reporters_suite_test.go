package reporters_test

import (
	. "github.com/FidelityInternational/virgil/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/FidelityInternational/virgil/Godeps/_workspace/src/github.com/onsi/gomega"

	"testing"
)

func TestReporters(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Reporters Suite")
}
