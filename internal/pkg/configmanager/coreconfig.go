package configmanager

import (
	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/file"
)

type CoreConfig struct {
	State *State `yaml:"state"`
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
	// for k8s backend
	Namespace string `yaml:"namespace"`
	ConfigMap string `yaml:"configmap"`
}

func getCoreConfigFromConfigFile(fileBytes []byte) (*CoreConfig, error) {
	yamlPath := "$.config"
	yamlStr, err := file.GetYamlNodeStrByPath(fileBytes, yamlPath)
	if err != nil {
		return nil, err
	}

	var config *CoreConfig
	err = yaml.Unmarshal([]byte(yamlStr), &config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
