package jira

import (
	_ "embed"

	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

//go:embed tpl/githubIssueBuilder.tpl.yaml
var githubIssueBuilderContent string

// validate validates the options provided by the core.
func validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := newOptions(options)
	if err != nil {
		return nil, err
	}
	if err := validator.StructAllError(opts); err != nil {
		log.Errorf("Options error : %s.", err)
		return nil, fmt.Errorf("opts are illegal")
	}
	if !opts.Scm.IsGithubRepo() {
		return nil, fmt.Errorf("plugin jira-integ only support scm type github for now")
	}
	return options, nil
}

func setDefault(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	const jiraGithubLocation = ".github/workflows/jira-github-integ.yml"
	opts, err := newOptions(options)
	if err != nil {
		return nil, err
	}
	opts.CIFileConfig = &cifile.CIFileConfig{
		Type:             "github",
		ConfigContentMap: cifile.CIFileConfigMap{jiraGithubLocation: githubIssueBuilderContent},
		Vars:             cifile.CIFileVarsMap(options),
	}
	return mapz.DecodeStructToMap(opts)
}
