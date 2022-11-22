package ci

import (
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile/server"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/step"
	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

type PipelineConfig struct {
	ConfigLocation downloader.ResourceLocation `mapstructure:"configLocation" validate:"required"`
	ImageRepo      *step.ImageRepoStepConfig   `mapstructure:"imageRepo"`
	Dingtalk       *step.DingtalkStepConfig    `mapstructure:"dingTalk"`
	Sonarqube      *step.SonarQubeStepConfig   `mapstructure:"sonarqube"`
	General        *step.GeneralStepConfig     `mapstructure:"general"`
}

type CIConfig struct {
	SCM      scm.SCMInfo    `mapstructure:"scm"  validate:"required"`
	Pipeline PipelineConfig `mapstructure:"pipeline"  validate:"required"`

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
		ConfigLocation: downloader.ResourceLocation(p.ConfigLocation),
	}
	// update ci render variables by plugins
	rawConfigVars := p.generateCIFileVars(repoInfo)
	rawConfigVars.Set("AppName", repoInfo.Repo)
	CIFileConfig.Vars = rawConfigVars
	log.Debugf("gitlab-ci pipeline get render vars: %+v", CIFileConfig)
	return CIFileConfig
}

func (p *PipelineConfig) generateCIFileVars(repoInfo *git.RepoInfo) cifile.CIFileVarsMap {
	// set default command for language
	if p.General != nil {
		p.General.SetDefault()
	}
	varMap, _ := mapz.DecodeStructToMap(p)
	globalVarsMap, _ := mapz.DecodeStructToMap(
		step.GetStepGlobalVars(repoInfo),
	)
	err := mergo.Merge(&varMap, globalVarsMap)
	if err != nil {
		log.Warnf("cifile merge CIFileVarsMap failed: %+v", err)
	}
	return varMap
}
