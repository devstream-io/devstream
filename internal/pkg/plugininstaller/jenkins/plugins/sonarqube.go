package plugins

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/base"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	sonarQubeTokenCredentialName = "sonarqubeTokenCredential"
)

type SonarQubeJenkinsConfig struct {
	base.SonarQubeStepConfig `mapstructure:",squash"`
}

func (g *SonarQubeJenkinsConfig) getDependentPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{
		{
			Name:    "sonar",
			Version: "2.14",
		},
	}
}

func (g *SonarQubeJenkinsConfig) config(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
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
func (g *SonarQubeJenkinsConfig) setRenderVars(vars *jenkins.JenkinsFileRenderInfo) {
	vars.SonarqubeEnable = true
}
