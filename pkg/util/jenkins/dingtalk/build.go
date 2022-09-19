package dingtalk

var DefaultNoticeOccasions = []string{
	"START",
	"ABORTED",
	"FAILURE",
	"SUCCESS",
	"UNSTABLE",
	"NOT_BUILT",
}

type (
	BotConfig struct {
		NoticeOccasions []string
		Verbose         bool
		ProxyConfig     ProxyConfig
		RobotConfigs    []BotInfoConfig
	}

	BotInfoConfig struct {
		ID            string
		Name          string
		Webhook       string
		SecurityType  string
		SecurityValue string
	}
)

func BuildDingTalkConfig(config BotConfig) (InnerBotConfig, error) {
	// set default notice occasions and proxy config
	if len(config.NoticeOccasions) == 0 {
		config.NoticeOccasions = DefaultNoticeOccasions
	}
	if config.ProxyConfig.Type == "" {
		config.ProxyConfig = ProxyConfig{
			Type: "DIRECT",
			Host: "",
			Port: "0",
		}
	}

	// build global config
	innerConfig := InnerBotConfig{
		NoticeOccasions: config.NoticeOccasions,
		Verbose:         config.Verbose,
		ProxyConfig:     config.ProxyConfig,
		CoreApply:       "true",
	}

	// build robot detail config
	for _, bot := range config.RobotConfigs {
		// set security type and value
		securityPolicyConfigs := []SecurityPolicyConfigs{
			{
				Value: "",
				Type:  SecurityTypeKey,
				Desc:  SecurityTypeKeyChinese,
			},
			{
				Value: "",
				Type:  SecurityTypeSecret,
				Desc:  SecurityTypeSecretChinese,
			},
		}
		switch bot.SecurityType {
		case SecurityTypeKey, SecurityTypeKeyChinese:
			securityPolicyConfigs[0].Value = bot.SecurityValue
		case SecurityTypeSecret, SecurityTypeSecretChinese:
			securityPolicyConfigs[1].Value = bot.SecurityValue
		}

		innerConfig.RobotConfigs = append(innerConfig.RobotConfigs, RobotConfigs{
			ID:                    bot.ID,
			Name:                  bot.Name,
			Webhook:               bot.Webhook,
			SecurityPolicyConfigs: securityPolicyConfigs,
		})
	}

	return innerConfig, nil
}
