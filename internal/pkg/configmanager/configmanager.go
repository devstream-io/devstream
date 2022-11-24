package configmanager

import (
	"errors"
	"io"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// Manager is used to load the config file from the ConfigFilePath and finally get the Config object.
type Manager struct {
	ConfigFilePath string
}

// NewManager takes configFilePath, then return a *Manager object.
func NewManager(configFilePath string) *Manager {
	return &Manager{
		ConfigFilePath: configFilePath,
	}
}

// LoadConfig is the only method that the caller of Manager needs to be concerned with, and this method returns a *Config finally.
// The main workflow of this method is:
// 1. Get the original config from the config file specified by ConfigFilePath;
// 2. Parsing tools from the apps, and merge that tools to Config.Tools;
// 3. Validation.
func (m *Manager) LoadConfig() (*Config, error) {
	// step 1
	c, err := m.getConfigFromFile()
	if err != nil {
		return nil, err
	}

	// step 2
	appTools, err := c.getToolsFromApps()
	if err != nil {
		return nil, err
	}

	err = c.renderToolsWithVars()
	if err != nil {
		return nil, err
	}
	c.Tools = append(c.Tools, appTools...)
	c.renderInstanceIDtoOptions()

	// step 3
	if err = c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

// getConfigFromFile gets Config from the config file specified by Manager.ConfigFilePath
func (m *Manager) getConfigFromFile() (*Config, error) {
	configBytes, err := os.ReadFile(m.ConfigFilePath)
	if err != nil {
		return nil, err
	}

	configBytesEscaped := escapeBrackets(configBytes)

	var c Config
	if err = yaml.Unmarshal(configBytesEscaped, &c); err != nil && !errors.Is(err, io.EOF) {
		log.Errorf("Please verify the format of your config. Error: %s.", err)
		return nil, err
	}

	return &c, nil
}

// escapeBrackets is used to escape []byte(": [[xxx]]xxx\n") to []byte(": \"[[xxx]]\"xxx\n")
func escapeBrackets(param []byte) []byte {
	re := regexp.MustCompile(`([^:]+:)(\s*)(\[\[[^\]]+\]\][^\s]*)`)
	return re.ReplaceAll(param, []byte("$1$2\"$3\""))
}
