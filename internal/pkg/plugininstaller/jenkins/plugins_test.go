package jenkins

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

type mockSuccessPluginsConfig struct{}

func (m *mockSuccessPluginsConfig) GetDependentPlugins() (plugins []*jenkins.JenkinsPlugin) {
	return
}
func (m *mockSuccessPluginsConfig) PreConfig(jenkins.JenkinsAPI) (config *jenkins.RepoCascConfig, err error) {
	return
}
func (m *mockSuccessPluginsConfig) UpdateJenkinsFileRenderVars(*jenkins.JenkinsFileRenderInfo) {}

var _ = Describe("installPlugins func", func() {
	var (
		mockClient *mockSuccessJenkinsClient
		s          *mockSuccessPluginsConfig
	)
	BeforeEach(func() {
		mockClient = &mockSuccessJenkinsClient{}
		s = &mockSuccessPluginsConfig{}
	})
	It("should work normal", func() {
		err := installPlugins(mockClient, []pluginConfigAPI{s}, false)
		Expect(err).Error().ShouldNot(HaveOccurred())
	})
})

var _ = Describe("preConfigPlugins func", func() {
	var (
		mockClient *mockSuccessJenkinsClient
		s          *mockSuccessPluginsConfig
	)
	BeforeEach(func() {
		mockClient = &mockSuccessJenkinsClient{}
		s = &mockSuccessPluginsConfig{}
	})
	It("should work normal", func() {
		err := preConfigPlugins(mockClient, []pluginConfigAPI{s})
		Expect(err).Error().ShouldNot(HaveOccurred())
	})

})
