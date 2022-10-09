package jenkins

import (
	"github.com/imdario/mergo"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type pluginConfigAPI interface {
	GetDependentPlugins() []*jenkins.JenkinsPlugin
	PreConfig(jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error)
	UpdateJenkinsFileRenderVars(*jenkins.JenkinsFileRenderInfo)
}

func installPlugins(jenkinsClient jenkins.JenkinsAPI, pluginConfigs []pluginConfigAPI, enableRestart bool) error {
	var plugins []*jenkins.JenkinsPlugin
	for _, pluginConfig := range pluginConfigs {
		plugins = append(plugins, pluginConfig.GetDependentPlugins()...)
	}
	return jenkinsClient.InstallPluginsIfNotExists(plugins, enableRestart)
}

func preConfigPlugins(jenkinsClient jenkins.JenkinsAPI, pluginConfigs []pluginConfigAPI) error {
	globalCascConfig := new(jenkins.RepoCascConfig)
	for _, pluginConfig := range pluginConfigs {
		cascConfig, err := pluginConfig.PreConfig(jenkinsClient)
		if err != nil {
			log.Debugf("jenkins plugin %+v preConfig error", pluginConfig)
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
