package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/show/config/plugin"
)

var pluginDefaultConfigs = map[string]string{
	"argocd":                         plugin.ArgocdDefaultConfig,
	"gitlabci-generic":               plugin.GitlabCIGenericDefaultConfig,
	"argocdapp":                      plugin.ArgocdappDefaultConfig,
	"github-repo-scaffolding-golang": plugin.GithubRepoScaffoldingGolangDefaultConfig,
	"githubactions-golang":           plugin.GithubActionsGolangDefaultConfig,
	"devlake":                        plugin.DevlakeDefaultConfig,
	"githubactions-python":           plugin.GithubActionsPythonDefaultConfig,
	"githubactions-nodejs":           plugin.GithubActionsNodejsDefaultConfig,
	"gitlabci-golang":                plugin.GitlabCIGolangDefaultConfig,
}

func Show() error {
	plugin := viper.GetString("plugin")
	if plugin == "" {
		return fmt.Errorf("empty plugin name. Maybe you forgot to add --plugin=PLUGIN_NAME?")
	}
	if config, ok := pluginDefaultConfigs[plugin]; ok {
		fmt.Println(config)
		return nil
	}
	return fmt.Errorf("illegal plugin name : < %s >", plugin)
}
