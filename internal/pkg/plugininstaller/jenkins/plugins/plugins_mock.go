package plugins

import (
	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

type mockPluginsConfig struct {
	configVar *jenkins.RepoCascConfig
	configErr error
}

func (m *mockPluginsConfig) getDependentPlugins() (plugins []*jenkins.JenkinsPlugin) {
	return
}
func (m *mockPluginsConfig) config(jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
	if m.configErr != nil {
		return nil, m.configErr
	}
	return m.configVar, nil
}
func (m *mockPluginsConfig) setRenderVars(*jenkins.JenkinsFileRenderInfo) {
}
