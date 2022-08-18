package jenkins

import (
	"fmt"
	"strings"

	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/helm"
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

// refer: https://github.com/jenkinsci/helm-charts/blob/30766b45faf639dbad45e2c66330ee5fcdc7e37f/charts/jenkins/templates/_helpers.tpl#L46-L57
func (opts *jenkinsOptions) getJenkinsFullName() string {
	if strings.Contains(opts.Chart.ChartName, opts.Chart.ReleaseName) {
		return opts.Chart.ReleaseName
	}

	return fmt.Sprintf("%s-%s", opts.Chart.ReleaseName, opts.Chart.ChartName)
}
