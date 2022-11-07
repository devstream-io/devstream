package step_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/step"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

var _ = Describe("GeneralStepConfig struct", func() {
	var (
		c          *step.GeneralStepConfig
		mockClient *jenkins.MockClient
	)
	BeforeEach(func() {
		c = &step.GeneralStepConfig{
			Language: "test",
		}
	})
	Context("GetJenkinsPlugins method", func() {
		It("should return empty plugins", func() {
			plugins := c.GetJenkinsPlugins()
			Expect(len(plugins)).Should(Equal(0))
		})
	})
	Context("ConfigJenkins method", func() {
		It("should return nil", func() {
			casc, err := c.ConfigJenkins(mockClient)
			Expect(err).Should(BeNil())
			Expect(casc).Should(BeNil())
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
