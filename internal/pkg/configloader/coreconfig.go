package configloader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// CoreConfig is the struct representing the complete original configuration YAML files.
type CoreConfig struct {
	// TODO(daniel-hutao): Relative path support
	VarFile string `yaml:"varFile"`
	// TODO(daniel-hutao): Relative path support
	ToolFile string `yaml:"toolFile"`
	State    *State `yaml:"state"`
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
}

func (c *CoreConfig) Validate() (bool, error) {
	if c.State == nil {
		return false, fmt.Errorf("state config is empty")
	}

	errors := make([]error, 0)

	log.Infof("Got Backend from config: %s", c.State.Backend)
	switch c.State.Backend {
	case "local":
		if c.State.Options.StateFile == "" {
			log.Debugf("The stateFile has not been set, default value %s will be used.", local.DefaultStateFile)
			c.State.Options.StateFile = local.DefaultStateFile
		}
	case "s3":
		if c.State.Options.Bucket == "" {
			errors = append(errors, fmt.Errorf("state s3 Bucket is empty"))
		}
		if c.State.Options.Region == "" {
			errors = append(errors, fmt.Errorf("state s3 Region is empty"))
		}
		if c.State.Options.Key == "" {
			errors = append(errors, fmt.Errorf("state s3 Key is empty"))
		}
	default:
		errors = append(errors, fmt.Errorf("backend type error"))
	}

	if len(errors) != 0 {
		var retErr []string
		log.Error("Config file validation failed.")
		for i, err := range errors {
			log.Errorf("%d -> %s.", i+1, err)
			retErr = append(retErr, err.Error())
		}
		return false, fmt.Errorf("%s", strings.Join(retErr, "; "))
	}

	return true, nil
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
		log.Errorf("Failed to get absolute path fo \"%s\".", filePath)
		return "", err
	}
	log.Debugf("Abs path is %s.", absFilePath)
	if fileExist(absFilePath) {
		return absFilePath, nil
	} else {
		return "", fmt.Errorf("file %s not exists", absFilePath)
	}
}
