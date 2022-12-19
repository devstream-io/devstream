package cifile

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/types"
)

type Options struct {
	CIFileConfig *CIFileConfig `mapstructure:"ci" validate:"required"`
	ProjectRepo  *git.RepoInfo `mapstructure:"scm" validate:"required"`
}

func NewOptions(options configmanager.RawOptions) (*Options, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

// FillDefaultValue config options default values by input defaultOptions
func (opts *Options) fillDefaultValue(defaultOptions *Options) {
	if opts.CIFileConfig == nil {
		opts.CIFileConfig = defaultOptions.CIFileConfig
	} else {
		types.FillStructDefaultValue(opts.CIFileConfig, defaultOptions.CIFileConfig)
	}
	if opts.ProjectRepo == nil {
		opts.ProjectRepo = defaultOptions.ProjectRepo
	} else {
		types.FillStructDefaultValue(opts.ProjectRepo, defaultOptions.ProjectRepo)
	}
}
