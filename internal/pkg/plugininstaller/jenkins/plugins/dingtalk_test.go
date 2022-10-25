package plugins

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/base"
	"github.com/devstream-io/devstream/pkg/util/jenkins"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DingtalkJenkinsConfig", func() {
	var (
		c          *DingtalkJenkinsConfig
		mockClient *jenkins.MockClient
		name       string
	)
	BeforeEach(func() {
		name = "test"
		c = &DingtalkJenkinsConfig{
			base.DingtalkStepConfig{
				Name: name,
			},
		}
	})
	Context("GetDependentPlugins method", func() {
		It("should return dingding plugin", func() {
			plugins := c.getDependentPlugins()
			Expect(len(plugins)).Should(Equal(1))
			plugin := plugins[0]
			Expect(plugin.Name).Should(Equal("dingding-notifications"))
		})
	})

	Context("config method", func() {
		When("all config work noraml", func() {
			BeforeEach(func() {
				mockClient = &jenkins.MockClient{}
			})
			It("should return nil", func() {
				_, err := c.config(mockClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})

	Context("setRenderVars method", func() {
		var (
			renderInfo *jenkins.JenkinsFileRenderInfo
			atUsers    string
		)
		BeforeEach(func() {
			atUsers = "testUser"
			renderInfo = &jenkins.JenkinsFileRenderInfo{}
			c.AtUsers = atUsers
		})
		It("should update renderInfo with DingtalkRobotID and DingtalkAtUser", func() {
			c.setRenderVars(renderInfo)
			Expect(renderInfo.DingtalkAtUser).Should(Equal(atUsers))
			Expect(renderInfo.DingtalkRobotID).Should(Equal(name))
		})
	})
})
