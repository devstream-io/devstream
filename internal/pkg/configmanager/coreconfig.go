package configmanager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	defaultNamespace     = "devstream"
	defaultConfigMapName = "devstream-state"
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

func (c *CoreConfig) ValidateAndDefault() error {
	if c.State == nil {
		return fmt.Errorf("state config is empty")
	}

	log.Infof("Got Backend from config: %s", c.State.Backend)

	errs := ValidateAndDefaultBackend(c.State)
	if len(errs) != 0 {
		var retErr []string
		log.Error("Config file validation failed.")
		for i, err := range errs {
			log.Errorf("%d -> %s.", i+1, err)
			retErr = append(retErr, err.Error())
		}
		return fmt.Errorf("%s", strings.Join(retErr, "; "))
	}

	return nil
}

func ValidateAndDefaultBackend(state *State) []error {
	var errs []error
	switch {
	case state.Backend == "local":
		if state.Options.StateFile == "" {
			log.Debugf("The stateFile has not been set, default value %s will be used.", local.DefaultStateFile)
			state.Options.StateFile = local.DefaultStateFile
		}
	case state.Backend == "s3":
		if state.Options.Bucket == "" {
			errs = append(errs, fmt.Errorf("state s3 Bucket is empty"))
		}
		if state.Options.Region == "" {
			errs = append(errs, fmt.Errorf("state s3 Region is empty"))
		}
		if state.Options.Key == "" {
			errs = append(errs, fmt.Errorf("state s3 Key is empty"))
		}
	case strings.ToLower(state.Backend) == "k8s" || strings.ToLower(state.Backend) == "kubernetes":
		state.Backend = "k8s"
		if state.Options.Namespace == "" {
			state.Options.Namespace = defaultNamespace
		}
		if state.Options.ConfigMap == "" {
			state.Options.ConfigMap = defaultConfigMapName
		}
	default:
		errs = append(errs, fmt.Errorf("the backend type < %s > is illegal", state.Backend))
	}

	return errs
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
