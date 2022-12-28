package configmanager

import (
	"bytes"
	"fmt"
	"regexp"
	"unicode"

	"gopkg.in/yaml.v3"
)

// rawConfig respent every valid config block for devstream
type rawConfig struct {
	apps              []byte
	pipelineTemplates []byte
	tools             []byte
	vars              []byte
	config            []byte
}

// newRawConfigFromConfigBytes will extract rawConfig from file content
func newRawConfigFromConfigBytes(fileText []byte) (*rawConfig, error) {
	topConfigLevelRegex := regexp.MustCompile(`(?m)^\w+\s*:`)
	allMatchIndex := topConfigLevelRegex.FindAllIndex(fileText, -1)
	totalMatchLength := len(allMatchIndex)
	cleanSpaceAndColonFunc := func(r rune) bool {
		return unicode.IsSpace(r) || string(r) == ":"
	}
	rawConfigData := &rawConfig{}
	// get map key and construct rawConfig options
	for i, currentIndex := range allMatchIndex {
		var matchStr []byte
		// if is last item, just get rest string
		if i == totalMatchLength-1 {
			matchStr = fileText[currentIndex[1]:]
		} else {
			nextIndex := allMatchIndex[i+1]
			matchStr = fileText[currentIndex[1]:nextIndex[0]]
		}
		configKey := bytes.TrimFunc(fileText[currentIndex[0]:currentIndex[1]], cleanSpaceAndColonFunc)
		switch string(configKey) {
		case "apps":
			rawConfigData.apps = append(rawConfigData.apps, matchStr...)
		case "tools":
			rawConfigData.tools = append(rawConfigData.tools, matchStr...)
		case "pipelineTemplates":
			rawConfigData.pipelineTemplates = append(rawConfigData.pipelineTemplates, matchStr...)
		case "config":
			if len(rawConfigData.config) != 0 {
				return nil, fmt.Errorf("<config> key can only be defined once")
			}
			rawConfigData.config = append(rawConfigData.config, matchStr...)
		case "vars":
			rawConfigData.vars = append(rawConfigData.vars, matchStr...)
		default:
			const errHint = "you may have filled in the wrong key or imported a yaml file that is not related to dtm"
			return nil, fmt.Errorf("invalid config key <%s>, %s", string(configKey), errHint)
		}
	}
	if err := rawConfigData.validate(); err != nil {
		return nil, err
	}
	return rawConfigData, nil
}

// validate will check config data is valid
func (c *rawConfig) validate() error {
	errorFmt := "config not valid; check the [%s] section of your config file"
	if (len(c.config)) == 0 {
		return fmt.Errorf(errorFmt, "config")
	}
	if len(c.apps) == 0 && len(c.tools) == 0 {
		return fmt.Errorf(errorFmt, "tools and apps")
	}
	return nil
}

// getVars will generate variables from vars config
func (c *rawConfig) getVars() (map[string]any, error) {
	var globalVars map[string]any
	err := yaml.Unmarshal(c.vars, &globalVars)
	return globalVars, err
}

// getToolsWithVars will generate Tools from tools config
func (c *rawConfig) getToolsWithVars(vars map[string]any) (Tools, error) {
	renderedTools, err := renderConfigWithVariables(string(c.tools), vars)
	if err != nil {
		return nil, err
	}
	var tools = make(Tools, 0)
	err = yaml.Unmarshal(renderedTools, &tools)
	return tools, err
}

// getAppToolsWithVars will generate tools from apps and pipelineTemplates
func (c *rawConfig) getAppToolsWithVars(vars map[string]any) (Tools, error) {
	// render and unmarshal app config
	apps, err := c.getAppsWithVars(vars)
	if err != nil {
		return nil, err
	}

	// get pipelineTemplateMap for app ci/cd config
	pipelineTemplateMap, err := c.getTemplatePipelineMap()
	if err != nil {
		return nil, err
	}

	var tools = make(Tools, 0)
	for _, a := range apps {
		appTools, err := a.getTools(vars, pipelineTemplateMap)
		if err != nil {
			return nil, err
		}
		tools = append(tools, appTools...)
	}
	return tools, err
}

// getTemplatePipelineMap will get map of templateName => templateStr
// TODO(steinliber) this func will transfer int to string, change it later
func (c *rawConfig) getTemplatePipelineMap() (map[string]string, error) {
	yamlRegex := regexp.MustCompile(`([^:]+:)(\s*)\'?((\[\[[^\]]+\]\][^\s\[]*)+)\'?[^\s#\n]*`)
	pipelineTemplateStr := yamlRegex.ReplaceAll(c.pipelineTemplates, []byte("$1$2\"$3\""))
	var pipelineTemplates = make([]*pipelineTemplate, 0)
	if err := yaml.Unmarshal(pipelineTemplateStr, &pipelineTemplates); err != nil {
		return nil, err
	}
	pipelineTemplateMap := map[string]string{}
	for _, t := range pipelineTemplates {
		rawPipeline, err := yaml.Marshal(t)
		if err != nil {
			return nil, err
		}
		if _, ok := pipelineTemplateMap[t.Name]; ok {
			return nil, fmt.Errorf("pipelineTemplate <%s> is duplicated", t.Name)
		}
		pipelineTemplateMap[t.Name] = string(rawPipeline)
	}
	return pipelineTemplateMap, nil
}

// getConfig will get config options
func (c *rawConfig) getConfig() (*CoreConfig, error) {
	var config *CoreConfig
	err := yaml.Unmarshal(c.config, &config)
	return config, err
}

// getAppsWithVars will get apps struct array
func (c *rawConfig) getAppsWithVars(vars map[string]any) ([]*app, error) {
	renderedApps, err := renderConfigWithVariables(string(c.apps), vars)
	if err != nil {
		return nil, err
	}
	var apps = make([]*app, 0)
	err = yaml.Unmarshal(renderedApps, &apps)
	return apps, err
}
