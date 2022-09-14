package jenkins

import (
	_ "embed"

	"github.com/devstream-io/devstream/pkg/util/jenkins/dingtalk"
)

func ExampleApplyDingTalkBot() {
	basicAuth := &BasicAuth{
		Username: "admin",
		Password: "changeme",
	}
	j, err := NewClient("http://localhost:32000", basicAuth)
	if err != nil {
		panic(err)
	}
	config := dingtalk.BotConfig{
		RobotConfigs: []dingtalk.BotInfoConfig{
			{
				// It is recommended to use a fixed id, otherwise the bot cannot be referenced in the pipeline
				ID:            "fd8f5e95-3b33-430e-84eb-d9e0f3027024",
				Name:          "通知测试机器人2",
				Webhook:       "https://oapi.dingtalk.com/robot/send?access_token=12d859fxxxxxxxxxxxxx5559b57b",
				SecurityType:  "SECRET",
				SecurityValue: "SEC64a21axxxxxxxxxxxxxxxxxxxxxdf435e72dec",
			},
		},
	}
	err = j.ApplyDingTalkBot(config)
	if err != nil {
		panic(err)
	}
}
