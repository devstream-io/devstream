package base

type ImageRepoStepConfig struct {
	URL  string `mapstructure:"url" validate:"url"`
	User string `mapstructure:"user"`
}
