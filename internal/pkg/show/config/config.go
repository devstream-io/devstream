package config

import (
	"fmt"

	"github.com/spf13/viper"
)

//go:generate go run gen_embed_var.go
func Show() error {
	plugin := viper.GetString("plugin")
	if plugin == "" {
		fmt.Println(DefaultConfig)
		return nil
	}
	if config, ok := pluginDefaultConfigs[plugin]; ok {
		fmt.Println(config)
		return nil
	}
	return fmt.Errorf("illegal plugin name : < %s >", plugin)
}
