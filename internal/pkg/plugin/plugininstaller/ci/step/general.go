package step

import (
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

type GeneralStepConfig struct {
	Language string `mapstructure:"language"`
}

// GetJenkinsPlugins return jenkins plugins info
func (g *GeneralStepConfig) GetJenkinsPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{}
}

// JenkinsConfig config jenkins and return casc config
func (g *GeneralStepConfig) ConfigJenkins(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
	return nil, nil
}

func (g *GeneralStepConfig) ConfigSCM(client scm.ClientOperation) error {
	return nil
}
