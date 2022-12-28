package ci

import (
	"github.com/imdario/mergo"
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile/server"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/config"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/step"
	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/types"
)

type PipelineConfig struct {
	ConfigLocation downloader.ResourceLocation `mapstructure:"configLocation" validate:"required"`
	ImageRepo      *step.ImageRepoStepConfig   `mapstructure:"imageRepo"`
	Dingtalk       *step.DingtalkStepConfig    `mapstructure:"dingTalk"`
	Sonarqube      *step.SonarQubeStepConfig   `mapstructure:"sonarqube"`
	Lanuage        config.LanguageOption       `mapstructure:"language"`
	Test           config.TestOption           `mapstructure:"test"`
}

type CIConfig struct {
	ProjectRepo *git.RepoInfo   `mapstructure:"scm"  validate:"required"`
	Pipeline    *PipelineConfig `mapstructure:"pipeline"  validate:"required"`

	// used in package
	CIFileConfig *cifile.CIFileConfig `mapstructure:"ci"`
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
	rawConfigVars := p.GenerateCIFileVars(repoInfo)
	CIFileConfig.Vars = rawConfigVars
	log.Debugf("gitlab-ci pipeline get render vars: %+v", CIFileConfig)
	return CIFileConfig
}

func (p *PipelineConfig) GenerateCIFileVars(repoInfo *git.RepoInfo) cifile.CIFileVarsMap {
	// set default command for language
	p.setDefault()
	varMap, _ := mapz.DecodeStructToMap(p)
	globalVarsMap, _ := mapz.DecodeStructToMap(
		step.GetStepGlobalVars(repoInfo),
	)
	err := mergo.Merge(&varMap, globalVarsMap)
	if err != nil {
		log.Warnf("cifile merge CIFileVarsMap failed: %+v", err)
	}
	varMap["AppName"] = repoInfo.Repo
	return varMap
}

func (p *PipelineConfig) setDefault() {
	if !p.Lanuage.IsConfigured() {
		return
	}
	// get language default options
	defaultOpt := p.Lanuage.GetGeneralDefaultOpt()
	if defaultOpt == nil {
		log.Debugf("pipeline language [%+v] done's have default options", p.Lanuage)
		return
	}
	if p.Test.Enable != types.Bool(false) {
		if err := mergo.Merge(&p.Test, defaultOpt.Test); err != nil {
			log.Warnf("ci merge default config failed: %+v", err)
			return
		}
	}
	// set image default url
	if p.ImageRepo != nil {
		p.ImageRepo.URL = p.ImageRepo.GetImageRepoURL()
	}
}
