package plugins

import (
	"os"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
)

const (
	githubCredentialName = "jenkinsGithubCredential"
)

type GithubJenkinsConfig struct {
	RepoOwner string `mapstructure:"repoOwner"`
}

func (g *GithubJenkinsConfig) getDependentPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{
		{
			Name:    "github-branch-source",
			Version: "1695.v88de84e9f6b_9",
		},
	}
}

func (g *GithubJenkinsConfig) config(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
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

func NewGithubPlugin(config *PluginGlobalConfig) *GithubJenkinsConfig {
	return &GithubJenkinsConfig{
		RepoOwner: config.RepoInfo.GetRepoOwner(),
	}
}

func (g *GithubJenkinsConfig) setRenderVars(vars *jenkins.JenkinsFileRenderInfo) {
}
