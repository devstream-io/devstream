package ci

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci/cifile/server"
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
	return options, nil
}

// SetSCMDefault is used for gitlab/github to set default options in ci
func SetDefault(ciType server.CIServerType) func(option configmanager.RawOptions) (configmanager.RawOptions, error) {
	return func(option configmanager.RawOptions) (configmanager.RawOptions, error) {
		opts, err := NewCIOptions(option)
		if err != nil {
			return nil, err
		}
		// set default value of repoInfo
		if err := opts.ProjectRepo.SetDefault(); err != nil {
			return nil, err
		}
		opts.CIFileConfig = opts.Pipeline.BuildCIFileConfig(ciType, opts.ProjectRepo)
		return mapz.DecodeStructToMap(opts)
	}
}
