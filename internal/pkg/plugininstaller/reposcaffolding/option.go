package reposcaffolding

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/common"
)

type Options struct {
	SourceRepo      *common.Repo `validate:"required" mapstructure:"sourceRepo"`
	DestinationRepo *common.Repo `validate:"required" mapstructure:"destinationRepo"`
	Vars            map[string]interface{}
}

func NewOptions(options plugininstaller.RawOptions) (*Options, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}

func (opts *Options) Encode() (map[string]interface{}, error) {
	var options map[string]interface{}
	if err := mapstructure.Decode(opts, &options); err != nil {
		return nil, err
	}
	return options, nil
}

func (opts *Options) renderTplConfig() map[string]interface{} {
	renderConfig := opts.DestinationRepo.BuildRepoRenderConfig()
	for k, v := range opts.Vars {
		renderConfig[k] = v
	}
	return renderConfig
}
