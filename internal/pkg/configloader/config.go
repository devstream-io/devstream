package configloader

import (
	"fmt"
	"io/ioutil"
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
}

// Tool is the struct for one section of the DevStream configuration file.
type Tool struct {
	// RFC 1123 - DNS Subdomain Names style
	// contain no more than 253 characters
	// contain only lowercase alphanumeric characters, '-' or '.'
	// start with an alphanumeric character
	// end with an alphanumeric character
	Name      string                 `yaml:"name"`
	Plugin    string                 `yaml:"plugin"`
	DependsOn []string               `yaml:"dependsOn"`
	Options   map[string]interface{} `yaml:"options"`
}

func (t *Tool) DeepCopy() *Tool {
	var retTool = Tool{
		Name:      t.Name,
		Plugin:    t.Plugin,
		DependsOn: t.DependsOn,
		Options:   map[string]interface{}{},
	}
	for k, v := range t.Options {
		retTool.Options[k] = v
	}
	return &retTool
}

// LoadConf reads an input file as a Config struct.
func LoadConf(configFileName, varFileName string) *Config {
	configFileBytes, err := ioutil.ReadFile(configFileName)
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

// GetPluginFileName creates the file name based on the tool's name and version
// If the plugin {githubactions 0.0.1}, the generated name will be "githubactions_0.0.1.so"
func GetPluginFileName(t *Tool) string {
	return fmt.Sprintf("%s-%s-%s_%s.so", t.Plugin, GOOS, GOARCH, version.Version)
}

// GetPluginMD5FileName  If the plugin {githubactions 0.0.1}, the generated name will be "githubactions_0.0.1.md5"
func GetPluginMD5FileName(t *Tool) string {
	return fmt.Sprintf("%s-%s-%s_%s.md5", t.Plugin, GOOS, GOARCH, version.Version)
}

// GetDtmMD5FileName format likes dtm-linux-amd64
func GetDtmMD5FileName() string {
	return fmt.Sprintf("%s-%s-%s.md5", "dtm", GOOS, GOARCH)
}
