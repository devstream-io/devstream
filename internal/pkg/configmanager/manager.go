package configmanager

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"

	"github.com/devstream-io/devstream/pkg/util/log"
)

type Manager struct {
	ConfigFile string
}

func NewManager(configFileName string) *Manager {
	return &Manager{
		ConfigFile: configFileName,
	}
}

// ConfigRaw is used to describe original raw configs read from files
type ConfigRaw struct {
	VarFile           string             `yaml:"varFile"`
	ToolFile          string             `yaml:"toolFile"`
	AppFile           string             `yaml:"appFile"`
	TemplateFile      string             `yaml:"templateFile"`
	PluginDir         string             `yaml:"pluginDir"`
	State             *State             `yaml:"state"`
	Tools             []Tool             `yaml:"tools"`
	AppsInConfig      []AppInConfig      `yaml:"apps"`
	PipelineTemplates []PipelineTemplate `yaml:"pipelineTemplates"`
	Vars              map[string]any     `yaml:"-"`
}

// LoadConfig reads an input file as a general config.
// It will return "non-nil, err" or "nil, err".
func (m *Manager) LoadConfig() (*Config, error) {
	// 1. get the whole config bytes from all the config files
	configBytesOrigin, err := m.getWholeConfigBytes()
	if err != nil {
		return nil, err
	}

	// get all globals vars
	globalVars, err := getVarsFromConfigBytes(configBytesOrigin)
	if err != nil {
		return nil, err
	}

	// 2. yaml unmarshal to get the whole config
	var config ConfigRaw
	err = yaml.Unmarshal(configBytesOrigin, &config)
	if err != nil {
		log.Errorf("Please verify the format of your config. Error: %s.", err)
		return nil, err
	}

	// 3. render ci/cd templates
	cis := make([][]PipelineTemplate, 0)
	cds := make([][]PipelineTemplate, 0)
	for i, app := range config.AppsInConfig {
		cisOneApp, err := renderCICDFromPipeTemplates(app.CIs, config.PipelineTemplates, globalVars, i, app.Name, "ci")
		if err != nil {
			return nil, err
		}
		cis = append(cis, cisOneApp)

		cdsOneApp, err := renderCICDFromPipeTemplates(app.CDs, config.PipelineTemplates, globalVars, i, app.Name, "cd")
		if err != nil {
			return nil, err
		}
		cds = append(cds, cdsOneApp)
	}

	// remove the pipeline templates, beacuse we don't need them anymore.
	// and because vars here main contains local vars,
	// it will cause error when rendered in the next step if we don't remove them.
	config.PipelineTemplates = nil

	// 4. re-generate the config bytes(without pipeline templates)
	configBytesWithoutPipelineTemplates, err := yaml.Marshal(config)
	if err != nil {
		return nil, err
	}

	// 5. render all vars to the whole config bytes
	renderedConfigBytes, err := renderConfigWithVariables(string(configBytesWithoutPipelineTemplates), globalVars)
	if err != nil {
		return nil, err
	}

	renderedConfigStr := string(renderedConfigBytes)
	log.Debugf("redenered config: %s\n", renderedConfigStr)

	// 6. yaml unmarshal again to get the whole config
	var configRendered ConfigRaw
	//renderedConfigStr := string(renderedConfigBytes)
	//fmt.Println(renderedConfigStr)
	err = yaml.Unmarshal(renderedConfigBytes, &configRendered)
	if err != nil {
		return nil, err
	}
	configFinal := &Config{}
	configFinal.PluginDir = configRendered.PluginDir
	configFinal.State = configRendered.State
	configFinal.Tools = configRendered.Tools

	// 7. restructure the apps
	for i, app := range configRendered.AppsInConfig {
		appFinal := App{
			Name:         app.Name,
			Spec:         app.Spec,
			Repo:         app.Repo,
			RepoTemplate: app.RepoTemplate,
		}
		appFinal.CIs = cis[i]
		appFinal.CDs = cds[i]
		configFinal.Apps = append(configFinal.Apps, appFinal)
	}

	errs := configFinal.Validate()
	if len(errs) > 0 {
		return nil, multierr.Combine(errs...)
	}

	return configFinal, nil
}

