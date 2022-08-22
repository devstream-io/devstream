package jenkins

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
)

type jenkinsOptions struct {
	helm.Options `mapstructure:",squash"`
}

func newOptions(options plugininstaller.RawOptions) (jenkinsOptions, error) {
	var opts jenkinsOptions
	if err := mapstructure.Decode(options, &opts); err != nil {
		return opts, err
	}
	return opts, nil
}

func (opts *jenkinsOptions) encode() (map[string]interface{}, error) {
	var options map[string]interface{}
	if err := mapstructure.Decode(opts, &options); err != nil {
		return nil, err
	}
	return options, nil
}

func setDefaultValue(defaultOpts *helm.Options) plugininstaller.MutableOperation {
	return func(options plugininstaller.RawOptions) (plugininstaller.RawOptions, error) {
		opts, err := newOptions(options)
		if err != nil {
			return nil, err
		}
		opts.FillDefaultValue(defaultOpts)
		return opts.encode()
	}
}
