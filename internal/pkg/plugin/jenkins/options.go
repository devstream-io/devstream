package jenkins

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugin/common/helm"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
)

type jenkinsOptions struct {
	// in test env, we will use hostpath to auto create persistent volume
	TestEnv      bool `mapstructure:"test_env"`
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
