package configloader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/internal/pkg/version"
	"github.com/devstream-io/devstream/pkg/util/log"
)

var (
	GOOS   string = runtime.GOOS
	GOARCH string = runtime.GOARCH
)

// Config is the struct for loading DevStream configuration YAML files.
type Config struct {
	Tools []Tool `yaml:"tools"`
	State *State
}

// ConfigFile is the struct for loading State and configuration YAML files.
type ConfigFile struct {
	VarFile  string `yaml:"varFile"`
	ToolFile string `yaml:"toolFile"`
	State    *State
}

type State struct {
	Backend string                 `yaml:"backend"`
	Options map[string]interface{} `yaml:"options"`
}

// Tool is the struct for one section of the DevStream configuration file.
type Tool struct {
	Name string `yaml:"name"`
	// RFC 1123 - DNS Subdomain Names style
	// contain no more than 253 characters
	// contain only lowercase alphanumeric characters, '-' or '.'
	// start with an alphanumeric character
	// end with an alphanumeric character
	InstanceID string                 `yaml:"instanceID"`
	DependsOn  []string               `yaml:"dependsOn"`
	Options    map[string]interface{} `yaml:"options"`
}

func (t *Tool) DeepCopy() *Tool {
	var retTool = Tool{
		Name:       t.Name,
		InstanceID: t.InstanceID,
		DependsOn:  t.DependsOn,
		Options:    map[string]interface{}{},
	}
	for k, v := range t.Options {
		retTool.Options[k] = v
	}
	return &retTool
}

// LoadConf reads an input file as a general config.
func LoadConf(configFileName string) (*Config, error) {
	configFileBytes, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Error(err)
		log.Info("Maybe the default file doesn't exist or you forgot to pass your config file to the \"-f\" option?")
		log.Info("See \"dtm help\" for more information.")
		return nil, err
	}
	log.Debugf("Original general config: \n%s\n", string(configFileBytes))

	var gConfig ConfigFile
	err = yaml.Unmarshal(configFileBytes, &gConfig)
	if err != nil {
		log.Error("Please verify the format of your general config file.")
		log.Errorf("Reading general config file failed. %s.", err)
		return nil, err
	}

	errs := validateConfigFile(&gConfig)
	if len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("General config validation failed: %s.", e)
		}
		return nil, nil
	}

	toolFilePath, varFilePath, err := genToolVarPath(configFileName, gConfig)
	if err != nil {
		return nil, err
	}

	cfg := LoadToolConf(toolFilePath, varFilePath)
	cfg.State = gConfig.State
	return cfg, nil
}

// LoadToolConf reads tool file rendering by var file as a Config struct.
func LoadToolConf(toolFileName, varFileName string) *Config {
	configFileBytes, err := ioutil.ReadFile(toolFileName)
	if err != nil {
		log.Error(err)
		log.Info("Maybe the default file doesn't exist or you forgot to pass your config file to the \"-f\" option?")
		log.Info("See \"dtm help\" for more information.")
		return nil
	}

	log.Debugf("Original config: \n%s\n", string(configFileBytes))

	// handle variables in the config file
	configFileBytesWithVarsRendered, err := renderVariables(varFileName, configFileBytes)
	if err != nil {
		log.Error(err)
		return nil
	}

	log.Debugf("Config file after rendering with variables: \n%s\n", string(configFileBytesWithVarsRendered))

	var config Config
	err = yaml.Unmarshal(configFileBytesWithVarsRendered, &config)
	if err != nil {
		log.Error("Please verify the format of your config file.")
		log.Errorf("Reading config file failed. %s.", err)
		return nil
	}

	errs := validateConfig(&config)
	if len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Config validation failed: %s.", e)
		}
		return nil
	}

	return &config
}

// genToolVarPath return the Abs path of tool file and var file, if var file is null, return variables.yaml
func genToolVarPath(configFileName string, gConfig ConfigFile) (string, string, error) {
	var absToolFilePath, absVarFilePath string
	var err error
	absToolFilePath, err = parseCustomPath(configFileName, gConfig.ToolFile)
	if err != nil {
		return "", "", err
	}

	absVarFilePath = "variables.yaml"
	if gConfig.VarFile != "" {
		absVarFilePath, err = parseCustomPath(configFileName, gConfig.VarFile)
		if err != nil {
			return "", "", err
		}
	}

	return absToolFilePath, absVarFilePath, nil
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

// GetPluginFileName creates the file name based on the tool's name and version
// If the plugin {githubactions 0.0.1}, the generated name will be "githubactions_0.0.1.so"
func GetPluginFileName(t *Tool) string {
	return fmt.Sprintf("%s-%s-%s_%s.so", t.Name, GOOS, GOARCH, version.Version)
}

// GetPluginMD5FileName  If the plugin {githubactions 0.0.1}, the generated name will be "githubactions_0.0.1.md5"
func GetPluginMD5FileName(t *Tool) string {
	return fmt.Sprintf("%s-%s-%s_%s.md5", t.Name, GOOS, GOARCH, version.Version)
}

// GetDtmMD5FileName format likes dtm-linux-amd64
func GetDtmMD5FileName() string {
	return fmt.Sprintf("%s-%s-%s.md5", "dtm", GOOS, GOARCH)
}
