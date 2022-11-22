package configmanager

import (
	"fmt"
)

// Config records rendered config values and is used as a general config in DevStream.
// Also see the rawConfig struct in rawconfig.go, it represents the original config.
type Config struct {
	// Command line flag have a higher priority than the config file.
	// If you used the `--plugin-dir` flag with `dtm`, then the "pluginDir" in the config file will be ignored.
	PluginDir string `yaml:"pluginDir"`
	Tools     Tools  `yaml:"tools"`
	State     *State `yaml:"state"`
}

func (c *Config) validate() error {
	if c.State == nil {
		return fmt.Errorf("state is not defined")
	}
	return nil
}

// State is the struct for reading the state configuration in the config file.
// It defines how the state is stored, specifies the type of backend and related options.
type State struct {
	Backend string             `yaml:"backend"`
	Options StateConfigOptions `yaml:"options"`
}

// StateConfigOptions is the struct for reading the options of the state backend.
type StateConfigOptions struct {
	// for s3 backend
	Bucket string `yaml:"bucket"`
	Region string `yaml:"region"`
	Key    string `yaml:"key"`
	// for local backend
	StateFile string `yaml:"stateFile"`
	// for k8s backend
	Namespace string `yaml:"namespace"`
	ConfigMap string `yaml:"configmap"`
}
