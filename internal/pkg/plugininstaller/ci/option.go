package ci

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
	"github.com/devstream-io/devstream/pkg/util/types"
)

type Options struct {
	CIConfig    *CIConfig    `mapstructure:"ci" validate:"required"`
	ProjectRepo *common.Repo `mapstructure:"projectRepo" validate:"required"`
}

func NewOptions(options plugininstaller.RawOptions) (*Options, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

// FillDefaultValue config options default values by input defaultOptions
func (opts *Options) FillDefaultValue(defaultOptions *Options) {
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

func (opts *Options) Encode() (map[string]interface{}, error) {
	var options map[string]interface{}
	if err := mapstructure.Decode(opts, &options); err != nil {
		return nil, err
	}
	return options, nil
}
