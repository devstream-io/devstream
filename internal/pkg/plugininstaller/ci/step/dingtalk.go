package step

import (
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/jenkins/dingtalk"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm/github"
)

const (
	dingTalkSecretKey = "DINGTALK_SECURITY_VALUE"
)

type DingtalkStepConfig struct {
	Name          string `mapstructure:"name"`
	Webhook       string `mapstructure:"webhook"`
	SecurityValue string `mapstructure:"securityValue" validate:"required"`
	SecurityType  string `mapstructure:"securityType" validate:"required,oneof=KEY SECRET"`
	AtUsers       string `mapstructure:"atUsers"`
}

// GetJenkinsPlugins return jenkins plugins info
func (g *DingtalkStepConfig) GetJenkinsPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{
		{
			Name:    "dingding-notifications",
			Version: "2.4.10",
		},
	}
}

// JenkinsConfig config jenkins and return casc config
func (g *DingtalkStepConfig) ConfigJenkins(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
	config := dingtalk.BotConfig{
		RobotConfigs: []dingtalk.BotInfoConfig{
			{
				// use dingtalk robot name as id
				ID:            g.Name,
				Name:          g.Name,
				Webhook:       g.Webhook,
				SecurityType:  g.SecurityType,
				SecurityValue: g.SecurityValue,
			},
		},
	}
	log.Info("jenkins plugin dingtalk start config...")
	return nil, jenkinsClient.ApplyDingTalkBot(config)
}

func (g *DingtalkStepConfig) ConfigGithub(client *github.Client) error {
	return client.AddRepoSecret(dingTalkSecretKey, g.SecurityValue)
}
