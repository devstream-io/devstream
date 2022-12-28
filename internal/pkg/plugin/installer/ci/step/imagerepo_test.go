package step

import (
	"fmt"

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
		When("image env is not set", func() {
			BeforeEach(func() {
				c = &ImageRepoStepConfig{
					URL:  url,
					User: user,
				}
			})
			It("should return error", func() {
				_, err := c.generateDockerAuthSecretData()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("config field password is not set"))
			})
		})
		When("image env is set", func() {
			BeforeEach(func() {
				c = &ImageRepoStepConfig{
					URL:      url,
					User:     user,
					Password: "test",
				}
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
}`, c.URL, "dGVzdF91c2VyOnRlc3Q=")
				Expect(string(configJson)).Should(Equal(expectStr))
			})
		})
	})

	Context("ConfigSCM method", func() {
		var (
			scmClient *scm.MockScmClient
		)
		When("imageRepoPassword is not valid", func() {
			BeforeEach(func() {
				c = &ImageRepoStepConfig{
					URL:  url,
					User: user,
				}
				scmClient = &scm.MockScmClient{}
			})
			It("should return error", func() {
				err := c.ConfigSCM(scmClient)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("config field password is not set"))
			})
		})
		When("all valid", func() {
			BeforeEach(func() {
				scmClient = &scm.MockScmClient{}
				c = &ImageRepoStepConfig{
					URL:      url,
					User:     user,
					Password: "test",
				}
			})
			It("should return nil", func() {
				err := c.ConfigSCM(scmClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})
})
