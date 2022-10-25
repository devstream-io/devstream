package base

type DingtalkStepConfig struct {
	Name          string `mapstructure:"name"`
	Webhook       string `mapstructure:"webhook"`
	SecurityValue string `mapstructure:"securityValue" validate:"required"`
	SecurityType  string `mapstructure:"securityType" validate:"required,oneof=KEY SECRET"`
	AtUsers       string `mapstructure:"atUsers"`
}
