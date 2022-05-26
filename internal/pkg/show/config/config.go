package config

import (
	"fmt"

	"github.com/devstream-io/devstream/examples"

	"github.com/spf13/viper"
)

var pluginDefaultConfigs = map[string]string{
	"argocd":                         examples.ArgocdDefaultConfig,
	"argocdapp":                      examples.ArgocdappDefaultConfig,
	"devlake":                        examples.DevlakeDefaultConfig,
	"github-repo-scaffolding-golang": examples.GithubRepoScaffoldingGolangDefaultConfig,
	"githubactions-golang":           examples.GithubActionsGolangDefaultConfig,
	"githubactions-nodejs":           examples.GithubActionsNodejsDefaultConfig,
	"githubactions-python":           examples.GithubActionsPythonDefaultConfig,
	"gitlabci-generic":               examples.GitlabCIGenericDefaultConfig,
	"gitlabci-golang":                examples.GitlabCIGolangDefaultConfig,
	"jenkins":                        examples.JenkinsDefaultConfig,
	"jira-github-integ":              examples.JiraGithubDefaultConfig,
	"kube-prometheus":                examples.KubePrometheusDefaultConfig,
	"openldap":                       examples.OpenldapDefaultConfig,
	"trello-github-integ":            examples.TrelloGithubDefaultConfig,
	"trello":                         examples.TrelloDefaultConfig,
	"helm-generic":                   examples.HelmGenericDefaultConfig,
	"gitlab-repo-scaffolding-golang": examples.GitLabRepoScaffoldingGolangDefaultConfig,
	"hashicorp-vault":                examples.VaultDefaultConfig,
}

func Show() error {
	plugin := viper.GetString("plugin")
	if plugin == "" {
		fmt.Println(DefaultConfig)
		return nil
	}
	if config, ok := pluginDefaultConfigs[plugin]; ok {
		fmt.Println(config)
		return nil
	}
	return fmt.Errorf("illegal plugin name : < %s >", plugin)
}
