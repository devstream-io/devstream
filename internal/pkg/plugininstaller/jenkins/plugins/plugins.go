package plugins

import (
	"github.com/imdario/mergo"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

type PluginConfigAPI interface {
	getDependentPlugins() []*jenkins.JenkinsPlugin
	config(jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error)
	setRenderVars(*jenkins.JenkinsFileRenderInfo)
}

type PluginGlobalConfig struct {
	RepoInfo *git.RepoInfo
}

func GetRepoCredentialsId(repoInfo *git.RepoInfo) string {
	if repoInfo.IsGitlab() && repoInfo.SSHPrivateKey != "" {
		return sshKeyCredentialName
	} else if repoInfo.IsGithub() {
		return githubCredentialName
	}
	return ""
}

func EnsurePluginInstalled(jenkinsClient jenkins.JenkinsAPI, pluginConfigs []PluginConfigAPI) error {
	var plugins []*jenkins.JenkinsPlugin
	for _, pluginConfig := range pluginConfigs {
		plugins = append(plugins, pluginConfig.getDependentPlugins()...)
	}
	return jenkinsClient.InstallPluginsIfNotExists(plugins)
}

func ConfigPlugins(jenkinsClient jenkins.JenkinsAPI, pluginConfigs []PluginConfigAPI) error {
	globalCascConfig := new(jenkins.RepoCascConfig)
	for _, pluginConfig := range pluginConfigs {
		cascConfig, err := pluginConfig.config(jenkinsClient)
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

func GetPluginsRenderVariables(pluginConfigs []PluginConfigAPI) *jenkins.JenkinsFileRenderInfo {
	jenkinsFileConfig := &jenkins.JenkinsFileRenderInfo{}
	for _, p := range pluginConfigs {
		p.setRenderVars(jenkinsFileConfig)
	}
	return jenkinsFileConfig
}
