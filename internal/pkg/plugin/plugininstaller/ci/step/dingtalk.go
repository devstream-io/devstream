package step

import (
	"fmt"
	"strings"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
	"github.com/devstream-io/devstream/pkg/util/jenkins/dingtalk"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

const (
	dingTalkSecretToken = "DINGTALK_SECURITY_TOKEN"
	dingTalkSecretVal   = "DINGTALK_SECURITY_VALUE"
)

type DingtalkStepConfig struct {
	Name          string `mapstructure:"name"`
	Webhook       string `mapstructure:"webhook"`
	SecurityValue string `mapstructure:"securityValue"`
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
	robotConfig := dingtalk.BotInfoConfig{
		// use dingtalk robot name as id
		ID:           g.Name,
		Name:         g.Name,
		Webhook:      g.Webhook,
		SecurityType: g.SecurityType,
	}
	if g.needSetSecurityValue() {
		robotConfig.SecurityValue = g.SecurityValue
	}
	config := dingtalk.BotConfig{
		RobotConfigs: []dingtalk.BotInfoConfig{
			robotConfig,
		},
	}
	log.Info("jenkins plugin dingtalk start config...")
	return nil, jenkinsClient.ApplyDingTalkBot(config)
}

// ConfigSCM is used to config github and gitlab variables
func (g *DingtalkStepConfig) ConfigSCM(client scm.ClientOperation) error {
	splitWebhook := strings.Split(g.Webhook, "=")
	if len(splitWebhook) < 2 {
		return fmt.Errorf("step scm dingTalk.webhook is not valid")
	}
	token := splitWebhook[len(splitWebhook)-1]
	if g.needSetSecurityValue() {
		err := client.AddRepoSecret(dingTalkSecretVal, g.SecurityValue)
		if err != nil {
			return err
		}
	}
	return client.AddRepoSecret(dingTalkSecretToken, token)
}

// if dingTalk use SECRET securityType, we need to set SecurityValue for use later
func (g *DingtalkStepConfig) needSetSecurityValue() bool {
	if g.SecurityType == "SECRET" && g.SecurityValue != "" {
		return true
	}
	return false
}
