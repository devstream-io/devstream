package gitlabci

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile/server"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/util"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validate validates the options provided by the core.
func validate(rawOptions configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts := new(options)
	if err := util.DecodePlugin(rawOptions, opts); err != nil {
		return nil, err
	}
	// check struct data
	if err := validator.CheckStructError(opts).Combine(); err != nil {
		return nil, err
	}

	// check repo is valid
	if opts.ProjectRepo.RepoType != "gitlab" {
		return nil, fmt.Errorf("gitlab ci only support gitlab repo")
	}
	return rawOptions, nil
}

func setDefault(rawOptions configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts := new(options)
	if err := util.DecodePlugin(rawOptions, opts); err != nil {
		return nil, err
	}
	// set default value of pipeline location
	if opts.Pipeline.ConfigLocation == "" {
		opts.Pipeline.ConfigLocation = "https://raw.githubusercontent.com/devstream-io/dtm-pipeline-templates/main/gitlab-ci/.gitlab-ci.yml"
	}
	opts.CIFileConfig = opts.Pipeline.BuildCIFileConfig(server.CIGitLabType, opts.ProjectRepo)
	return mapz.DecodeStructToMap(opts)
}
