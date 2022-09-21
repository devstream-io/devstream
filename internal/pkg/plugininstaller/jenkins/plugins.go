package jenkins

import (
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type pluginConfigAPI interface {
	GetDependentPlugins() []*jenkins.JenkinsPlugin
	PreConfig(jenkins.JenkinsAPI) error
}

func installPlugins(jenkinsClient jenkins.JenkinsAPI, pluginConfigs []pluginConfigAPI, enableRestart bool) error {
	var plugins []*jenkins.JenkinsPlugin
	for _, pluginConfig := range pluginConfigs {
		plugins = append(plugins, pluginConfig.GetDependentPlugins()...)
	}
	return jenkinsClient.InstallPluginsIfNotExists(plugins, enableRestart)
}

func preConfigPlugins(jenkinsClient jenkins.JenkinsAPI, pluginConfigs []pluginConfigAPI) error {
	for _, pluginConfig := range pluginConfigs {
		err := pluginConfig.PreConfig(jenkinsClient)
		if err != nil {
			log.Debugf("jenkins plugin %+v preConfig error", pluginConfig)
			return err
		}
	}
	return nil
}
