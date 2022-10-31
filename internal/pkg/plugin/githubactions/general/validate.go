package general

import (
	"fmt"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// validate validates the options provided by the core.
func validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := newActionOptions(options)
	if err != nil {
		return nil, err
	}
	if err = validator.StructAllError(opts); err != nil {
		return nil, err
	}

	// check repo is valid
	if opts.ProjectRepo.RepoType == "gitlab" {
		return nil, fmt.Errorf("github action don't support gitlab repo")
	}
	if err := opts.ProjectRepo.CheckValid(); err != nil {
		log.Debugf("github action validate repo invalid: %+v", err)
		return nil, err
	}
	return options, nil
}

func setDefault(option configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := newActionOptions(option)
	if err != nil {
		return nil, err
	}
	projectRepo, err := opts.SCM.BuildRepoInfo()
	if err != nil {
		return nil, err
	}
	opts.ProjectRepo = projectRepo
	opts.CIConfig = opts.Action.buildCIConfig(projectRepo)
	return mapz.DecodeStructToMap(opts)
}
