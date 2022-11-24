package configmanager

import (
	"bytes"
	"errors"
	"io"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/file"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type Manager struct {
	ConfigFilePath string
}

func NewManager(configFilePath string) *Manager {
	return &Manager{
		ConfigFilePath: configFilePath,
	}
}

// LoadConfig reads an input file as a general config.
// It will return "non-nil, nil" or "nil, err".
func (m *Manager) LoadConfig() (*Config, error) {
	rawConf, err := m.getRawConfigFromFile()
	if err != nil {
		return nil, err
	}

	tools, err := rawConf.getAllTools()
	if err != nil {
		return nil, err
	}

	config := &Config{
		State: rawConf.State,
		Tools: tools,
	}
	if err = config.validate(); err != nil {
		return nil, err
	}

	return config, nil
}

// getRawConfigFromFile generate new rawConfig options
func (m *Manager) getRawConfigFromFile() (*rawConfig, error) {
	// 1. get baseDir from configFile
	baseDir, err := file.GetFileAbsDirPath(m.ConfigFilePath)
	if err != nil {
		return nil, err
	}

	// 2. read the original main config file
	configBytes, err := os.ReadFile(m.ConfigFilePath)
	if err != nil {
		return nil, err
	}

	// TODO(daniel-hutao): We should change the documents to delete all "---" with config file. After a while, delete the following line of code, or prompt the user with "This is the wrong way" when "---" be detected

	// replace all "---", otherwise yaml.Unmarshal can only read the content before the first "---"
	configBytes = bytes.Replace(configBytes, []byte("---"), []byte("\n"), -1)

	// 3. decode total yaml files to get rawConfig
	var rawConf rawConfig
	if err := yaml.Unmarshal(configBytes, &rawConf); err != nil && !errors.Is(err, io.EOF) {
		log.Errorf("Please verify the format of your config. Error: %s.", err)
		return nil, err
	}

	rawConf.configFileBaseDir = baseDir
	rawConf.totalConfigBytes = configBytes
	return &rawConf, nil
}
