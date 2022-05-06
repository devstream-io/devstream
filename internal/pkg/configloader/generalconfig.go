package configloader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// GeneralConfig is the struct for loading State and configuration YAML files.
type GeneralConfig struct {
	VarFile  string `yaml:"varFile"`
	ToolFile string `yaml:"toolFile"`
	State    *State
}

type State struct {
	Backend string                 `yaml:"backend"`
	Options map[string]interface{} `yaml:"options"`
}

// LoadGeneralConf reads an input file as a GeneralConfig struct.
func LoadGeneralConf(configFileName string) (*GeneralConfig, error) {
	configFileBytes, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Error(err)
		log.Info("Maybe the default file doesn't exist or you forgot to pass your config file to the \"-f\" option?")
		log.Info("See \"dtm help\" for more information.")
		return nil, err
	}
	log.Debugf("Original general config: \n%s\n", string(configFileBytes))

	var gConfig GeneralConfig
	err = yaml.Unmarshal(configFileBytes, &gConfig)
	if err != nil {
		log.Error("Please verify the format of your general config file.")
		log.Errorf("Reading general config file failed. %s.", err)
		return nil, err
	}

	errs := validateGeneralConfig(&gConfig)
	if len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("General config validation failed: %s.", e)
		}
		return nil, nil
	}

	absToolFilePath, err := parseCustomPath(configFileName, gConfig.ToolFile)
	if err != nil {
		return nil, err
	}
	gConfig.ToolFile = absToolFilePath

	if gConfig.VarFile != "" {
		absVarFilePath, err := parseCustomPath(configFileName, gConfig.VarFile)
		if err != nil {
			return nil, err
		}
		gConfig.VarFile = absVarFilePath
		return &gConfig, nil
	}

	gConfig.VarFile = "variables.yaml"
	return &gConfig, nil
}

// parseCustomPath  parse the path of tools.yaml or variable.yaml
func parseCustomPath(configFileName, customPath string) (string, error) {
	if filepath.IsAbs(customPath) {
		log.Debugf("Abs path is %s.", customPath)
		if err := fileExists(customPath); err != nil {
			log.Errorf("file %s not exists.", customPath)
			return "", err
		}
		return customPath, nil
	}
	configFilePath, err := filepath.Abs(configFileName)
	if err != nil {
		return "", err
	}
	absFilePath := filepath.Join(filepath.Dir(configFilePath), customPath)
	log.Debugf("Abs path is %s.", absFilePath)
	if err := fileExists(absFilePath); err != nil {
		log.Errorf("file %s not exists.", absFilePath)
		return "", err
	}
	return absFilePath, nil
}

func fileExists(path string) error {
	if _, err := os.Stat(path); err != nil {
		return err
	}
	return nil
}

// validateGeneralConfig validate all the general config items
func validateGeneralConfig(c *GeneralConfig) []error {
	errors := make([]error, 0)

	if c.ToolFile == "" {
		errors = append(errors, fmt.Errorf("tool file is empty"))
	}

	if c.State == nil {
		errors = append(errors, fmt.Errorf("state config is empty"))
	}

	if c.State.Options == nil {
		errors = append(errors, fmt.Errorf("state options is empty"))
	}

	if c.State.Backend == "" {
		errors = append(errors, fmt.Errorf("backend is empty"))
	}
	return errors
}
