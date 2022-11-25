package step

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

var _ = Describe("GeneralStepConfig struct", func() {
	var (
		c          *GeneralStepConfig
		mockClient *jenkins.MockClient
	)
	BeforeEach(func() {
		c = &GeneralStepConfig{}
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
	Context("SetDefault method", func() {
		When("language is empty", func() {
			BeforeEach(func() {
				c.Language = language{}
			})
			It("should return GeneralStepConfig", func() {
				c.SetDefault()
				Expect(c.Test.Enable).Should(BeNil())
			})
		})
		When("language is valid", func() {
			BeforeEach(func() {
				c.Language = language{
					Name: "go",
				}
			})
			It("should set test options", func() {
				goTestOption, exist := languageDefaultOptionMap["go"]
				Expect(exist).Should(BeTrue())
				Expect(goTestOption).ShouldNot(BeNil())
				c.SetDefault()
				Expect(c.Test).Should(Equal(*goTestOption.testOption))
			})
		})
	})
})
