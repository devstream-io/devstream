package plugins

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/ci/base"
	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/jenkins/dingtalk"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type DingtalkJenkinsConfig struct {
	base.DingtalkStepConfig `mapstructure:",squash"`
}

func (g *DingtalkJenkinsConfig) getDependentPlugins() []*jenkins.JenkinsPlugin {
	return []*jenkins.JenkinsPlugin{
		{
			Name:    "dingding-notifications",
			Version: "2.4.10",
		},
	}
}

func (g *DingtalkJenkinsConfig) config(jenkinsClient jenkins.JenkinsAPI) (*jenkins.RepoCascConfig, error) {
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

func (g *DingtalkJenkinsConfig) setRenderVars(vars *jenkins.JenkinsFileRenderInfo) {
	vars.DingtalkRobotID = g.Name
	vars.DingtalkAtUser = g.AtUsers
}
