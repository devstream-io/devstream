package configloader

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"

	"github.com/merico-dev/stream/internal/pkg/util/log"
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
	Name    string                 `yaml:"name"`
	Plugin  Plugin                 `yaml:"plugin"`
	Options map[string]interface{} `yaml:"options"`
}

// Plugin is the struct for the plugin section of each tool of the DevStream configuration file.
type Plugin struct {
	Kind    string `mapstructure:"kind"`
	Version string `mapstructure:"version"`
}

func (t *Tool) DeepCopy() *Tool {
	var retTool = Tool{
		Name:    t.Name,
		Plugin:  t.Plugin,
		Options: map[string]interface{}{},
	}
	for k, v := range t.Options {
		retTool.Options[k] = v
	}
	return &retTool
}

// LoadConf reads an input file as a Config struct.
func LoadConf(fname string) *Config {
	fileBytes, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Error(err)
		log.Info("Perhaps you forgot to specify the path of the config file by using the \"-f\" parameter?")
		log.Fatal("See more help by running \"dtm help\"")
	}

	var config Config
	err = yaml.Unmarshal(fileBytes, &config)
	if err != nil {
		log.Errorf("unmarshal config failed: %s", err)
		return nil
	}

	errs := config.Validate()

	if len(errs) != 0 {
		for _, e := range errs {
			fmt.Printf("Config validation failed: %s", e)
		}
		return nil
	}

	return &config
}

// GetPluginFileName creates the file name based on the tool's name and version
// If the plugin {githubactions 0.0.1}, the generated name will be "githubactions_0.0.1.so"
func GetPluginFileName(t *Tool) string {
	return fmt.Sprintf("%s_%s.so", t.Plugin.Kind, t.Plugin.Version)
}
