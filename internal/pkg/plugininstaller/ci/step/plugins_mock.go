package step

import (
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

type MockPluginsConfig struct {
	ConfigVar *jenkins.RepoCascConfig
	ConfigErr error
}

func (m *MockPluginsConfig) GetJenkinsPlugins() (plugins []*jenkins.JenkinsPlugin) {
	return
}
func (m *MockPluginsConfig) ConfigJenkins(jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
	if m.ConfigErr != nil {
		return nil, m.ConfigErr
	}
	return m.ConfigVar, nil
}
func (m *MockPluginsConfig) ConfigSCM(scm.ClientOperation) error {
	return nil
}
