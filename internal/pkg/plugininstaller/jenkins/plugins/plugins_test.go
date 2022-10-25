package plugins

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("EnsurePluginInstalled func", func() {
	var (
		mockClient *jenkins.MockClient
		s          *mockPluginsConfig
	)
	BeforeEach(func() {
		mockClient = &jenkins.MockClient{}
		s = &mockPluginsConfig{}
	})
	It("should work normal", func() {
		err := EnsurePluginInstalled(mockClient, []PluginConfigAPI{s})
		Expect(err).Error().ShouldNot(HaveOccurred())
	})
})

var _ = Describe("ConfigPlugins func", func() {
	var (
		mockClient *jenkins.MockClient
		s          *mockPluginsConfig
		f          *mockPluginsConfig
	)
	When("cascConfig is valid", func() {
		BeforeEach(func() {
			mockClient = &jenkins.MockClient{}
			s = &mockPluginsConfig{}
		})
		It("should work normal", func() {
			err := ConfigPlugins(mockClient, []PluginConfigAPI{s})
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})
	When("cascConfig is not valid", func() {
		BeforeEach(func() {
			mockClient = &jenkins.MockClient{}
			f = &mockPluginsConfig{
				configErr: fmt.Errorf("test error"),
			}
		})
		It("should return error", func() {
			err := ConfigPlugins(mockClient, []PluginConfigAPI{f})
			Expect(err).Error().Should(HaveOccurred())
		})
	})
})

var _ = Describe("GetRepoCredentialsId func", func() {
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
			sshKey := GetRepoCredentialsId(repoInfo)
			Expect(sshKey).Should(Equal(sshKeyCredentialName))
		})
	})
	When("repo type is github", func() {
		BeforeEach(func() {
			repoInfo.RepoType = "github"
		})
		It("should return github ssh key", func() {
			sshKey := GetRepoCredentialsId(repoInfo)
			Expect(sshKey).Should(Equal(githubCredentialName))
		})
	})
	When("repo type is not valid", func() {
		BeforeEach(func() {
			repoInfo.RepoType = "not exist"
		})
		It("should return empty", func() {
			sshKey := GetRepoCredentialsId(repoInfo)
			Expect(sshKey).Should(Equal(""))
		})
	})
})
