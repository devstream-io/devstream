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

	// AppInConfig is the raw structured data in config file.
	// The main difference between App and AppInConfig is "CI" and "CD" field.
	// The "CIRawConfigs" is the raw data of "CI" field defined in config file,
	// which will be rendered to "CIPipelines" field in App with "PipelineTemplates".
	// The "CDRawConfigs" is similar to "CIRawConfigs".
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
