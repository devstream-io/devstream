package gitlab_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/git"
	"github.com/devstream-io/devstream/pkg/util/gitlab"
)

var _ = Describe("NewClient func", func() {
	var (
		gitlabToken string
		repoInfo    git.RepoInfo
	)
	When("gitlab Token is not set", func() {
		BeforeEach(func() {
			gitlabToken = os.Getenv("GITLAB_TOKEN")
			Expect(gitlabToken).Should(BeEmpty())
		})
		It("should return error", func() {
			_, err := gitlab.NewClient(&repoInfo)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("failed to read GITLAB_TOKEN from environment variable"))
		})
	})
	When("gitlab token is set", func() {
		BeforeEach(func() {
			os.Setenv("GITLAB_TOKEN", "test")
		})
		When("repoInfo field baseURL is empty", func() {
			BeforeEach(func() {
				repoInfo.BaseURL = ""
			})
			It("should return client with gitlab url", func() {
				client, err := gitlab.NewClient(&repoInfo)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.Client.BaseURL().Host).Should(Equal("gitlab.com"))
			})
		})
		When("repoInfo field baseURL is set", func() {
			var baseURL string
			BeforeEach(func() {
				baseURL = "test.com"
				repoInfo.BaseURL = fmt.Sprintf("http://%s", baseURL)
			})
			It("should return self host url", func() {
				client, err := gitlab.NewClient(&repoInfo)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.Client.BaseURL().Host).Should(Equal(baseURL))
			})
			AfterEach(func() {
				repoInfo.BaseURL = ""
			})
		})
		AfterEach(func() {
			os.Unsetenv("GITLAB_TOKEN")
		})
	})
})
