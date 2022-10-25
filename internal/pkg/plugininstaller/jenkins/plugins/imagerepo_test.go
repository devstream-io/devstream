package plugins

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/base"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

var _ = Describe("ImageRepoJenkinsConfig", func() {
	var (
		c         *ImageRepoJenkinsConfig
		url, user string
	)
	BeforeEach(func() {
		url = "jenkins_test"
		user = "test_user"
		c = &ImageRepoJenkinsConfig{
			base.ImageRepoStepConfig{
				URL:  url,
				User: user,
			},
		}
	})

	Context("getDependentPlugins method", func() {
		It("should be empty", func() {
			plugins := c.getDependentPlugins()
			Expect(len(plugins)).Should(BeZero())
		})
	})

	Context("generateDockerAuthSecretData method", func() {
		var (
			existImageEnv   string
			testImageEnvVar string
		)
		BeforeEach(func() {
			existImageEnv = os.Getenv("IMAGE_REPO_PASSWORD")
			if existImageEnv != "" {
				err := os.Unsetenv("IMAGE_REPO_PASSWORD")
				Expect(err).Error().ShouldNot(HaveOccurred())
			}
		})
		When("image env is not set", func() {
			It("should return error", func() {
				_, err := c.generateDockerAuthSecretData()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("the environment variable IMAGE_REPO_PASSWORD is not set"))
			})
		})
		When("image env is set", func() {
			BeforeEach(func() {
				testImageEnvVar = "test_image_repo_env"
				err := os.Setenv("IMAGE_REPO_PASSWORD", testImageEnvVar)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
			It("should return image auth data", func() {
				d, err := c.generateDockerAuthSecretData()
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(len(d)).Should(Equal(1))
				configJson, ok := d["config.json"]
				Expect(ok).Should(BeTrue())

				expectStr := fmt.Sprintf(`{
  "auths": {
    "%s": {
      "auth": "%s"
    }
  }
}`, c.URL, "dGVzdF91c2VyOnRlc3RfaW1hZ2VfcmVwb19lbnY=")
				Expect(string(configJson)).Should(Equal(expectStr))
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
		It("should update image info", func() {
			c.setRenderVars(renderInfo)
			Expect(renderInfo.ImageRepositoryURL).Should(Equal(fmt.Sprintf("%s/library", url)))
			Expect(renderInfo.ImageAuthSecretName).Should(Equal(imageRepoSecretName))
		})
	})
})
