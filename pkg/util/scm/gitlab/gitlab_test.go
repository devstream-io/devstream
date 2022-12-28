package gitlab_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/scm/gitlab"
)

var _ = Describe("NewClient func", func() {
	var (
		repoInfo *git.RepoInfo
	)
	When("gitlab Token is not set", func() {
		BeforeEach(func() {
			repoInfo = &git.RepoInfo{
				BaseURL: "test",
				Repo:    "no_token",
			}
		})
		It("should return error", func() {
			_, err := gitlab.NewClient(repoInfo)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("config field scm.token is not setted"))
		})
	})
	When("gitlab token is set", func() {
		When("repoInfo field baseURL is empty", func() {
			BeforeEach(func() {
				repoInfo = &git.RepoInfo{
					BaseURL: "",
					Token:   "exist_token",
				}
			})
			It("should return client with gitlab url", func() {
				client, err := gitlab.NewClient(repoInfo)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.Client.BaseURL().Host).Should(Equal("gitlab.com"))
			})
		})
		When("repoInfo field baseURL is set", func() {
			var baseURL string
			BeforeEach(func() {
				baseURL = "test.com"
				repoInfo = &git.RepoInfo{
					BaseURL: fmt.Sprintf("http://%s", baseURL),
					Token:   "exist",
				}
			})
			It("should return self host url", func() {
				client, err := gitlab.NewClient(repoInfo)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(client.Client.BaseURL().Host).Should(Equal(baseURL))
			})
		})
	})
})
