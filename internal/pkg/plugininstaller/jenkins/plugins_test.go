package jenkins

import (
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

type mockSuccessPluginsConfig struct{}

func (m *mockSuccessPluginsConfig) GetDependentPlugins() (plugins []*jenkins.JenkinsPlugin) {
	return
}
func (m *mockSuccessPluginsConfig) PreConfig(jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
	return &jenkins.RepoCascConfig{
		RepoType:     "gitlab",
		CredentialID: "test",
	}, nil
}
func (m *mockSuccessPluginsConfig) UpdateJenkinsFileRenderVars(*jenkins.JenkinsFileRenderInfo) {}

type mockErrorPluginsConfig struct{}

func (m *mockErrorPluginsConfig) GetDependentPlugins() (plugins []*jenkins.JenkinsPlugin) {
	return
}
func (m *mockErrorPluginsConfig) PreConfig(jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
	return nil, errors.New("test_error")
}
func (m *mockErrorPluginsConfig) UpdateJenkinsFileRenderVars(*jenkins.JenkinsFileRenderInfo) {}

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
		f          *mockErrorPluginsConfig
	)
	When("cascConfig is valid", func() {
		BeforeEach(func() {
			mockClient = &mockSuccessJenkinsClient{}
			s = &mockSuccessPluginsConfig{}
		})
		It("should work normal", func() {
			err := preConfigPlugins(mockClient, []pluginConfigAPI{s})
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})
	When("cascConfig is not valid", func() {
		BeforeEach(func() {
			mockClient = &mockSuccessJenkinsClient{}
			f = &mockErrorPluginsConfig{}
		})
		It("should return error", func() {
			err := preConfigPlugins(mockClient, []pluginConfigAPI{f})
			Expect(err).Error().Should(HaveOccurred())
		})
	})
})
