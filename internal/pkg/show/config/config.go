package config

import (
	"fmt"

	"github.com/spf13/viper"
)

var templates = map[string]string{
	"quickstart": QuickStart,
	"gitops":     GitOps,
	"apps":       Apps,
}

//go:generate go run gen_embed_var.go
func Show() error {
	// at first, check is template arg is set
	template := viper.GetString("template")
	if template != "" {
		if tmpl, ok := templates[template]; ok {
			fmt.Println(tmpl)
			return nil
		}
		return fmt.Errorf("illegal template name : < %s >", template)
	}

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
