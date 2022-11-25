package common

import "github.com/devstream-io/devstream/pkg/util/log"

func HandleErrLogsWithPlugin(err error, pluginName string) {
	log.Errorf("An error have occurred while processing the plugin \"%s\" -> %s", pluginName, err)
}
