package pluginmanager

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/merico-dev/stream/internal/pkg/configloader"
)

func DownloadPlugins(conf *configloader.Config) error {
	// create plugins dir if not exist
	pluginsDir := filepath.Join(".", "plugins")

	// download all plugins that don't exist locally
	client := NewDownloadClient()

	for _, tool := range conf.Tools {
		pluginFileName := configloader.GetPluginFileName(&tool)
		if _, err := os.Stat(filepath.Join(pluginsDir, pluginFileName)); errors.Is(err, os.ErrNotExist) {
			// plugin does not exist
			err := client.download(pluginsDir, pluginFileName, tool.Version)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