func getVarsFromConfigBytes(configBytes []byte) (map[string]any, error) {
	// how to get all the global vars:
	// we regard the whole config bytes as a map,
	// and we regard all the key-value as global vars(except some special keys)
	allVars := make(map[string]any)
	if err := yaml.Unmarshal(configBytes, allVars); err != nil {
		return nil, err
	}

	excludeKeys := []string{
		"varFile", "toolFile", "appFile", "templateFile", "pluginDir",
		"state", "tools", "apps", "pipelineTemplates",
	}

	for _, key := range excludeKeys {
		// delete the special keys
		delete(allVars, key)
	}

	return allVars, nil
}

func mergeMaps(m1 map[string]any, m2 map[string]any) map[string]any {
	all := make(map[string]any)
	for k, v := range m1 {
		all[k] = v
	}
	for k, v := range m2 {
		all[k] = v
	}
	return all
}

func (m *Manager) loadMainConfigFile() ([]byte, error) {
	configBytes, err := os.ReadFile(m.ConfigFile)
	if err != nil {
		log.Errorf("Failed to read the config file. Error: %s", err)
		log.Info(`Maybe the default file (config.yaml) doesn't exist or you forgot to pass your config file to the "-f" option?`)
		log.Info(`See "dtm help" for more information."`)
		return nil, err
	}
	log.Debugf("Original config: \n%s\n", string(configBytes))
	return configBytes, err
}

func (m *Manager) getWholeConfigBytes() ([]byte, error) {
	// 1. read the original main config file
	configBytes, err := m.loadMainConfigFile()
	if err != nil {
		return nil, err
	}

	// 2. yaml unmarshal to get the varFile, toolFile, appFile, templateFile
	var config ConfigRaw
	dec := yaml.NewDecoder(strings.NewReader(string(configBytes)))
	//dec.KnownFields(true)
	for err == nil {
		err = dec.Decode(&config)
	}
	if !errors.Is(err, io.EOF) {
		log.Errorf("Please verify the format of your config. Error: %s.", err)
		return nil, err
	}

	// combine bytes from all files
	if config.ToolFile != "" {
		if configBytes, err = m.combineBytesFromFile(configBytes, config.ToolFile); err != nil {
			return nil, err
		}
	}

	if config.VarFile != "" {
		if configBytes, err = m.combineBytesFromFile(configBytes, config.VarFile); err != nil {
			return nil, err
		}
	}

	if config.AppFile != "" {
		if configBytes, err = m.combineBytesFromFile(configBytes, config.AppFile); err != nil {
			return nil, err
		}
	}

	if config.TemplateFile != "" {
		if configBytes, err = m.combineBytesFromFile(configBytes, config.TemplateFile); err != nil {
			return nil, err
		}
	}

	configBytesStr := string(configBytes)
	log.Debugf("The final whole config is: \n%s\n", configBytesStr)

	return configBytes, nil
}

func (m *Manager) combineBytesFromFile(origin []byte, file string) ([]byte, error) {
	// get main config file path
	mainConfigFileAbs, err := filepath.Abs(m.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("%s not exists", m.ConfigFile)
	}
	// refer other config file path by main config file
	fileAbs, err := genAbsFilePath(filepath.Dir(mainConfigFileAbs), file)
	if err != nil {
		return nil, err
	}
	bytes, err := os.ReadFile(fileAbs)
	if err != nil {
		return nil, err
	}
	return append(origin, bytes...), nil
}

// genAbsFilePath return all of the path with a given file name
func genAbsFilePath(baseDir, file string) (string, error) {
	file = filepath.Join(baseDir, file)

	fileExist := func(path string) bool {
		if _, err := os.Stat(file); err != nil {
			log.Errorf("File %s not exists. Error: %s", file, err)
			return false
		}
		return true
	}

	absFilePath, err := filepath.Abs(file)
	if err != nil {
		log.Errorf(`Failed to get absolute path fo "%s".`, file)
		return "", err
	}
	log.Debugf("Abs path is %s.", absFilePath)
	if fileExist(absFilePath) {
		return absFilePath, nil
	} else {
		return "", fmt.Errorf("file %s not exists", absFilePath)
	}
}
