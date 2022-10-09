package plugins

import (
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/jenkins/dingtalk"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type DingtalkJenkinsConfig struct {
	Name          string `mapstructure:"name"`
	Webhook       string `mapstructure:"webhook"`
	SecurityValue string `mapstructure:"securityValue" validate:"required"`
	SecurityType  string `mapstructure:"securityType" validate:"required,oneof=KEY SECRET"`
	AtUsers       string `mapstructure:"atUsers"`
}

func (g *DingtalkJenkinsConfig) GetDependentPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{
		{
			Name:    "dingding-notifications",
			Version: "2.4.10",
		},
	}
}

func (g *DingtalkJenkinsConfig) PreConfig(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
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

func (g *DingtalkJenkinsConfig) UpdateJenkinsFileRenderVars(vars *jenkins.JenkinsFileRenderInfo) {
	vars.DingtalkRobotID = g.Name
	vars.DingtalkAtUser = g.AtUsers
}
