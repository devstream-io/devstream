package configloader

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

// Config is the struct for loading DevStream configuration YAML files.
type Config struct {
	Tools []Tool `yaml:"tools"`
}

// Tool is the struct for one section of the DevStream configuration file.
type Tool struct {
	Name    string                 `yaml:"name"`
	Version string                 `yaml:"version"`
	Options map[string]interface{} `yaml:"options"`
}

// LoadConf reads an input file as a Config struct.
func LoadConf(fname string) *Config {
	f, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Print("It seems you don't have a config file. Maybe you have it in another directory and forgot to use the -f option? See dsm -h for more help.")
		log.Fatal(err)
	}

	conf := Config{}
	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		log.Fatal(err)
	}

	return &conf
}

// GetPluginFileName creates the file name based on the tool's name and version
// If tool is {githubactions 0.0.1}, the generated name will be "githubactions_0.0.1.so"
func GetPluginFileName(t *Tool) string {
	return fmt.Sprintf("%s_%s.so", t.Name, t.Version)
}
