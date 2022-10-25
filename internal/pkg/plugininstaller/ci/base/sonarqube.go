package base

type SonarQubeStepConfig struct {
	Name  string `mapstructure:"name"`
	Token string `mapstructure:"token"`
	URL   string `mapstructure:"url" validate:"url"`
}
