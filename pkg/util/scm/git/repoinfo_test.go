package git_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("RepoInfo struct", func() {
	var (
		repoName, branch, owner, org string
		repoInfo                     *git.RepoInfo
	)
	BeforeEach(func() {
		repoName = "test_repo"
		branch = "test_branch"
		owner = "test_owner"
		org = "test_org"
		repoInfo = &git.RepoInfo{
			Repo:   repoName,
			Branch: branch,
			Owner:  owner,
			Org:    org,
		}
	})
	Context("GetRepoOwner method", func() {
		It("should return owner", func() {
			result := repoInfo.GetRepoOwner()
			Expect(result).Should(Equal(repoInfo.Org))
		})
	})
	Context("GetRepoPath method", func() {
		It("should return repo path", func() {
			result := repoInfo.GetRepoPath()
			Expect(result).Should(Equal(fmt.Sprintf("%s/%s", repoInfo.GetRepoOwner(), repoName)))
		})
	})

	Context("BuildRepoRenderConfig method", func() {
		It("should return map", func() {
			config := repoInfo.BuildRepoRenderConfig()
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

	Context("GetBranchWithDefault method", func() {
		When("repo is gitlab and branch is empty", func() {
			BeforeEach(func() {
				repoInfo.Branch = ""
				repoInfo.RepoType = "gitlab"
			})
			It("should return master branch", func() {
				branch := repoInfo.GetBranchWithDefault()
				Expect(branch).Should(Equal("master"))
			})
		})
		When("repo is github and branch is empty", func() {
			BeforeEach(func() {
				repoInfo.Branch = ""
				repoInfo.RepoType = "github"
			})
			It("should return main branch", func() {
				branch := repoInfo.GetBranchWithDefault()
				Expect(branch).Should(Equal("main"))
			})
		})
	})

	Context("BuildScmURL method", func() {
		When("repo is empty", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "not_exist"
			})
			It("should return empty url", func() {
				url := repoInfo.BuildScmURL()
				Expect(url).Should(BeEmpty())
			})
		})
		When("repo is github", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "github"
			})
			It("should return github url", func() {
				url := repoInfo.BuildScmURL()
				Expect(url).Should(Equal(fmt.Sprintf("https://github.com/%s/%s", repoInfo.Org, repoInfo.Repo)))
			})
		})
		When("repo is gitlab", func() {
			BeforeEach(func() {
				repoInfo.RepoType = "gitlab"
				repoInfo.BaseURL = "http://test.com"
				repoInfo.Org = ""
			})
			It("should return gitlab url", func() {
				url := repoInfo.BuildScmURL()
				Expect(url).Should(Equal(fmt.Sprintf("%s/%s/%s.git", repoInfo.BaseURL, repoInfo.Owner, repoInfo.Repo)))
			})
		})
	})
})
