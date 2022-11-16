package ci

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/ci/cifile/server"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/ci/step"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

type PipelineConfig struct {
	ConfigLocation string                    `mapstructure:"configLocation" validate:"required"`
	ImageRepo      *step.ImageRepoStepConfig `mapstructure:"imageRepo"`
	Dingtalk       *step.DingtalkStepConfig  `mapstructure:"dingTalk"`
	Sonarqube      *step.SonarQubeStepConfig `mapstructure:"sonarqube"`
	General        *step.GeneralStepConfig   `mapstructure:"general"`
}

type CIConfig struct {
	SCM      scm.SCMInfo    `mapstructure:"scm"`
	Pipeline PipelineConfig `mapstructure:"pipeline"`

	// used in package
	CIFileConfig *cifile.CIFileConfig `mapstructure:"ci"`
	ProjectRepo  *git.RepoInfo        `mapstructure:"projectRepo"`
}

func NewCIOptions(options configmanager.RawOptions) (*CIConfig, error) {
	var opts CIConfig
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

func (p *PipelineConfig) BuildCIFileConfig(ciType server.CIServerType, repoInfo *git.RepoInfo) *cifile.CIFileConfig {
	CIFileConfig := &cifile.CIFileConfig{
		Type:           ciType,
		ConfigLocation: p.ConfigLocation,
	}
	// update ci render variables by plugins
	rawConfigVars := step.GenerateCIFileVars(p, repoInfo)
	rawConfigVars.Set("AppName", repoInfo.Repo)
	CIFileConfig.Vars = rawConfigVars
	log.Debugf("gitlab-ci pipeline get render vars: %+v", CIFileConfig)
	return CIFileConfig
}
