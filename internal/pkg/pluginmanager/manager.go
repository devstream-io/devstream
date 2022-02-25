package pluginmanager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/pkg/util/log"
)

func DownloadPlugins(conf *configloader.Config) error {
	// create plugins dir if not exist
	pluginDir := viper.GetString("plugin-dir")
	if pluginDir == "" {
		return fmt.Errorf("plugins directory should not be \"\"")
	}
	log.Infof("Using dir <%s> to store plugins.", pluginDir)

	// download all plugins that don't exist locally
	dc := NewPbDownloadClient()

	for _, tool := range conf.Tools {
		pluginFileName := configloader.GetPluginFileName(&tool)
		if _, err := os.Stat(filepath.Join(pluginDir, pluginFileName)); errors.Is(err, os.ErrNotExist) {
			// plugin does not exist
			err := dc.download(pluginDir, pluginFileName, tool.Plugin.Version)
			if err != nil {
				return err
			}
			continue
		}
		// check md5
		dup, err := checkFileMD5(filepath.Join(pluginDir, pluginFileName), dc, pluginFileName, tool.Plugin.Version)
		if err != nil {
			return err
		}

		if dup {
			log.Infof("Plugin: %s already exists, no need to download.", pluginFileName)
			continue
		}

		log.Infof("Plugin: %s changed and will be overrided.", pluginFileName)
		if err = os.Remove(filepath.Join(pluginDir, pluginFileName)); err != nil {
			return err
		}
		if err = dc.download(pluginDir, pluginFileName, tool.Plugin.Version); err != nil {
			return err
		}
	}

	return nil
}

func CheckLocalPlugins(conf *configloader.Config) error {
	pluginDir := viper.GetString("plugin-dir")
	if pluginDir == "" {
		return fmt.Errorf("plugins directory doesn't exist")
	}

	log.Infof("Using dir <%s> to store plugins.", pluginDir)

	for _, tool := range conf.Tools {
		pluginFileName := configloader.GetPluginFileName(&tool)
		if _, err := os.Stat(filepath.Join(pluginDir, pluginFileName)); errors.Is(err, os.ErrNotExist) {
			if err != nil {
				return err
			}
			return fmt.Errorf("plugin %s doesn't exist", tool.Name)
		}
	}

	return nil
}

func checkFileMD5(file string, dc *PbDownloadClient, pluginFileName string, version string) (bool, error) {
	localmd5, err := LocalContentMD5(file)
	if err != nil {
		return false, err
	}
	remotemd5, err := dc.fetchContentMD5(pluginFileName, version)
	if err != nil {
		return false, err
	}

	if localmd5 == remotemd5 {
		return true, nil
	}
	return false, nil
}
