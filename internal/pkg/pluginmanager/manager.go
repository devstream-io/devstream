package pluginmanager

import (
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/merico-dev/stream/internal/pkg/configloader"
)

func DownloadPlugins(conf *configloader.Config) error {
	// create plugins dir if not exist
	pluginsDir := filepath.Join(".", "plugins")

	// download all plugins that don't exist locally
	dc := NewPbDownloadClient()

	for _, tool := range conf.Tools {
		pluginFileName := configloader.GetPluginFileName(&tool)
		if _, err := os.Stat(filepath.Join(pluginsDir, pluginFileName)); errors.Is(err, os.ErrNotExist) {
			// plugin does not exist
			log.Printf("=== now downloading plugin: %s ,version: %s === \n", pluginFileName, tool.Version)
			err := dc.download(pluginsDir, pluginFileName, tool.Version)
			if err != nil {
				return err
			}
			log.Printf("=== plugin: %s ,version: %s downloaded === \n", pluginFileName, tool.Version)
		}
	}

	return nil
}
