package step

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("GitlabStepConfig", func() {
	var (
		c                          *GitlabStepConfig
		mockClient                 *jenkins.MockClient
		sshKey, repoOwner, baseURL string
	)
	BeforeEach(func() {
		baseURL = "jenkins_test"
		repoOwner = "test_user"
		sshKey = "test_key"
		c = &GitlabStepConfig{
			BaseURL:   baseURL,
			RepoOwner: repoOwner,
		}
	})

	Context("GetJenkinsPlugins method", func() {
		It("should return gitlab plugin", func() {
			plugins := c.GetJenkinsPlugins()
			Expect(len(plugins)).Should(Equal(1))
			plugin := plugins[0]
			Expect(plugin.Name).Should(Equal("gitlab-plugin"))
		})
	})

	Context("ConfigJenkins method", func() {
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
				_, err := c.ConfigJenkins(mockClient)
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
					_, err := c.ConfigJenkins(mockClient)
					Expect(err).Error().Should(HaveOccurred())
					Expect(err.Error()).Should(Equal(createErrMsg))
				})
			})
			When("create gitlabCredential success", func() {
				It("should return repoCascConfig", func() {
					cascConfig, err := c.ConfigJenkins(mockClient)
					Expect(err).Error().ShouldNot(HaveOccurred())
					Expect(cascConfig.RepoType).Should(Equal("gitlab"))
					Expect(cascConfig.CredentialID).Should(Equal("gitlabCredential"))
					Expect(cascConfig.GitLabConnectionName).Should(Equal("gitlabConnection"))
					Expect(cascConfig.GitlabURL).Should(Equal(baseURL))
				})
			})
		})
	})
	Context("ConfigSCM method", func() {
		It("should return nil", func() {
			scmClient := &scm.MockScmClient{}
			err := c.ConfigSCM(scmClient)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})
})

var _ = Describe("newGitlabStep func", func() {
	var (
		pluginConfig           *StepGlobalOption
		sshKey, owner, baseURL string
	)
	BeforeEach(func() {
		owner = "test_owner"
		sshKey = "test_key"
		baseURL = "http://base.com"
		pluginConfig = &StepGlobalOption{
			RepoInfo: &git.RepoInfo{
				Owner:         owner,
				SSHPrivateKey: sshKey,
				BaseURL:       baseURL,
			},
		}
	})
	It("should return gitlab plugin", func() {
		githubPlugin := newGitlabStep(pluginConfig)
		Expect(githubPlugin.RepoOwner).Should(Equal(owner))
		Expect(githubPlugin.SSHPrivateKey).Should(Equal(sshKey))
		Expect(githubPlugin.BaseURL).Should(Equal(baseURL))
	})
})
