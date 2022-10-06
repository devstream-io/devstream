package plugins

import (
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	sonarQubeTokenCredentialName = "sonarqubeTokenCredential"
)

type SonarQubeJenkinsConfig struct {
	Name  string `mapstructure:"name"`
	Token string `mapstructure:"token"`
	URL   string `mapstructure:"url" validate:"url"`
}

func (g *SonarQubeJenkinsConfig) GetDependentPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{
		{
			Name:    "sonar",
			Version: "2.14",
		},
	}
}

func (g *SonarQubeJenkinsConfig) PreConfig(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
	log.Info("jenkins plugin sonarqube start config...")
	// 1. install token credential
	err := jenkinsClient.CreateSecretCredential(
		sonarQubeTokenCredentialName,
		g.Token,
	)
	if err != nil {
		log.Debugf("jenkins preinstall github credentials failed: %s", err)
		return nil, err
	}
	// 2. config sonarqube casc
	return &jenkins.RepoCascConfig{
		SonarqubeURL:           g.URL,
		SonarqubeName:          g.Name,
		SonarTokenCredentialID: sonarQubeTokenCredentialName,
	}, nil

}

func (g *SonarQubeJenkinsConfig) UpdateJenkinsFileRenderVars(vars *jenkins.JenkinsFileRenderInfo) {
	vars.SonarqubeEnable = true
}
