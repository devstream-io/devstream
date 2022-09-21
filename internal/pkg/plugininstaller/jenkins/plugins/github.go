package plugins

import (
	"os"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
)

const (
	GithubCredentialName     = "jenkinsGithubCredential"
	githubCredentialUserName = "jenkinsGithubDevstreamUser"
)

type GithubJenkinsConfig struct {
	JenkinsURL string
}

func (g *GithubJenkinsConfig) GetDependentPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{
		{
			Name:    "github-branch-source",
			Version: "1695.v88de84e9f6b_9",
		},
	}
}

func (g *GithubJenkinsConfig) PreConfig(jenkinsClient jenkins.JenkinsAPI) error {
	// 1. create github credentials by github token
	err := jenkinsClient.CreatePasswordCredential(
		GithubCredentialName,
		githubCredentialUserName,
		os.Getenv(github.TokenEnvKey),
	)
	if err != nil {
		log.Debugf("jenkins preinstall github credentials failed: %s", err)
		return err
	}
	// 2. config github plugin casc
	return jenkinsClient.ConfigCascForRepo(&jenkins.RepoCascConfig{
		RepoType:     "github",
		CredentialID: GithubCredentialName,
		JenkinsURL:   g.JenkinsURL,
	})
}
