package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/show/config/plugin"
)

var pluginDefaultConfigs = map[string]string{
	"argocd": plugin.ArgocdDefaultConfig,
}

func Show() error {
	plugin := viper.GetString("plugin")
	if config, ok := pluginDefaultConfigs[plugin]; ok {
		fmt.Println(config)
		return nil
	}
	return fmt.Errorf("plugin name error: < %s >", plugin)
}
