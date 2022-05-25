package configloader

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"

	"github.com/devstream-io/devstream/internal/pkg/version"

	"gopkg.in/yaml.v3"

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

// LoadConfig reads an input file as a general config.
func LoadConfig(configFileName string) (*Config, error) {
	// 1. read the original config file
	originalConfigFileBytes, err := ioutil.ReadFile(configFileName)
	if err != nil {
		log.Errorf("Failed to read the config file. Error: %s", err)
		log.Info("Maybe the default file (config.yaml) doesn't exist or you forgot to pass your config file to the \"-f\" option?")
		log.Info("See \"dtm help\" for more information.")
		return nil, err
	}
	log.Debugf("Original general config: \n%s\n", string(originalConfigFileBytes))

	// 2. unmarshal original config file
	var originalConfig OriginalConfig
	err = yaml.Unmarshal(originalConfigFileBytes, &originalConfig)
	if err != nil {
		log.Error("Please verify the format of your config file.")
		log.Errorf("Reading original config file failed. Error: %s.", err)
		return nil, err
	}

	// 3. validation
	errs := validateOriginalConfigFile(&originalConfig)
	if len(errs) != 0 {
		for _, e := range errs {
			log.Errorf("Config file validation failed: %s.", e)
		}
		return nil, nil
	}

	// 4. get the toolFilePath and varFilePath
	toolFilePath, varFilePath, err := genToolVarPath(configFileName, originalConfig)
	if err != nil {
		return nil, err
	}

	// 5. render the tool config & state
	cfg, err := renderToolConfigWithVarsToConfig(toolFilePath, varFilePath)
	if err != nil {
		return nil, err
	}

	cfg.State = originalConfig.State
	return cfg, nil
}

// genToolVarPath return the Abs path of tool file and var file, return (absToolFilePath, "") if var file is empty.
func genToolVarPath(configFileName string, originalConfig OriginalConfig) (string, string, error) {
	var absToolFilePath, absVarFilePath string
	var err error

	absToolFilePath, err = parseCustomPath(configFileName, originalConfig.ToolFile)
	if err != nil {
		return "", "", err
	}
	// if var file is empty, just return ""
	if originalConfig.VarFile != "" {
		absVarFilePath, err = parseCustomPath(configFileName, originalConfig.VarFile)
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
