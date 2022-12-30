package jira

import (
	_ "embed"
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile/server"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

//go:embed tpl/issuebuilder.tpl.yml
var issueBuilderWorkflow string

func validate(rawOptions configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := newOptions(rawOptions)
	if err != nil {
		return nil, err
	}
	if err := validator.CheckStructError(opts).Combine(); err != nil {
		return nil, err
	}
	return rawOptions, nil
}

func setDefault(rawOptions configmanager.RawOptions) (configmanager.RawOptions, error) {
	const workflowFileName = "jiraIssueBuilder.yml"
	opts, err := newOptions(rawOptions)
	if err != nil {
		return nil, err
	}
	opts.CIFileConfig = &cifile.CIFileConfig{
		Type: server.CIGithubType,
		ConfigContentMap: cifile.CIFileConfigMap{
			fmt.Sprintf("%s/%s", server.CiGitHubWorkConfigLocation, workflowFileName): issueBuilderWorkflow,
		},
		Vars: cifile.CIFileVarsMap(rawOptions),
	}
	return mapz.DecodeStructToMap(opts)
}
