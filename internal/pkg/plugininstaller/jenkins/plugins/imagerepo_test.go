package plugins_test

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/jenkins/plugins"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

var _ = Describe("ImageRepoJenkinsConfig", func() {
	var (
		c                        *plugins.ImageRepoJenkinsConfig
		mockClient               *jenkins.MockClient
		authNamespace, url, user string
	)
	BeforeEach(func() {
		url = "jenkins_test"
		user = "test_user"
		authNamespace = "test_namespace"
		c = &plugins.ImageRepoJenkinsConfig{
			URL:           url,
			User:          user,
			AuthNamespace: authNamespace,
		}
	})

	Context("GetDependentPlugins method", func() {
		It("should be empty", func() {
			plugins := c.GetDependentPlugins()
			Expect(len(plugins)).Should(BeZero())
		})
	})

	Context("PreConfig method", func() {
		When("image env is not set", func() {
			BeforeEach(func() {
				err := os.Unsetenv("IMAGE_REPO_PASSWORD")
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
			It("should return error", func() {
				_, err := c.PreConfig(mockClient)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("the environment variable IMAGE_REPO_PASSWORD is not set"))
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
			Expect(renderInfo.ImageRepositoryURL).Should(Equal(fmt.Sprintf("%s/library", url)))
			Expect(renderInfo.ImageAuthSecretName).Should(Equal("docker-config"))
		})
	})
})
