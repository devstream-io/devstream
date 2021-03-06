package jenkinsgithub

import (
	"context"
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

var plugins = []string{
	"ghprb",                     // "GitHub Pull Request Builder" plugin, https://plugins.jenkins.io/ghprb/
	"antisamy-markup-formatter", // "OWASP Markup Formatter" plugin, https://plugins.jenkins.io/antisamy-markup-formatter/
}

func getPluginExistsMap(j *jenkins.Jenkins) (map[string]bool, error) {
	res := make(map[string]bool)
	for _, pluginName := range plugins {
		plugin, err := j.HasPlugin(context.Background(), pluginName)
		if err != nil {
			return nil, fmt.Errorf("failed to check plugin %s: %s", pluginName, err)
		}
		res[pluginName] = plugin != nil
	}

	return res, nil
}

func installPluginsIfNotExists(j *jenkins.Jenkins) error {
	hasPlugins, err := getPluginExistsMap(j)
	if err != nil {
		return err
	}
	for _, pluginName := range plugins {
		if hasPlugins[pluginName] {
			continue
		}
		if err := j.InstallPlugin(context.Background(), pluginName, "latest"); err != nil {
			return fmt.Errorf("failed to install plugin %s: %s", pluginName, err)
		}
	}

	return nil
}
