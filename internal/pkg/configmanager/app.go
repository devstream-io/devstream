package configmanager

type (
	App struct {
		Name         string             `yaml:"name" mapstructure:"name"`
		Spec         RawOptions         `yaml:"spec" mapstructure:"spec"`
		Repo         *Repo              `yaml:"repo" mapstructure:"repo"`
		RepoTemplate *RepoTemplate      `yaml:"repoTemplate" mapstructure:"repoTemplate"`
		CIPipelines  []PipelineTemplate `yaml:"ci" mapstructure:"ci"`
		CDPipelines  []PipelineTemplate `yaml:"cd" mapstructure:"cd"`
	}
	Apps []App

	// AppInConfig is the raw structured data in config file
	AppInConfig struct {
		Name         string        `yaml:"name" mapstructure:"name"`
		Spec         RawOptions    `yaml:"spec" mapstructure:"spec"`
		Repo         *Repo         `yaml:"repo" mapstructure:"repo"`
		RepoTemplate *RepoTemplate `yaml:"repoTemplate" mapstructure:"repoTemplate"`
		CIRawConfigs []CICD        `yaml:"ci" mapstructure:"ci"`
		CDRawConfigs []CICD        `yaml:"cd" mapstructure:"cd"`
	}

	CICD struct {
		Type         string     `yaml:"type" mapstructure:"type"`
		TemplateName string     `yaml:"templateName" mapstructure:"templateName"`
		Options      RawOptions `yaml:"options" mapstructure:"options"`
		Vars         RawOptions `yaml:"vars" mapstructure:"vars"`
	}

	// ConfigRaw is used to describe original raw configs read from files
	ConfigRaw struct {
		VarFile           string             `yaml:"varFile"`
		ToolFile          string             `yaml:"toolFile"`
		AppFile           string             `yaml:"appFile"`
		TemplateFile      string             `yaml:"templateFile"`
		PluginDir         string             `yaml:"pluginDir"`
		State             *State             `yaml:"state"`
		Tools             []Tool             `yaml:"tools"`
		AppsInConfig      []AppInConfig      `yaml:"apps"`
		PipelineTemplates []PipelineTemplate `yaml:"pipelineTemplates"`
		GlobalVars        map[string]any     `yaml:"-"`
	}
)

func (apps Apps) validate() (errs []error) {
	for _, app := range apps {
		if err := app.validate(); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

func (app *App) validate() error {
	if app.Repo.RepoInfo != nil && app.Repo.RepoInfo.Name == "" {
		app.Repo.RepoInfo.Name = app.Name
	}

	err := app.Repo.FillAndValidate()
	if err != nil {
		return err
	}

	if app.RepoTemplate != nil {
		err = app.RepoTemplate.FillAndValidate()
		if err != nil {
			return err
		}
	}

	return nil
}

func (config *ConfigRaw) constructApps(ciPipelines, cdPipelines [][]PipelineTemplate) *Config {
	configFinal := &Config{}
	configFinal.PluginDir = config.PluginDir
	configFinal.State = config.State
	configFinal.Tools = config.Tools

	for i, app := range config.AppsInConfig {
		appFinal := App{
			Name:         app.Name,
			Spec:         app.Spec,
			Repo:         app.Repo,
			RepoTemplate: app.RepoTemplate,
		}
		appFinal.CIPipelines = ciPipelines[i]
		appFinal.CDPipelines = cdPipelines[i]
		configFinal.Apps = append(configFinal.Apps, appFinal)
	}

	return configFinal
}
