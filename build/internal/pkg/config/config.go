package config

import (
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
		log.Fatal(err)
	}

	conf := Config{}
	err = yaml.Unmarshal(f, &conf)
	if err != nil {
		log.Fatal(err)
	}

	return &conf
}
