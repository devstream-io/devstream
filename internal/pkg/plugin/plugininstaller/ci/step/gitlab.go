package step

import (
	"os"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

const (
	sshKeyCredentialName = "gitlabSSHKeyCredential"
	gitlabCredentialName = "gitlabCredential"
	gitlabConnectionName = "gitlabConnection"
)

type GitlabStepConfig struct {
	SSHPrivateKey string `mapstructure:"sshPrivateKey"`
	RepoOwner     string `mapstructure:"repoOwner"`
	BaseURL       string `mapstructure:"baseURL"`
}

func (g *GitlabStepConfig) GetJenkinsPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{
		{
			Name:    "gitlab-plugin",
			Version: "1.5.35",
		},
	}
}

func (g *GitlabStepConfig) ConfigJenkins(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
	// 1. create ssh credentials
	if g.SSHPrivateKey == "" {
		log.Warnf("jenkins gitlab ssh key not config, private repo can't be clone")
	} else {
		err := jenkinsClient.CreateSSHKeyCredential(
			sshKeyCredentialName, g.RepoOwner, g.SSHPrivateKey,
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
		GitLabConnectionName: gitlabConnectionName,
		GitlabURL:            g.BaseURL,
	}, nil
}

func (g *GitlabStepConfig) ConfigSCM(client scm.ClientOperation) error {
	return nil
}

func newGitlabStep(config *StepGlobalOption) *GitlabStepConfig {
	return &GitlabStepConfig{
		SSHPrivateKey: config.RepoInfo.SSHPrivateKey,
		RepoOwner:     config.RepoInfo.GetRepoOwner(),
		BaseURL:       config.RepoInfo.BaseURL,
	}
}
