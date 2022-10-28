package step_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/step"
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
			Expect(v.CredentialID).Should(Equal("gitlabSSHKeyCredential"))
		})
	})
	When("repo type is not valid", func() {
		BeforeEach(func() {
			repoInfo.RepoType = "not exist"
		})
		It("should return empty", func() {
			v := step.GetStepGlobalVars(repoInfo)
			Expect(v.ImageRepoSecret).Should(Equal("repo-auth"))
		})
	})
})
