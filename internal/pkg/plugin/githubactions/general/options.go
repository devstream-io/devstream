package general

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/step"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

// actionOptions is the struct for configurations of the github-actions plugin.
type actionOptions struct {
	SCM    scm.SCMInfo `mapstructure:"scm"`
	Action action      `mapstructure:"action"`

	// used in package
	CIConfig    *cifile.CIConfig `mapstructure:"ci"`
	ProjectRepo *git.RepoInfo    `mapstructure:"projectRepo"`
}

type action struct {
	ConfigLocation string                    `mapstructure:"configLocation" validate:"required"`
	ImageRepo      *step.ImageRepoStepConfig `mapstructure:"imageRepo"`
	Dingtalk       *step.DingtalkStepConfig  `mapstructure:"dingTalk"`
	Sonarqube      *step.SonarQubeStepConfig `mapstructure:"sonarqube"`
}

// newOptions create options by raw options
func newActionOptions(options configmanager.RawOptions) (actionOptions, error) {
	var opts actionOptions
	if err := mapstructure.Decode(options, &opts); err != nil {
		return opts, err
	}
	return opts, nil
}

func (a *action) buildCIConfig(repoInfo *git.RepoInfo) *cifile.CIConfig {
	ciConfig := &cifile.CIConfig{
		Type:           "github",
		ConfigLocation: a.ConfigLocation,
	}
	// update ci render variables by plugins
	rawConfigVars := step.GenerateCIFileVars(a, repoInfo)
	rawConfigVars.Set("AppName", repoInfo.Repo)
	ciConfig.Vars = rawConfigVars
	log.Debugf("github-actions pipeline get render vars: %+v", ciConfig.Vars)
	return ciConfig
}
