package scm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm"
)

var _ = Describe("NewRepoFromURL func", func() {
	var (
		scmInfo *scm.SCMInfo
	)
	BeforeEach(func() {
		scmInfo = &scm.SCMInfo{}
	})
	When("is github repo", func() {
		BeforeEach(func() {
			scmInfo.CloneURL = "git@github.com:test/dtm-test.git"
			scmInfo.Branch = "test"
		})
		It("should return github repo info", func() {
			repo, err := scmInfo.BuildRepoInfo()
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
			scmInfo.CloneURL = "git@github.comtest/dtm-test.git"
			scmInfo.Branch = "test"
		})
		It("should return error", func() {
			_, err := scmInfo.BuildRepoInfo()
			Expect(err).Error().Should(HaveOccurred())
		})
	})
	When("is gitlab repo", func() {
		BeforeEach(func() {
			scmInfo.Type = "gitlab"
		})
		When("apiURL is not set, cloneURL is ssh format", func() {
			BeforeEach(func() {
				scmInfo.CloneURL = "git@gitlab.test.com:root/test-demo.git"
				scmInfo.APIURL = ""
				scmInfo.Branch = "test"
			})
			It("should return error", func() {
				_, err := scmInfo.BuildRepoInfo()
				Expect(err).Error().Should(HaveOccurred())
			})
		})
		When("apiURL is not set, cloneURL is http format", func() {
			BeforeEach(func() {
				scmInfo.CloneURL = "http://gitlab.test.com:3000/root/test-demo.git"
				scmInfo.APIURL = ""
				scmInfo.Branch = "test"
			})
			It("should return error", func() {
				repo, err := scmInfo.BuildRepoInfo()
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(repo.BaseURL).Should(Equal("http://gitlab.test.com:3000"))
				Expect(repo.Owner).Should(Equal("root"))
				Expect(repo.Repo).Should(Equal("test-demo"))
				Expect(repo.Branch).Should(Equal("test"))
			})
		})
		When("apiURL is set", func() {
			BeforeEach(func() {
				scmInfo.CloneURL = "git@gitlab.test.com:root/test-demo.git"
				scmInfo.APIURL = "http://gitlab.http.com"
				scmInfo.Branch = "test"
			})
			It("should set apiURL as BaseURL", func() {
				repo, err := scmInfo.BuildRepoInfo()
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(repo.BaseURL).Should(Equal("http://gitlab.http.com"))
				Expect(repo.Owner).Should(Equal("root"))
				Expect(repo.Repo).Should(Equal("test-demo"))
				Expect(repo.Branch).Should(Equal("test"))
			})

		})
	})
})
