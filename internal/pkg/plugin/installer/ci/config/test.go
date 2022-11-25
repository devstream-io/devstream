package config

type TestOption struct {
	Enable                *bool    `mapstructure:"enable"`
	Command               []string `mapstructure:"command"`
	ContainerName         string   `mapstructure:"containerName"`
	CoverageCommand       string   `mapstructure:"coverageCommand"`
	CoverageStatusCommand string   `mapstructure:"CoverageStatusCommand"`
}
