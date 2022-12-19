package trello

import (
	"fmt"

	_ "embed"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile/server"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/types"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

//go:embed tpl/issuebuilder.tpl.yml
var issueBuilderWorkflow string

func validate(rawOptions configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := newOptions(rawOptions)
	if err != nil {
		return nil, err
	}
	if err := validator.StructAllError(opts); err != nil {
		return nil, err
	}
	return rawOptions, nil
}

func setDefault(rawOptions configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := newOptions(rawOptions)
	if err != nil {
		return nil, err
	}
	// set board dedefault value
	var boardDefaultConfig = &board{
		Name:        fmt.Sprintf("%s/%s", opts.Scm.GetRepoOwner(), opts.Scm.GetRepoName()),
		Description: fmt.Sprintf("Description is managed by DevStream, please don't modify. %s/%s", opts.Scm.GetRepoOwner(), opts.Scm.GetRepoName()),
	}
	if opts.Board == nil {
		opts.Board = &board{}
	}
	types.FillStructDefaultValue(opts.Board, boardDefaultConfig)

	// set CIFileConfig
	const workflowFileName = "trelloIssueBuilder.yml"
	opts.CIFileConfig = &cifile.CIFileConfig{
		Type: server.CIGithubType,
		ConfigContentMap: cifile.CIFileConfigMap{
			fmt.Sprintf("%s/%s", server.CiGitHubWorkConfigLocation, workflowFileName): issueBuilderWorkflow,
		},
	}
	return mapz.DecodeStructToMap(opts)
}
