package plugins

import (
	"os"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	SSHKeyCredentialName = "gitlabSSHKeyCredential"
	gitlabCredentialName = "gitlabCredential"
	GitlabConnectionName = "gitlabConnection"
)

type GitlabJenkinsConfig struct {
	SSHPrivateKey string
	RepoOwner     string
	BaseURL       string
}

func (g *GitlabJenkinsConfig) GetDependentPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{
		{
			Name:    "gitlab-plugin",
			Version: "1.5.35",
		},
	}
}

func (g *GitlabJenkinsConfig) PreConfig(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
	// 1. create ssh credentials
	if g.SSHPrivateKey == "" {
		log.Warnf("jenkins gitlab ssh key not config, private repo can't be clone")
	} else {
		err := jenkinsClient.CreateSSHKeyCredential(
			SSHKeyCredentialName, g.RepoOwner, g.SSHPrivateKey,
		)
		if err != nil {
			return nil, err
		}
	}
	// 2. create gitlab connection by casc
	err := jenkinsClient.CreateGiltabCredential(gitlabCredentialName, os.Getenv("GITLAB_TOKEN"))
	if err != nil {
		log.Debugf("jenkins preinstall gitlab credentials failed: %s", err)
		return nil, err
	}
	// 3. config gitlab casc
	return &jenkins.RepoCascConfig{
		RepoType:             "gitlab",
		CredentialID:         gitlabCredentialName,
		GitLabConnectionName: GitlabConnectionName,
		GitlabURL:            g.BaseURL,
	}, nil
}

func (g *GitlabJenkinsConfig) UpdateJenkinsFileRenderVars(vars *jenkins.JenkinsFileRenderInfo) {
}
