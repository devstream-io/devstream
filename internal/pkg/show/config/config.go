package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/show/config/plugin"
)

var pluginDefaultConfigs = map[string]string{
	"argocd":                         plugin.ArgocdDefaultConfig,
	"argocdapp":                      plugin.ArgocdappDefaultConfig,
	"devlake":                        plugin.DevlakeDefaultConfig,
	"github-repo-scaffolding-golang": plugin.GithubRepoScaffoldingGolangDefaultConfig,
	"githubactions-golang":           plugin.GithubActionsGolangDefaultConfig,
	"githubactions-nodejs":           plugin.GithubActionsNodejsDefaultConfig,
	"githubactions-python":           plugin.GithubActionsPythonDefaultConfig,
	"gitlabci-generic":               plugin.GitlabCIGenericDefaultConfig,
	"gitlabci-golang":                plugin.GitlabCIGolangDefaultConfig,
	"jenkins":                        plugin.JenkinsDefaultConfig,
	"jira-github-integ":              plugin.JiraGithubDefaultConfig,
	"kube-prometheus":                plugin.KubePrometheusDefaultConfig,
	"openldap":                       plugin.OpenldapDefaultConfig,
	"trello-github-integ":            plugin.TrelloGithubDefaultConfig,
	"trello":                         plugin.TrelloDefaultConfig,
	"helm-generic":                   plugin.HelmGenericDefaultConfig,
	"gitlab-repo-scaffolding-golang": plugin.GitLabRepoScaffoldingGolangDefaultConfig,
	"hashicorp-vault":                plugin.VaultDefaultConfig,
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
