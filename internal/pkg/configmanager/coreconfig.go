package configmanager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// CoreConfig is the struct representing the complete original configuration YAML files.
type CoreConfig struct {
	// TODO(daniel-hutao): Relative path support
	VarFile string `yaml:"varFile"`
	// TODO(daniel-hutao): Relative path support
	ToolFile string `yaml:"toolFile"`
	// abs path of the plugin dir
	PluginDir string `yaml:"pluginDir"`
	State     *State `yaml:"state"`
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

func (c *CoreConfig) Validate() error {
	if c.State == nil {
		return fmt.Errorf("state config is empty")
	}

	return nil
}

func (c *CoreConfig) ParseVarFilePath() error {
	var err error

	log.Debugf("Original varFile path: '%s'.", c.VarFile)
	if c.VarFile == "" {
		return nil
	}
	c.VarFile, err = c.genAbsFilePath(c.VarFile)
	if err != nil {
		return err
	}
	log.Debugf("Absolute varFile path: '%s'.", c.VarFile)

	return nil
}

func (c *CoreConfig) ParseToolFilePath() error {
	var err error

	log.Debugf("Original toolFile path: '%s'.", c.ToolFile)
	if c.VarFile == "" {
		return nil
	}
	c.ToolFile, err = c.genAbsFilePath(c.ToolFile)
	if err != nil {
		return err
	}
	log.Debugf("Absolute toolFile path: '%s'.", c.ToolFile)

	return nil
}

// genAbsFilePath return all of the path with a given file name
func (c *CoreConfig) genAbsFilePath(filePath string) (string, error) {
	fileExist := func(path string) bool {
		if _, err := os.Stat(filePath); err != nil {
			log.Errorf("File %s not exists. Error: %s", filePath, err)
			return false
		}
		return true
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		log.Errorf(`Failed to get absolute path fo "%s".`, filePath)
		return "", err
	}
	log.Debugf("Abs path is %s.", absFilePath)
	if fileExist(absFilePath) {
		return absFilePath, nil
	} else {
		return "", fmt.Errorf("file %s not exists", absFilePath)
	}
}
