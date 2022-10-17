package plugins_test

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/jenkins/plugins"
	"github.com/devstream-io/devstream/pkg/util/jenkins"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("SonarQubeJenkinsConfig", func() {
	var (
		c                *plugins.SonarQubeJenkinsConfig
		mockClient       *jenkins.MockClient
		name, url, token string
	)
	BeforeEach(func() {
		name = "test"
		token = "test_token"
		url = "test_url"
		c = &plugins.SonarQubeJenkinsConfig{
			Name:  name,
			Token: token,
			URL:   url,
		}
	})
	Context("GetDependentPlugins method", func() {
		It("should return sonar plugin", func() {
			plugins := c.GetDependentPlugins()
			Expect(len(plugins)).Should(Equal(1))
			plugin := plugins[0]
			Expect(plugin.Name).Should(Equal("sonar"))
		})
	})

	Context("PreConfig method", func() {
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
				_, err := c.PreConfig(mockClient)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err).Error().Should(HaveOccurred())
			})
		})
		When("all config work noraml", func() {
			BeforeEach(func() {
				mockClient = &jenkins.MockClient{}
			})
			It("should return nil", func() {
				cascConfig, err := c.PreConfig(mockClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(cascConfig.SonarqubeURL).Should(Equal(url))
				Expect(cascConfig.SonarqubeName).Should(Equal(name))
				Expect(cascConfig.SonarTokenCredentialID).Should(Equal("sonarqubeTokenCredential"))
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
		It("should update renderInfo with DingtalkRobotID and DingtalkAtUser", func() {
			c.UpdateJenkinsFileRenderVars(renderInfo)
			Expect(renderInfo.SonarqubeEnable).Should(BeTrue())
		})
	})
})
