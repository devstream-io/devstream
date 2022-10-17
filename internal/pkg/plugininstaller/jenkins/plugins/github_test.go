package plugins_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/jenkins/plugins"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

var _ = Describe("GithubJenkinsConfig", func() {
	var (
		c                     *plugins.GithubJenkinsConfig
		mockClient            *jenkins.MockClient
		jenkinsURL, repoOwner string
	)
	BeforeEach(func() {
		jenkinsURL = "jenkins_test"
		repoOwner = "test_user"
		c = &plugins.GithubJenkinsConfig{
			JenkinsURL: jenkinsURL,
			RepoOwner:  repoOwner,
		}
	})

	Context("GetDependentPlugins method", func() {
		It("should return github plugin", func() {
			plugins := c.GetDependentPlugins()
			Expect(len(plugins)).Should(Equal(1))
			plugin := plugins[0]
			Expect(plugin.Name).Should(Equal("github-branch-source"))
		})
	})

	Context("PreConfig method", func() {
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
				_, err := c.PreConfig(mockClient)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(createErrMsg))
			})
		})
		When("create password success", func() {
			BeforeEach(func() {
				mockClient = &jenkins.MockClient{}
			})
			It("should return repoCascConfig", func() {
				cascConfig, err := c.PreConfig(mockClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(cascConfig.RepoType).Should(Equal("github"))
				Expect(cascConfig.CredentialID).Should(Equal(plugins.GithubCredentialName))
				Expect(cascConfig.JenkinsURL).Should(Equal(jenkinsURL))
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
