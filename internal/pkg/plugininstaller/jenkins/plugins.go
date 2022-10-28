package jenkins

import (
	"github.com/imdario/mergo"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/step"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
)

func ensurePluginInstalled(jenkinsClient jenkins.JenkinsAPI, pluginConfigs []step.StepConfigAPI) error {
	var plugins []*jenkins.JenkinsPlugin
	for _, pluginConfig := range pluginConfigs {
		plugins = append(plugins, pluginConfig.GetJenkinsPlugins()...)
	}
	return jenkinsClient.InstallPluginsIfNotExists(plugins)
}

func configPlugins(jenkinsClient jenkins.JenkinsAPI, pluginConfigs []step.StepConfigAPI) error {
	globalCascConfig := new(jenkins.RepoCascConfig)
	for _, pluginConfig := range pluginConfigs {
		cascConfig, err := pluginConfig.ConfigJenkins(jenkinsClient)
		if err != nil {
			log.Debugf("jenkins plugin %+v config error", pluginConfig)
			return err
		}
		if cascConfig != nil {
			err := mergo.Merge(globalCascConfig, cascConfig, mergo.WithOverride)
			if err != nil {
				log.Debugf("jenins merge casc config failed: %+v", err)
				return err
			}
		}
	}
	return jenkinsClient.ConfigCascForRepo(globalCascConfig)
}
