package ci

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/ci/cifile/server"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/mapz"
	"github.com/devstream-io/devstream/pkg/util/validator"
)

// Validate validates the options provided by the core.
func Validate(options configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := NewCIOptions(options)
	if err != nil {
		return nil, err
	}
	if err = validator.StructAllError(opts); err != nil {
		return nil, err
	}

	if err := opts.ProjectRepo.CheckValid(); err != nil {
		log.Debugf("github action validate repo invalid: %+v", err)
		return nil, err
	}
	return options, nil
}

// SetSCMDefault is used for gitlab/github to set default options in ci
func SetSCMDefault(option configmanager.RawOptions) (configmanager.RawOptions, error) {
	opts, err := NewCIOptions(option)
	if err != nil {
		return nil, err
	}
	projectRepo, err := opts.SCM.BuildRepoInfo()
	if err != nil {
		return nil, err
	}
	opts.ProjectRepo = projectRepo
	ciType := server.CIServerType(projectRepo.RepoType)
	opts.CIFileConfig = opts.Pipeline.BuildCIFileConfig(ciType, projectRepo)
	return mapz.DecodeStructToMap(opts)
}
