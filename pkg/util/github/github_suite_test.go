package github_test

import (
	"net/http"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	util_github "github.com/devstream-io/devstream/pkg/util/github"
)

var (
	mux       *http.ServeMux
	serverURL string
	teardown  func()
)

var _ = BeforeSuite(func() {
	mux, serverURL, teardown = util_github.Setup()
})

var _ = AfterSuite(func() {
	teardown()
})

func TestPlanmanager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GitHub Suite")
}
