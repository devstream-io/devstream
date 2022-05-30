package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var pluginDefaultConfigs = map[string]string{
	"argocd":                         ArgocdDefaultConfig,
	"argocdapp":                      ArgocdappDefaultConfig,
	"devlake":                        DevlakeDefaultConfig,
	"github-repo-scaffolding-golang": GithubRepoScaffoldingGolangDefaultConfig,
	"githubactions-golang":           GithubActionsGolangDefaultConfig,
	"githubactions-nodejs":           GithubActionsNodejsDefaultConfig,
	"githubactions-python":           GithubActionsPythonDefaultConfig,
	"gitlabci-generic":               GitlabCIGenericDefaultConfig,
	"gitlabci-golang":                GitlabCIGolangDefaultConfig,
	"jenkins":                        JenkinsDefaultConfig,
	"jira-github-integ":              JiraGithubDefaultConfig,
	"kube-prometheus":                KubePrometheusDefaultConfig,
	"openldap":                       OpenldapDefaultConfig,
	"trello-github-integ":            TrelloGithubDefaultConfig,
	"trello":                         TrelloDefaultConfig,
	"helm-generic":                   HelmGenericDefaultConfig,
	"gitlab-repo-scaffolding-golang": GitLabRepoScaffoldingGolangDefaultConfig,
	"hashicorp-vault":                VaultDefaultConfig,
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
