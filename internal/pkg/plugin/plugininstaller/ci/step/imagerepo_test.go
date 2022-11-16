package step

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm"
)

var _ = Describe("ImageRepoStepConfig", func() {
	var (
		c         *ImageRepoStepConfig
		url, user string
	)
	BeforeEach(func() {
		url = "jenkins_test"
		user = "test_user"
		c = &ImageRepoStepConfig{
			URL:  url,
			User: user,
		}
	})

	Context("GetJenkinsPlugins method", func() {
		It("should be empty", func() {
			plugins := c.GetJenkinsPlugins()
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

	Context("ConfigSCM method", func() {
		var (
			scmClient             *scm.MockScmClient
			errMsg, existImageEnv string
		)
		BeforeEach(func() {
			existImageEnv = os.Getenv("IMAGE_REPO_PASSWORD")
			if existImageEnv != "" {
				err := os.Unsetenv("IMAGE_REPO_PASSWORD")
				Expect(err).Error().ShouldNot(HaveOccurred())
			}
		})
		When("imageRepoPassword is not valid", func() {
			BeforeEach(func() {
				errMsg = "the environment variable IMAGE_REPO_PASSWORD is not set"
				scmClient = &scm.MockScmClient{}
			})
			It("should return error", func() {
				err := c.ConfigSCM(scmClient)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(errMsg))
			})
		})
		When("all valid", func() {
			BeforeEach(func() {
				scmClient = &scm.MockScmClient{}
				os.Setenv("IMAGE_REPO_PASSWORD", "test")
			})
			It("should return nil", func() {
				err := c.ConfigSCM(scmClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
		AfterEach(func() {
			os.Setenv("IMAGE_REPO_PASSWORD", existImageEnv)
		})
	})

})
