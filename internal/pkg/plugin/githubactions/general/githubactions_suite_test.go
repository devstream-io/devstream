package general_test

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
	githubEnv, gitlabEnv string
)

var _ = BeforeSuite(func() {
	githubEnv = os.Getenv("GITHUB_TOKEN")
	gitlabEnv = os.Getenv("GITLAB_TOKEN")
	err := os.Unsetenv("GITHUB_TOKEN")
	Expect(err).Error().ShouldNot(HaveOccurred())
	err = os.Unsetenv("GITLAB_TOKEN")
	Expect(err).Error().ShouldNot(HaveOccurred())
})

var _ = AfterSuite(func() {
	if githubEnv != "" {
		os.Setenv("GITHUB_TOKEN", githubEnv)
	}
	if gitlabEnv != "" {
		os.Setenv("GITLAB_TOKEN", gitlabEnv)
	}
})
