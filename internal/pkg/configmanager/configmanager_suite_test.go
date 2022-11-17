package configmanager

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPlanmanager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "configmanager Suite")
}
