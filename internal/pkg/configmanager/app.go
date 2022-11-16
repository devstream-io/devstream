package configmanager

type (
	App struct {
		Name         string             `yaml:"name" mapstructure:"name"`
		Spec         RawOptions         `yaml:"spec" mapstructure:"spec"`
		Repo         *Repo              `yaml:"repo" mapstructure:"repo"`
		RepoTemplate *RepoTemplate      `yaml:"repoTemplate" mapstructure:"repoTemplate"`
		CIs          []PipelineTemplate `yaml:"ci" mapstructure:"ci"`
		CDs          []PipelineTemplate `yaml:"cd" mapstructure:"cd"`
	}
	Apps []App

	// AppInConfig is the raw structured data in config file
	AppInConfig struct {
		Name         string        `yaml:"name" mapstructure:"name"`
		Spec         RawOptions    `yaml:"spec" mapstructure:"spec"`
		Repo         *Repo         `yaml:"repo" mapstructure:"repo"`
		RepoTemplate *RepoTemplate `yaml:"repoTemplate" mapstructure:"repoTemplate"`
		CIs          []CICD        `yaml:"ci" mapstructure:"ci"`
		CDs          []CICD        `yaml:"cd" mapstructure:"cd"`
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
	// set "app.repo.name = app.name" if "app.repo.name" is empty
	if app.Repo.RepoCommon != nil && app.Repo.RepoCommon.Name == "" {
		app.Repo.RepoCommon.Name = app.Name
	}

	// fill and validate "app.repo"
	err := app.Repo.FillAndValidate()
	if err != nil {
		return err
	}

	// fill and validate "app.repoTemplate"
	if app.RepoTemplate != nil {
		err = app.RepoTemplate.FillAndValidate()
		if err != nil {
			return err
		}
	}

	return nil
}
