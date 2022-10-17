package plugins_test

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/jenkins/plugins"
	"github.com/devstream-io/devstream/pkg/util/jenkins"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DingtalkJenkinsConfig", func() {
	var (
		c             *plugins.DingtalkJenkinsConfig
		mockClient    *jenkins.MockClient
		name, atUsers string
	)
	BeforeEach(func() {
		name = "test"
		c = &plugins.DingtalkJenkinsConfig{
			Name: name,
		}
	})
	Context("GetDependentPlugins method", func() {
		It("should return dingding plugin", func() {
			plugins := c.GetDependentPlugins()
			Expect(len(plugins)).Should(Equal(1))
			plugin := plugins[0]
			Expect(plugin.Name).Should(Equal("dingding-notifications"))
		})
	})

	Context("PreConfig method", func() {
		When("all config work noraml", func() {
			BeforeEach(func() {
				mockClient = &jenkins.MockClient{}
			})
			It("should return nil", func() {
				_, err := c.PreConfig(mockClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})

	Context("UpdateJenkinsFileRenderVars method", func() {
		var (
			renderInfo *jenkins.JenkinsFileRenderInfo
		)
		BeforeEach(func() {
			atUsers = "testUser"
			renderInfo = &jenkins.JenkinsFileRenderInfo{}
			c.AtUsers = atUsers
		})
		It("should update renderInfo with DingtalkRobotID and DingtalkAtUser", func() {
			c.UpdateJenkinsFileRenderVars(renderInfo)
			Expect(renderInfo.DingtalkAtUser).Should(Equal(atUsers))
			Expect(renderInfo.DingtalkRobotID).Should(Equal(name))
		})
	})
})
