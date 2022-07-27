package jenkinsgithub

import (
	"context"
	"fmt"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

var plugins = [2]string{
	"ghprb",
	"antisamy-markup-formatter",
}

func getPluginExistsMap(j *jenkins.Jenkins) (map[string]bool, error) {
	res := make(map[string]bool)
	for _, pluginName := range plugins {
		plugin, err := j.HasPlugin(context.Background(), pluginName)
		if err != nil {
			return nil, fmt.Errorf("failed to check plugin %s: %s", pluginName, err)
		}
		if plugin != nil {
			res[pluginName] = true
		} else {
			res[pluginName] = false
		}
	}

	return res, nil
}

func installPluginsIfNotExists(j *jenkins.Jenkins) error {
	hasPlugins, err := getPluginExistsMap(j)
	if err != nil {
		return err
	}
	for _, pluginName := range plugins {
		if !hasPlugins[pluginName] {
			if err := j.InstallPlugin(context.Background(), pluginName, "latest"); err != nil {
				return fmt.Errorf("failed to install plugin %s: %s", pluginName, err)
			}
		}
	}
	return nil
}
