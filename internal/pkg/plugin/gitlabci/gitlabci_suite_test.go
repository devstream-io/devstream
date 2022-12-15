package gitlabci_test

import (
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestGitlabcedocker(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Plugin github-actions Suite")
}

var (
	gitlabEnv string
)

var _ = BeforeSuite(func() {
	gitlabEnv = os.Getenv("GITLAB_TOKEN")
	err := os.Unsetenv("GITLAB_TOKEN")
	Expect(err).Error().ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	if gitlabEnv != "" {
		os.Setenv("GITLAB_TOKEN", gitlabEnv)
	}
})
