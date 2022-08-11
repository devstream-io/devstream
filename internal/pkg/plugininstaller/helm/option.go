package helm

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/pkg/util/helm"
)

// Options is the struct for parameters used by the helm install config.
type Options struct {
	CreateNamespace *bool      `mapstructure:"create_namespace"`
	Repo            helm.Repo  `mapstructure:"repo"`
	Chart           helm.Chart `mapstructure:"chart"`
}

func (opts *Options) GetHelmParam() *helm.HelmParam {
	return &helm.HelmParam{
		Repo:  opts.Repo,
		Chart: opts.Chart,
	}
}

func (opts *Options) CheckIfCreateNamespace() bool {
	if opts.CreateNamespace == nil {
		return false
	}
	return *opts.CreateNamespace
}

func (opts *Options) GetNamespace() string {
	return opts.Chart.Namespace
}

func (opts *Options) GetReleaseName() string {
	return opts.Chart.ReleaseName
}

func (opts *Options) Encode() (map[string]interface{}, error) {
	var options map[string]interface{}
	if err := mapstructure.Decode(opts, &options); err != nil {
		return nil, err
	}
	return options, nil
}

func (opts *Options) fillDefaultValue(defaultOpts *Options) {
	chart := &opts.Chart
	chart.FillDefaultValue(&defaultOpts.Chart)
	repo := &opts.Repo
	repo.FillDefaultValue(&defaultOpts.Repo)
	if opts.CreateNamespace == nil && defaultOpts.CreateNamespace != nil {
		opts.CreateNamespace = defaultOpts.CreateNamespace
	}
}

// NewOptions create options by raw options
func NewOptions(options plugininstaller.RawOptions) (*Options, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}
