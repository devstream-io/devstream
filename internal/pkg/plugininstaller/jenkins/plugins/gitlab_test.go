package plugins_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/jenkins/plugins"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

var _ = Describe("GitlabJenkinsConfig", func() {
	var (
		c                          *plugins.GitlabJenkinsConfig
		mockClient                 *jenkins.MockClient
		sshKey, repoOwner, baseURL string
	)
	BeforeEach(func() {
		baseURL = "jenkins_test"
		repoOwner = "test_user"
		sshKey = "test_key"
		c = &plugins.GitlabJenkinsConfig{
			BaseURL:   baseURL,
			RepoOwner: repoOwner,
		}
	})

	Context("GetDependentPlugins method", func() {
		It("should return gitlab plugin", func() {
			plugins := c.GetDependentPlugins()
			Expect(len(plugins)).Should(Equal(1))
			plugin := plugins[0]
			Expect(plugin.Name).Should(Equal("gitlab-plugin"))
		})
	})

	Context("PreConfig method", func() {
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
				_, err := c.PreConfig(mockClient)
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
					_, err := c.PreConfig(mockClient)
					Expect(err).Error().Should(HaveOccurred())
					Expect(err.Error()).Should(Equal(createErrMsg))
				})
			})
			When("create gitlabCredential success", func() {
				It("should return repoCascConfig", func() {
					cascConfig, err := c.PreConfig(mockClient)
					Expect(err).Error().ShouldNot(HaveOccurred())
					Expect(cascConfig.RepoType).Should(Equal("gitlab"))
					Expect(cascConfig.CredentialID).Should(Equal("gitlabCredential"))
					Expect(cascConfig.GitLabConnectionName).Should(Equal(plugins.GitlabConnectionName))
					Expect(cascConfig.GitlabURL).Should(Equal(baseURL))
				})
			})
		})
	})

	Context("UpdateJenkinsFileRenderVars method", func() {
		var (
			renderInfo *jenkins.JenkinsFileRenderInfo
		)
		BeforeEach(func() {
			renderInfo = &jenkins.JenkinsFileRenderInfo{}
		})
		It("should update nothing", func() {
			c.UpdateJenkinsFileRenderVars(renderInfo)
			Expect(*renderInfo).Should(Equal(jenkins.JenkinsFileRenderInfo{}))
		})
	})
})
