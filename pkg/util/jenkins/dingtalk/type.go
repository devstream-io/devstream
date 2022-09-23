package dingtalk

type (
	InnerBotConfig struct {
		NoticeOccasions []string       `json:"noticeOccasions"`
		Verbose         bool           `json:"verbose"`
		ProxyConfig     ProxyConfig    `json:"proxyConfig"`
		RobotConfigs    []RobotConfigs `json:"robotConfigs"`
		CoreApply       string         `json:"core:apply"`
	}

	ProxyConfig struct {
		Type string `json:"type"`
		Host string `json:"host"`
		Port string `json:"port"`
	}
	SecurityPolicyConfigs struct {
		Value string `json:"value"`
		Type  string `json:"type"`
		Desc  string `json:"desc"`
	}
	RobotConfigs struct {
		ID                    string                  `json:"id"`
		Name                  string                  `json:"name"`
		Webhook               string                  `json:"webhook"`
		SecurityPolicyConfigs []SecurityPolicyConfigs `json:"securityPolicyConfigs"`
	}
)

const (
	SecurityTypeKey           = "KEY"
	SecurityTypeKeyChinese    = "关键字"
	SecurityTypeSecret        = "SECRET"
	SecurityTypeSecretChinese = "加密"
)
