package step

import (
	"github.com/imdario/mergo"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/types"
)

type test struct {
	Enable                *bool  `mapstructure:"enable"`
	Command               string `mapstructure:"command"`
	ContainerName         string `mapstructure:"containerName"`
	CoverageCommand       string `mapstructure:"coverageCommand"`
	CoverageStatusCommand string `mapstructure:"CoverageStatusCommand"`
}

type GeneralStepConfig struct {
	Language *Language `mapstructure:"language"`
	Test     *test     `mapstructure:"test"`
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

func (g *GeneralStepConfig) SetDefault() {
	if g.Language != nil && g.Language.Name == "" {
		return
	}
	generalDefaultOpts := g.Language.getGeneralDefaultOption()
	if g.Test.Enable != types.Bool(false) {
		if err := mergo.Merge(g.Test, generalDefaultOpts); err != nil {
			log.Warnf("step general merge default config failed: %+v", err)
		}
	}
}
