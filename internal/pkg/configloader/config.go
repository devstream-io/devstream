package configloader

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

// Config is the struct for loading DevStream configuration YAML files.
type Config struct {
	Tools []Tool `yaml:"tools"`
}

// Tool is the struct for one section of the DevStream configuration file.
type Tool struct {
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
	if fname != "" {
		viper.SetConfigFile(fname)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Print(err)
		log.Print("Perhaps you forgot to specify the path of the config file by using the \"-f\" parameter?")
		log.Fatal("See more help by running \"dtm help\"")
	}

	var tools = make([]Tool, 0)

	if err := viper.UnmarshalKey("tools", &tools); err != nil {
		log.Fatal(err)
	}

	return &Config{Tools: tools}
}

// GetPluginFileName creates the file name based on the tool's name and version
// If the plugin {githubactions 0.0.1}, the generated name will be "githubactions_0.0.1.so"
func GetPluginFileName(t *Tool) string {
	return fmt.Sprintf("%s_%s.so", t.Plugin.Kind, t.Plugin.Version)
}
