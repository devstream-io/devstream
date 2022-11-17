package configmanager

import (
	"fmt"
	"runtime"
)

var (
	GOOS   string = runtime.GOOS
	GOARCH string = runtime.GOARCH
)

// Config records rendered config values and is used as a general config in DevStream.
type Config struct {
	// Command line flag have a higher priority than the config file.
	// If you used the `--plugin-dir` flag with `dtm`, then the "pluginDir" in the config file will be ignored.
	PluginDir string
	Tools     Tools `yaml:"tools"`
	Apps      Apps  `yaml:"apps"`
	State     *State
}

// ConfigRaw is used to describe original raw configs read from files
type ConfigRaw struct {
	VarFile           string             `yaml:"varFile"`
	ToolFile          string             `yaml:"toolFile"`
	AppFile           string             `yaml:"appFile"`
	TemplateFile      string             `yaml:"templateFile"`
	PluginDir         string             `yaml:"pluginDir"`
	State             *State             `yaml:"state"`
	Tools             []Tool             `yaml:"tools"`
	AppsInConfig      []AppInConfig      `yaml:"apps"`
	PipelineTemplates []PipelineTemplate `yaml:"pipelineTemplates"`
	GlobalVars        map[string]any     `yaml:"-"`
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
	// for ConfigMap backend
	Namespace string `yaml:"namespace"`
	ConfigMap string `yaml:"configmap"`
}

func (c *Config) Validate() (errs []error) {
	errs = append(errs, c.Tools.validate()...)
	errs = append(errs, c.Tools.validateDependency()...)
	errs = append(errs, c.Apps.validate()...)

	if c.State == nil {
		errs = append(errs, fmt.Errorf("state is not defined"))
	}

	return errs
}
