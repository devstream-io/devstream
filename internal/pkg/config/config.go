package config

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Tools []Tool `yaml:"tools"`
}

type Tool struct {
	Name    string                 `yaml:"name"`
	Version string                 `yaml:"version"`
	Options map[string]interface{} `yaml:"options"`
}

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
