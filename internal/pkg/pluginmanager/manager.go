package pluginmanager

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/merico-dev/stream/internal/pkg/configloader"
)

func DownloadPlugins(conf *configloader.Config) error {
	// create plugins dir if not exist
	pluginDir := viper.GetString("plugin-dir")
	if pluginDir == "" {
		return fmt.Errorf("plugins directory should not be \"\"")
	}
	log.Printf("Using dir <%s> to store plugins.", pluginDir)

	// download all plugins that don't exist locally
	dc := NewPbDownloadClient()

	for _, tool := range conf.Tools {
		pluginFileName := configloader.GetPluginFileName(&tool)
		if _, err := os.Stat(filepath.Join(pluginDir, pluginFileName)); errors.Is(err, os.ErrNotExist) {
			// plugin does not exist
			err := dc.download(pluginDir, pluginFileName, tool.Version)
			if err != nil {
				return err
			}
			continue
		}
		log.Printf("Plugin: %s already exists, no need to download.", pluginFileName)
	}

	return nil
}
