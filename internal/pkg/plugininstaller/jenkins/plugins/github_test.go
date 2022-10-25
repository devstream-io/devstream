package plugins

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("GithubJenkinsConfig", func() {
	var (
		c          *GithubJenkinsConfig
		mockClient *jenkins.MockClient
		repoOwner  string
	)
	BeforeEach(func() {
		repoOwner = "test_user"
		c = &GithubJenkinsConfig{
			RepoOwner: repoOwner,
		}
	})

	Context("GetDependentPlugins method", func() {
		It("should return github plugin", func() {
			plugins := c.getDependentPlugins()
			Expect(len(plugins)).Should(Equal(1))
			plugin := plugins[0]
			Expect(plugin.Name).Should(Equal("github-branch-source"))
		})
	})

	Context("config method", func() {
		When("create password failed", func() {
			var (
				createErrMsg string
			)
			BeforeEach(func() {
				createErrMsg = "create password failed"
				mockClient = &jenkins.MockClient{
					CreatePasswordCredentialError: fmt.Errorf(createErrMsg),
				}
			})
			It("should return error", func() {
				_, err := c.config(mockClient)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(createErrMsg))
			})
		})
		When("create password success", func() {
			BeforeEach(func() {
				mockClient = &jenkins.MockClient{}
			})
			It("should return repoCascConfig", func() {
				cascConfig, err := c.config(mockClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(cascConfig.RepoType).Should(Equal("github"))
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

var _ = Describe("NewGithubPlugin func", func() {
	var (
		pluginConfig *PluginGlobalConfig
	)
	BeforeEach(func() {
		pluginConfig = &PluginGlobalConfig{
			RepoInfo: &git.RepoInfo{
				Owner: "test",
			},
		}
	})
	It("should return github plugin", func() {
		githubPlugin := NewGithubPlugin(pluginConfig)
		Expect(githubPlugin.RepoOwner).Should(Equal("test"))
	})
})
