package gitlabci

import "github.com/devstream-io/devstream/internal/pkg/plugin/installer/ci"

type options struct {
	Runner      *runnerOptions `mapstructure:"runner"`
	ci.CIConfig `mapstructure:",squash"`
}

type runnerOptions struct {
	Enable bool `mapstructure:"enable"`
}

func (o *options) needCreateRunner() bool {
	if o.Runner != nil {
		return o.Runner.Enable
	}
	return false
}
