package config

import _ "embed"

var DefaultConfig = `# config.yaml sample:
# var file path, you can set it to absolute path or relative path.
varFile: variables.yaml # here is a relative path. (defaults is ./variables.yaml)
# tool file path, you can set it to absolute path or relative path.
toolFile: tools.yaml # here is a relative path.
# state config
state:
  backend: local # backend can be local or s3
  options:
    stateFile: devstream.state

# tools.yaml sample:
tools:
- name: github-repo-scaffolding-golang
  instanceID: default
  options:
    owner: [[ githubUsername ]]
    org: ""
    repo: [[ repoName ]]

# variables.yaml sample:
githubUsername: daniel-hutao
repo: go-webapp-demo`

// plugin default config
var (
	//go:embed plugins/argocd.yaml
	ArgocdDefaultConfig string

	//go:embed plugins/argocdapp.yaml
	ArgocdappDefaultConfig string

	//go:embed plugins/devlake.yaml
	DevlakeDefaultConfig string

	//go:embed plugins/github-repo-scaffolding-golang.yaml
	GithubRepoScaffoldingGolangDefaultConfig string

	//go:embed plugins/githubactions-golang.yaml
	GithubActionsGolangDefaultConfig string

	//go:embed plugins/githubactions-nodejs.yaml
	GithubActionsNodejsDefaultConfig string

	//go:embed plugins/githubactions-python.yaml
	GithubActionsPythonDefaultConfig string

	//go:embed plugins/gitlabci-generic.yaml
	GitlabCIGenericDefaultConfig string

	//go:embed plugins/gitlabci-golang.yaml
	GitlabCIGolangDefaultConfig string

	//go:embed plugins/jenkins.yaml
	JenkinsDefaultConfig string

	//go:embed plugins/jira-github-integ.yaml
	JiraGithubDefaultConfig string

	//go:embed plugins/kube-prometheus.yaml
	KubePrometheusDefaultConfig string

	//go:embed plugins/openldap.yaml
	OpenldapDefaultConfig string

	//go:embed plugins/trello-github-integ.yaml
	TrelloGithubDefaultConfig string

	//go:embed plugins/trello.yaml
	TrelloDefaultConfig string

	//go:embed plugins/helm-generic.yaml
	HelmGenericDefaultConfig string

	//go:embed plugins/gitlab-repo-scaffolding-golang.yaml
	GitLabRepoScaffoldingGolangDefaultConfig string

	//go:embed plugins/hashicorp-vault.yaml
	VaultDefaultConfig string
)
