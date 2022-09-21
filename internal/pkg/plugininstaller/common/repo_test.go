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
