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
	log.Printf("prepare to use dir < %s > to hold the plugins", pluginDir)

	// download all plugins that don't exist locally
	dc := NewDownloadClient()

	for _, tool := range conf.Tools {
		pluginFileName := configloader.GetPluginFileName(&tool)
		if _, err := os.Stat(filepath.Join(pluginDir, pluginFileName)); errors.Is(err, os.ErrNotExist) {
			// plugin does not exist
			log.Printf("=== downloading plugin: %s ===", pluginFileName)
			err := dc.download(pluginDir, pluginFileName, tool.Version)
			if err != nil {
				return err
			}
			log.Printf("=== plugin: %s has been downloaded ===", pluginFileName)
			continue
		}
		log.Printf("=== plugin: %s exists ===", pluginFileName)
	}

	return nil
}
