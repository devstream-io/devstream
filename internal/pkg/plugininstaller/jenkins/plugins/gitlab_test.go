package plugins

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("GitlabJenkinsConfig", func() {
	var (
		c                          *GitlabJenkinsConfig
		mockClient                 *jenkins.MockClient
		sshKey, repoOwner, baseURL string
	)
	BeforeEach(func() {
		baseURL = "jenkins_test"
		repoOwner = "test_user"
		sshKey = "test_key"
		c = &GitlabJenkinsConfig{
			BaseURL:   baseURL,
			RepoOwner: repoOwner,
		}
	})

	Context("getDependentPlugins method", func() {
		It("should return gitlab plugin", func() {
			plugins := c.getDependentPlugins()
			Expect(len(plugins)).Should(Equal(1))
			plugin := plugins[0]
			Expect(plugin.Name).Should(Equal("gitlab-plugin"))
		})
	})

	Context("config method", func() {
		When("create sshKey failed", func() {
			var (
				createErrMsg string
			)
			BeforeEach(func() {
				createErrMsg = "create ssh key failed"
				mockClient = &jenkins.MockClient{
					CreateSSHKeyCredentialError: fmt.Errorf(createErrMsg),
				}
				c.SSHPrivateKey = sshKey
			})
			It("should return error", func() {
				_, err := c.config(mockClient)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(createErrMsg))
			})
		})
		When("not use ssh key", func() {
			BeforeEach(func() {
				mockClient = &jenkins.MockClient{}
			})
			When("create gitlabCredential failed", func() {
				var (
					createErrMsg string
				)
				BeforeEach(func() {
					createErrMsg = "create gitlabCredential failed"
					mockClient.CreateGiltabCredentialError = fmt.Errorf(createErrMsg)
				})
				It("should return error", func() {
					_, err := c.config(mockClient)
					Expect(err).Error().Should(HaveOccurred())
					Expect(err.Error()).Should(Equal(createErrMsg))
				})
			})
			When("create gitlabCredential success", func() {
				It("should return repoCascConfig", func() {
					cascConfig, err := c.config(mockClient)
					Expect(err).Error().ShouldNot(HaveOccurred())
					Expect(cascConfig.RepoType).Should(Equal("gitlab"))
					Expect(cascConfig.CredentialID).Should(Equal("gitlabCredential"))
					Expect(cascConfig.GitLabConnectionName).Should(Equal(GitlabConnectionName))
					Expect(cascConfig.GitlabURL).Should(Equal(baseURL))
				})
			})
		})
	})

	Context("setRenderVars method", func() {
		var (
			renderInfo *jenkins.JenkinsFileRenderInfo
		)
		BeforeEach(func() {
			renderInfo = &jenkins.JenkinsFileRenderInfo{}
		})
		It("should update nothing", func() {
			c.setRenderVars(renderInfo)
			Expect(*renderInfo).Should(Equal(jenkins.JenkinsFileRenderInfo{}))
		})
	})
})

var _ = Describe("NewGitlabPlugin func", func() {
	var (
		pluginConfig           *PluginGlobalConfig
		sshKey, owner, baseURL string
	)
	BeforeEach(func() {
		owner = "test_owner"
		sshKey = "test_key"
		baseURL = "http://base.com"
		pluginConfig = &PluginGlobalConfig{
			RepoInfo: &git.RepoInfo{
				Owner:         owner,
				SSHPrivateKey: sshKey,
				BaseURL:       baseURL,
			},
		}
	})
	It("should return gitlab plugin", func() {
		githubPlugin := NewGitlabPlugin(pluginConfig)
		Expect(githubPlugin.RepoOwner).Should(Equal(owner))
		Expect(githubPlugin.SSHPrivateKey).Should(Equal(sshKey))
		Expect(githubPlugin.BaseURL).Should(Equal(baseURL))
	})
})
