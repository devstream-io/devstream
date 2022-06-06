package configloader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/pkg/util/log"
)

// CoreConfig is the config YAML user writes
type CoreConfig struct {
	// TODO(daniel-hutao): Relative path support
	VarFile string `yaml:"varFile"`
	// TODO(daniel-hutao): Relative path support
	ToolFile string `yaml:"toolFile"`
	State    *State `yaml:"state"`
}

type State struct {
	Backend string             `yaml:"backend"`
	Options StateConfigOptions `yaml:"options"`
}

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

	c.VarFile, err = c.genAbsFilePath(c.VarFile)
	if err != nil {
		return err
	}

	return nil
}

func (c *CoreConfig) ParseToolFilePath() error {
	var err error

	c.ToolFile, err = c.genAbsFilePath(c.ToolFile)
	if err != nil {
		return err
	}

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

	if filepath.IsAbs(filePath) {
		log.Debugf("Abs path is %s.", filePath)
		if fileExist(filePath) {
			return filePath, nil
		} else {
			return "", fmt.Errorf("file %s not exists", filePath)
		}
		return filePath, nil
	}

	absFilePath := filepath.Join(filepath.Dir(filePath), filePath)
	log.Debugf("Abs path is %s.", absFilePath)
	if fileExist(absFilePath) {
		return absFilePath, nil
	} else {
		return "", fmt.Errorf("file %s not exists", absFilePath)
	}
	return absFilePath, nil
}
