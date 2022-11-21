package jenkinspipeline

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/step"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

var _ = Describe("ensurePluginInstalled func", func() {
	var (
		mockClient *jenkins.MockClient
		s          *step.MockPluginsConfig
	)
	BeforeEach(func() {
		mockClient = &jenkins.MockClient{}
		s = &step.MockPluginsConfig{}
	})
	It("should work normal", func() {
		err := ensurePluginInstalled(mockClient, []step.StepConfigAPI{s})
		Expect(err).Error().ShouldNot(HaveOccurred())
	})
})

var _ = Describe("configPlugins func", func() {
	var (
		mockClient *jenkins.MockClient
		s          *step.MockPluginsConfig
		f          *step.MockPluginsConfig
	)
	When("configPlugins is valid", func() {
		BeforeEach(func() {
			mockClient = &jenkins.MockClient{}
			s = &step.MockPluginsConfig{}
		})
		It("should work normal", func() {
			err := configPlugins(mockClient, []step.StepConfigAPI{s})
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})
	When("configPlugins is not valid", func() {
		BeforeEach(func() {
			mockClient = &jenkins.MockClient{}
			f = &step.MockPluginsConfig{
				ConfigErr: fmt.Errorf("test error"),
			}
		})
		It("should return error", func() {
			err := configPlugins(mockClient, []step.StepConfigAPI{f})
			Expect(err).Error().Should(HaveOccurred())
		})
	})
})
