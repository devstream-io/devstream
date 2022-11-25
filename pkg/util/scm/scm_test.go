package scm_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("BuildRepoInfo func", func() {
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
		When("apiURL is not set, url is ssh format", func() {
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
		When("apiURL is not set, url is http format", func() {
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
				scmInfo.Org = "cover_org"
			})
			It("should set apiURL as BaseURL", func() {
				repo, err := scmInfo.BuildRepoInfo()
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(repo.BaseURL).Should(Equal("http://gitlab.http.com"))
				Expect(repo.Owner).Should(Equal(""))
				Expect(repo.Repo).Should(Equal("test-demo"))
				Expect(repo.Branch).Should(Equal("test"))
				Expect(repo.Org).Should(Equal("cover_org"))
			})

		})
	})
	When("scm repo has fields", func() {
		BeforeEach(func() {
			scmInfo = &scm.SCMInfo{
				CloneURL: "",
				Name:     "test",
				Type:     "github",
				Org:      "test_org",
			}
		})
		It("should return repoInfo", func() {
			r, err := scmInfo.BuildRepoInfo()
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(r).Should(Equal(&git.RepoInfo{
				Branch:   "main",
				Repo:     "test",
				RepoType: "github",
				NeedAuth: true,
				CloneURL: "https://github.com/test_org/test",
				Org:      "test_org",
			}))
		})
	})
})

var _ = Describe("SCMInfo struct", func() {
	var s *scm.SCMInfo
	Context("Encode method", func() {
		BeforeEach(func() {
			s = &scm.SCMInfo{
				Name:  "test",
				Type:  "github",
				Owner: "test_user",
			}
		})
		It("should return map", func() {
			Expect(s.Encode()).Should(Equal(map[string]any{
				"owner":   "test_user",
				"scmType": "github",
				"name":    "test",
			}))
		})
	})
})
