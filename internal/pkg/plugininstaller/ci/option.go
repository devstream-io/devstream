package ci

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/scm/git"
	"github.com/devstream-io/devstream/pkg/util/types"
)

type Options struct {
	CIConfig    *CIConfig     `mapstructure:"ci" validate:"required"`
	ProjectRepo *git.RepoInfo `mapstructure:"projectRepo" validate:"required"`
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
	if opts.CIConfig == nil {
		opts.CIConfig = defaultOptions.CIConfig
	} else {
		types.FillStructDefaultValue(opts.CIConfig, defaultOptions.CIConfig)
	}
	if opts.ProjectRepo == nil {
		opts.ProjectRepo = defaultOptions.ProjectRepo
	} else {
		types.FillStructDefaultValue(opts.ProjectRepo, defaultOptions.ProjectRepo)
	}
}
