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

	Context("getRepoDownloadURL method", func() {
		It("should return url", func() {
			githubURL := repo.getRepoDownloadURL()
			Expect(githubURL).Should(Equal("https://codeload.github.com/test_org/test_repo/zip/refs/heads/test_branch"))
		})
	})

	Context("BuildURL method", func() {
		When("repo is empty", func() {
			BeforeEach(func() {
				repo.RepoType = "not_exist"
			})
			It("should return empty url", func() {
				url := repo.BuildURL()
				Expect(url).Should(BeEmpty())
			})
		})
		When("repo is github", func() {
			BeforeEach(func() {
				repo.RepoType = "github"
			})
			It("should return github url", func() {
				url := repo.BuildURL()
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
				url := repo.BuildURL()
				Expect(url).Should(Equal(fmt.Sprintf("%s/%s/%s.git", repo.BaseURL, repo.Owner, repo.Repo)))
			})
		})
	})
})

var _ = Describe("NewRepoFromURL func", func() {
	var (
		repoType, apiURL, cloneURL, branch string
	)
	When("is github repo", func() {
		BeforeEach(func() {
			cloneURL = "git@github.com:test/dtm-test.git"
			branch = "test"
		})
		It("should return github repo info", func() {
			repo, err := NewRepoFromURL(repoType, apiURL, cloneURL, branch)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(repo).ShouldNot(BeNil())
			Expect(repo.Repo).Should(Equal("dtm-test"))
			Expect(repo.Owner).Should(Equal("test"))
			Expect(repo.RepoType).Should(Equal("github"))
			Expect(repo.Branch).Should(Equal("test"))
		})
	})
	When("clone url is not valid", func() {
		BeforeEach(func() {
			cloneURL = "git@github.comtest/dtm-test.git"
			branch = "test"
		})
		It("should return error", func() {
			_, err := NewRepoFromURL(repoType, apiURL, cloneURL, branch)
			Expect(err).Error().Should(HaveOccurred())
		})
	})
	When("is gitlab repo", func() {
		BeforeEach(func() {
			repoType = "gitlab"
		})
		When("apiURL is not set, cloneURL is ssh format", func() {
			BeforeEach(func() {
				cloneURL = "git@gitlab.test.com:root/test-demo.git"
				apiURL = ""
				branch = "test"
			})
			It("should return error", func() {
				_, err := NewRepoFromURL(repoType, apiURL, cloneURL, branch)
				Expect(err).Error().Should(HaveOccurred())
			})
		})
		When("apiURL is not set, cloneURL is http format", func() {
			BeforeEach(func() {
				cloneURL = "http://gitlab.test.com:3000/root/test-demo.git"
				apiURL = ""
				branch = "test"
			})
			It("should return error", func() {
				repo, err := NewRepoFromURL(repoType, apiURL, cloneURL, branch)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(repo.BaseURL).Should(Equal("http://gitlab.test.com:3000"))
				Expect(repo.Owner).Should(Equal("root"))
				Expect(repo.Repo).Should(Equal("test-demo"))
				Expect(repo.Branch).Should(Equal("test"))
			})
		})
		When("apiURL is set", func() {
			BeforeEach(func() {
				cloneURL = "git@gitlab.test.com:root/test-demo.git"
				apiURL = "http://gitlab.http.com"
				branch = "test"
			})
			It("should set apiURL as BaseURL", func() {
				repo, err := NewRepoFromURL(repoType, apiURL, cloneURL, branch)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(repo.BaseURL).Should(Equal("http://gitlab.http.com"))
				Expect(repo.Owner).Should(Equal("root"))
				Expect(repo.Repo).Should(Equal("test-demo"))
				Expect(repo.Branch).Should(Equal("test"))
			})

		})
	})
})
