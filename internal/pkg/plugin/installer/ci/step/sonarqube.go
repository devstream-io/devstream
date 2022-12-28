package step

import (
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

const (
	sonarSecretKey = "SONAR_SECRET_TOKEN"
)

type SonarQubeStepConfig struct {
	Name  string `mapstructure:"name" validate:"required"`
	Token string `mapstructure:"token"`
	URL   string `mapstructure:"url" validate:"url"`
}

func (g *SonarQubeStepConfig) GetJenkinsPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{
		{
			Name:    "sonar",
			Version: "2.14",
		},
	}
}

func (g *SonarQubeStepConfig) ConfigJenkins(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
	log.Info("jenkins plugin sonarqube start config...")
	// 1. install token credential
	err := jenkinsClient.CreateSecretCredential(
		g.Name,
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
		SonarTokenCredentialID: sonarSecretKey,
	}, nil
}

func (s *SonarQubeStepConfig) ConfigSCM(client scm.ClientOperation) error {
	return client.AddRepoSecret(sonarSecretKey, s.Token)
}
