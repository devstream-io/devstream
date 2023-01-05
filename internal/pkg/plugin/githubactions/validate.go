package general

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile/server"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/util"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validate validates the options provided by the core.
func validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts := new(ci.CIConfig)
	if err := util.DecodePlugin(options, opts); err != nil {
		return nil, err
	}
	// check struct data
	if err := validator.CheckStructError(opts).Combine(); err != nil {
		return nil, err
	}

	// check repo is valid
	if opts.ProjectRepo.RepoType != "github" {
		return nil, fmt.Errorf("github action only support github repo")
	}
	return options, nil
}

func setDefault(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts := new(ci.CIConfig)
	if err := util.DecodePlugin(options, opts); err != nil {
		return nil, err
	}
	// set default value of pipeline location
	if opts.Pipeline.ConfigLocation == "" {
		opts.Pipeline.ConfigLocation = "https://raw.githubusercontent.com/devstream-io/dtm-pipeline-templates/main/github-actions/workflows/main.yml"
	}
	opts.CIFileConfig = opts.Pipeline.BuildCIFileConfig(server.CIGithubType, opts.ProjectRepo)
	return mapz.DecodeStructToMap(opts)
}
