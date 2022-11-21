package helm

import (
	"github.com/mitchellh/mapstructure"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/helm"
)

// Options is the struct for parameters used by the helm install config.
type Options struct {
	InstanceID string     `mapstructure:"instanceID"`
	Repo       helm.Repo  `mapstructure:"repo"`
	Chart      helm.Chart `mapstructure:"chart"`
	ValuesYaml string     `mapstructure:"valuesYaml" validate:"yaml"`
}

func (opts *Options) GetHelmParam() *helm.HelmParam {
	return &helm.HelmParam{
		Repo:  opts.Repo,
		Chart: opts.Chart,
	}
}

func (opts *Options) GetNamespace() string {
	return opts.Chart.Namespace
}

func (opts *Options) GetReleaseName() string {
	return opts.Chart.ReleaseName
}

func (opts *Options) FillDefaultValue(defaultOpts *Options) {
	chart := &opts.Chart
	chart.FillDefaultValue(&defaultOpts.Chart)
	repo := &opts.Repo
	repo.FillDefaultValue(&defaultOpts.Repo)
}

// NewOptions create options by raw options
func NewOptions(options configmanager.RawOptions) (*Options, error) {
	var opts Options
	if err := mapstructure.Decode(options, &opts); err != nil {
		return nil, err
	}
	return &opts, nil
}
