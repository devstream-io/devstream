package step

import (
	"os"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
)

const (
	githubCredentialName = "githubCredential"
)

type GithubStepConfig struct {
	RepoOwner string `mapstructure:"repoOwner"`
}

func (g *GithubStepConfig) GetJenkinsPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{
		{
			Name:    "github-branch-source",
			Version: "1695.v88de84e9f6b_9",
		},
	}
}

func (g *GithubStepConfig) ConfigJenkins(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
	// 1. create github credentials by github token
	err := jenkinsClient.CreatePasswordCredential(
		githubCredentialName,
		g.RepoOwner,
		os.Getenv(github.TokenEnvKey),
	)
	if err != nil {
		log.Debugf("jenkins preinstall github credentials failed: %s", err)
		return nil, err
	}
	// 2. config github plugin casc
	return &jenkins.RepoCascConfig{
		RepoType:     "github",
		CredentialID: githubCredentialName,
		JenkinsURL:   jenkinsClient.GetBasicInfo().URL,
	}, nil
}

func (g *GithubStepConfig) ConfigSCM(client scm.ClientOperation) error {
	return nil
}

func newGithubStep(config *StepGlobalOption) *GithubStepConfig {
	return &GithubStepConfig{
		RepoOwner: config.RepoInfo.GetRepoOwner(),
	}
}
