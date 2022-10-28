package step_test

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/step"
	"github.com/devstream-io/devstream/pkg/util/jenkins"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DingtalkStepConfig", func() {
	var (
		c          *step.DingtalkStepConfig
		mockClient *jenkins.MockClient
		name       string
	)
	BeforeEach(func() {
		name = "test"
		c = &step.DingtalkStepConfig{
			Name: name,
		}
	})
	Context("GetJenkinsPlugins method", func() {
		It("should return dingding plugin", func() {
			plugins := c.GetJenkinsPlugins()
			Expect(len(plugins)).Should(Equal(1))
			plugin := plugins[0]
			Expect(plugin.Name).Should(Equal("dingding-notifications"))
		})
	})

	Context("ConfigJenkins method", func() {
		When("all config work noraml", func() {
			BeforeEach(func() {
				mockClient = &jenkins.MockClient{}
			})
			It("should return nil", func() {
				_, err := c.ConfigJenkins(mockClient)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
		})
	})

})
