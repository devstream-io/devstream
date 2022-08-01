package pluginmanager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/internal/pkg/version"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/md5"
)

var errorDev = errors.New("dev version plugins can't be downloaded from the remote plugin repo; please run `make build-plugin.PLUGIN_NAME` to build them locally")

func DownloadPlugins(conf *configloader.Config) error {
	// create plugins dir if not exist
	pluginDir := viper.GetString("plugin-dir")
	if pluginDir == "" {
		return fmt.Errorf("plugins directory should not be \"\"")
	}
	log.Infof("Using dir <%s> to store plugins.", pluginDir)

	// download all plugins that don't exist locally
	dc := NewPbDownloadClient(defaultReleaseUrl)

	for _, tool := range conf.Tools {
		pluginName := configloader.GetPluginName(&tool)
		pluginFileName := configloader.GetPluginFileName(&tool)
		pluginMD5FileName := configloader.GetPluginMD5FileName(&tool)

		// a new line to make outputs more beautiful
		fmt.Println()
		log.Separator(pluginName)

		_, pluginFileErr := os.Stat(filepath.Join(pluginDir, pluginFileName))
		_, pluginMD5FileErr := os.Stat(filepath.Join(pluginDir, pluginMD5FileName))

		// plugin does not exist
		if pluginFileErr != nil {
			if !errors.Is(pluginFileErr, os.ErrNotExist) {
				return pluginFileErr
			}
			// download .so file
			if version.Dev {
				return errorDev
			}
			if err := dc.download(pluginDir, pluginFileName, version.Version); err != nil {
				return err
			}
			log.Successf("[%s] download succeeded.", pluginFileName)
		}
		// .md5 does not exist
		if pluginMD5FileErr != nil {
			if !errors.Is(pluginMD5FileErr, os.ErrNotExist) {
				return pluginMD5FileErr
			}
			// download .md5 file
			if version.Dev {
				return errorDev
			}
			if err := dc.download(pluginDir, pluginMD5FileName, version.Version); err != nil {
				return err
			}
			log.Successf("[%s] download succeeded.", pluginMD5FileName)
		}

		// check if the plugin matches with .md5
		isMD5Match, err := md5.FileMatchesMD5(filepath.Join(pluginDir, pluginFileName), filepath.Join(pluginDir, pluginMD5FileName))
		if err != nil {
			return err
		}

		if !isMD5Match {
			// if existing .so doesn't matches with .md5, re-download
			log.Infof("Plugin: [%s] doesn't match with .md5 and will be downloaded.", pluginFileName)
			if err = redownloadPlugins(dc, pluginDir, pluginFileName, pluginMD5FileName, version.Version); err != nil {
				return err
			}
			// check if the downloaded plugin md5 matches with .md5
			if err = pluginAndMD5Matches(pluginDir, pluginFileName, pluginMD5FileName, tool.Name); err != nil {
				return err
			}
		}

		log.Infof("Initialize [%s] finished.", pluginName)
		log.Separatorf(pluginName)
	}
	return nil
}

// CheckLocalPlugins checks if the local plugins exists, and matches with md5 value.
func CheckLocalPlugins(conf *configloader.Config) error {
	pluginDir := viper.GetString("plugin-dir")
	if pluginDir == "" {
		return fmt.Errorf("plugins directory doesn't exist")
	}

	log.Infof("Using dir <%s> to store plugins.", pluginDir)

	for _, tool := range conf.Tools {
		pluginFileName := configloader.GetPluginFileName(&tool)
		pluginMD5FileName := configloader.GetPluginMD5FileName(&tool)
		if _, err := os.Stat(filepath.Join(pluginDir, pluginFileName)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("plugin %s doesn't exist", tool.Name)
			}
			return err
		}
		if _, err := os.Stat(filepath.Join(pluginDir, pluginMD5FileName)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf(".md5 file of plugin %s doesn't exist", tool.Name)
			}
			return err
		}
		if err := pluginAndMD5Matches(pluginDir, pluginFileName, pluginMD5FileName, tool.Name); err != nil {
			return err
		}
	}
	return nil
}

// pluginAndMD5Matches checks if the plugins match with .md5
func pluginAndMD5Matches(pluginDir, soFileName, md5FileName, tooName string) error {
	isMD5Match, err := md5.FileMatchesMD5(filepath.Join(pluginDir, soFileName), filepath.Join(pluginDir, md5FileName))
	if err != nil {
		return err
	}
	if !isMD5Match {
		return fmt.Errorf("plugin %s doesn't match with .md5", tooName)
	}
	return nil
}

// redownloadPlugins re-download from remote
func redownloadPlugins(dc *PbDownloadClient, pluginDir, pluginFileName, pluginMD5FileName, version string) error {
	if err := os.Remove(filepath.Join(pluginDir, pluginFileName)); err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(pluginDir, pluginMD5FileName)); err != nil {
		return err
	}
	// download .so file
	if err := dc.download(pluginDir, pluginFileName, version); err != nil {
		return err
	}
	log.Successf("[%s] download succeeded.", pluginFileName)
	// download .md5 file
	if err := dc.download(pluginDir, pluginMD5FileName, version); err != nil {
		return err
	}
	log.Successf("[%s] download succeeded.", pluginMD5FileName)
	return nil
}
