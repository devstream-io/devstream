package scm

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Repo Struct", func() {
	var (
		repo                                                        *Repo
		owner, org, repoName, branch, repoType, baseUrl, visibility string
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
			Owner:      owner,
			Org:        org,
			Repo:       repoName,
			Branch:     branch,
			RepoType:   repoType,
			BaseURL:    baseUrl,
			Visibility: visibility,
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

	Context("getBranch method", func() {
		When("repo is gitlab and branch is empty", func() {
			BeforeEach(func() {
				repo.Branch = ""
				repo.RepoType = "gitlab"
			})
			It("should return master branch", func() {
				branch := repo.getBranch()
				Expect(branch).Should(Equal("master"))
			})
		})
		When("repo is github and branch is empty", func() {
			BeforeEach(func() {
				repo.Branch = ""
				repo.RepoType = "github"
			})
			It("should return main branch", func() {
				branch := repo.getBranch()
				Expect(branch).Should(Equal("main"))
			})
		})
	})

	Context("BuildScmURL method", func() {
		When("repo is empty", func() {
			BeforeEach(func() {
				repo.RepoType = "not_exist"
			})
			It("should return empty url", func() {
				url := repo.BuildScmURL()
				Expect(url).Should(BeEmpty())
			})
		})
		When("repo is github", func() {
			BeforeEach(func() {
				repo.RepoType = "github"
			})
			It("should return github url", func() {
				url := repo.BuildScmURL()
				Expect(url).Should(Equal(fmt.Sprintf("https://github.com/%s/%s", repo.Org, repo.Repo)))
			})
		})
		When("repo is gitlab", func() {
			BeforeEach(func() {
				repo.RepoType = "gitlab"
				repo.BaseURL = "http://test.com"
				repo.Org = ""
			})
			It("should return gitlab url", func() {
				url := repo.BuildScmURL()
				Expect(url).Should(Equal(fmt.Sprintf("%s/%s/%s.git", repo.BaseURL, repo.Owner, repo.Repo)))
			})
		})
	})
})
