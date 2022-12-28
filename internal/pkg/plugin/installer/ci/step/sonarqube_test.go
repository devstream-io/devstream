package step_test

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/step"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/scm"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SonarQubeStepConfig", func() {
	var (
		c                *step.SonarQubeStepConfig
		mockClient       *jenkins.MockClient
		name, url, token string
	)
	BeforeEach(func() {
		name = "test"
		token = "test_token"
		url = "test_url"
		c = &step.SonarQubeStepConfig{
			Name:  name,
			Token: token,
			URL:   url,
		}
	})
	Context("GetJenkinsPlugins method", func() {
		It("should return sonar plugin", func() {
			plugins := c.GetJenkinsPlugins()
			Expect(len(plugins)).Should(Equal(1))
			plugin := plugins[0]
			Expect(plugin.Name).Should(Equal("sonar"))
		})
	})

	Context("ConfigJenkins method", func() {
		When("Create secret failed", func() {
			var (
				errMsg string
			)
			BeforeEach(func() {
				errMsg = "create secret failed"
				mockClient = &jenkins.MockClient{
					CreateSecretCredentialError: fmt.Errorf(errMsg),
				}
			})
			It("should return error", func() {
				_, err := c.ConfigJenkins(mockClient)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err).Error().Should(HaveOccurred())
			})
		})
		When("all config work noraml", func() {
			BeforeEach(func() {
				mockClient = &jenkins.MockClient{}
			})
			It("should return nil", func() {
				cascConfig, err := c.ConfigJenkins(mockClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(cascConfig.SonarqubeURL).Should(Equal(url))
				Expect(cascConfig.SonarqubeName).Should(Equal(name))
				Expect(cascConfig.SonarTokenCredentialID).Should(Equal("SONAR_SECRET_TOKEN"))
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
