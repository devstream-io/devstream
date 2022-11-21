package step_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/step"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("GetStepGlobalVars func", func() {
	var (
		repoInfo *git.RepoInfo
	)
	BeforeEach(func() {
		repoInfo = &git.RepoInfo{}
	})
	When("repo type is gitlab and ssh key is not empty", func() {
		BeforeEach(func() {
			repoInfo.RepoType = "gitlab"
			repoInfo.SSHPrivateKey = "test"
		})
		It("should return gitlab ssh key", func() {
			v := step.GetStepGlobalVars(repoInfo)
			Expect(v.CredentialID).Should(Equal("gitlabCredential"))
		})
	})
	When("repo type is github", func() {
		BeforeEach(func() {
			repoInfo.RepoType = "github"
		})
		It("should return github ssh key", func() {
			v := step.GetStepGlobalVars(repoInfo)
			Expect(v.CredentialID).Should(Equal("githubCredential"))
		})
	})
	When("repo type is not valid", func() {
		BeforeEach(func() {
			repoInfo.RepoType = "not exist"
		})
		It("should return empty", func() {
			v := step.GetStepGlobalVars(repoInfo)
			Expect(v.ImageRepoSecret).Should(Equal("IMAGE_REPO_SECRET"))
		})
	})
})

var _ = Describe("ExtractValidStepConfig func", func() {
	type mockPlugin struct {
		ImageRepo *step.ImageRepoStepConfig
	}
	var (
		p        mockPlugin
		imageURL string
	)
	BeforeEach(func() {
		imageURL = "test"
		p = mockPlugin{
			ImageRepo: &step.ImageRepoStepConfig{
				URL: imageURL,
			},
		}
	})
	When("input type is pointer", func() {
		It("should return field", func() {
			stepAPI := step.ExtractValidStepConfig(&p)
			Expect(len(stepAPI)).Should(Equal(1))
		})
	})
	When("input type is struct", func() {
		It("should return field", func() {
			stepAPI := step.ExtractValidStepConfig(p)
			Expect(len(stepAPI)).Should(Equal(1))
		})
	})
})

var _ = Describe("GetRepoStepConfig func", func() {
	var r *git.RepoInfo
	When("repo type is github", func() {
		BeforeEach(func() {
			r = &git.RepoInfo{
				RepoType: "github",
			}
		})
		It("should return github stepConfig", func() {
			s := step.GetRepoStepConfig(r)
			Expect(len(s)).Should(Equal(1))
		})
	})
	When("repo type is gitlab", func() {
		BeforeEach(func() {
			r = &git.RepoInfo{
				RepoType: "gitlab",
			}
		})
		It("should return gitlab stepConfig", func() {
			s := step.GetRepoStepConfig(r)
			Expect(len(s)).Should(Equal(1))
		})
	})
	When("repo type is not valid", func() {
		BeforeEach(func() {
			r = &git.RepoInfo{
				RepoType: "not_exist",
			}
		})
		It("should return empty stepConfig", func() {
			s := step.GetRepoStepConfig(r)
			Expect(len(s)).Should(Equal(0))
		})
	})
})
