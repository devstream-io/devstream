package jenkins

import "github.com/devstream-io/devstream/internal/pkg/plugin/common/helm"

type Options struct {
	// in test env, we will use hostpath to auto create persistent volume
	TestEnv      bool `mapstructure:"test_env"`
	helm.Options `mapstructure:",squash"`
}
