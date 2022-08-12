package common

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo Struct", func() {
	var (
		repo                                                                           *Repo
		owner, org, repoName, branch, pathWithNamespace, repoType, baseUrl, visibility string
	)
	BeforeEach(func() {
		owner = "test_owner"
		org = "test_org"
		repoName = "test_repo"
		branch = "test_branch"
		repoType = "github"
		baseUrl = "http://test.gitlab.com"
		visibility = "public"
		repo = &Repo{
			Owner:             owner,
			Org:               org,
			Repo:              repoName,
			Branch:            branch,
			PathWithNamespace: pathWithNamespace,
			RepoType:          repoType,
			BaseURL:           baseUrl,
			Visibility:        visibility,
		}
	})

	Context("BuildRepoRenderConfig method", func() {
		It("should return map", func() {
			config := repo.BuildRepoRenderConfig()
			appName, ok := config["AppName"]
			Expect(ok).Should(BeTrue())
			Expect(appName).Should(Equal(repoName))
			repoInfo, ok := config["Repo"]
			Expect(ok).Should(BeTrue())
			repoInfoMap := repoInfo.(map[string]string)
			repoName, ok := repoInfoMap["Name"]
			Expect(ok).Should(BeTrue())
			Expect(repoName).Should(Equal(repoName))
			owner, ok := repoInfoMap["Owner"]
			Expect(ok).Should(BeTrue())
			Expect(owner).Should(Equal(org))
		})
	})
})
