package gitlab_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	existGitlabToken, apiRootPath string
)

func TestPlanmanager(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GitLab Suite")
}

var _ = BeforeSuite(func() {
	apiRootPath = "/api/v4/"
	existGitlabToken := os.Getenv("GITLAB_TOKEN")
	if existGitlabToken != "" {
		os.Unsetenv("GITLAB_TOKEN")
	}
})

var _ = AfterSuite(func() {
	if existGitlabToken != "" {
		os.Setenv("GITLAB_TOKEN", existGitlabToken)
	}
})
