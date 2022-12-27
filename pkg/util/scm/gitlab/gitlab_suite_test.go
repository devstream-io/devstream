package gitlab_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	apiRootPath string
)

func TestPlanmanager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GitLab Suite")
}

var _ = BeforeSuite(func() {
	apiRootPath = "/api/v4/"
})
