package common

import (
	"fmt"

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

	Context("CreateGithubClient method", func() {
		It("should return client", func() {
			client, err := repo.CreateGithubClient(false)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(client).ShouldNot(BeNil())
		})
	})

	//TODO(steinliber) add CreateAndRenderLocalRepo test

	Context("getBranch method", func() {
		When("branch is not set", func() {
			BeforeEach(func() {
				repo.Branch = ""
			})
			It("should return main branch", func() {
				repoBranch := repo.getBranch()
				Expect(repoBranch).Should(Equal("main"))
			})
			AfterEach(func() {
				repo.Branch = branch
			})
		})
	})

	Context("getRepoNameWithBranch method", func() {
		It("should return repo and branch name", func() {
			repoNameWithURL := repo.getRepoNameWithBranch()
			Expect(repoNameWithURL).Should(Equal(fmt.Sprintf("%s-%s", repoName, branch)))
		})
	})

	Context("getRepoOwner method", func() {
		When("org is not empty", func() {
			It("should return org", func() {
				ownerName := repo.getRepoOwner()
				Expect(ownerName).Should(Equal(org))
			})
		})
		When("org is empty", func() {
			BeforeEach(func() {
				repo.Org = ""
			})
			It("should return owner", func() {
				ownerName := repo.getRepoOwner()
				Expect(ownerName).Should(Equal(owner))
			})
			AfterEach(func() {
				repo.Org = org
			})
		})
	})
})
