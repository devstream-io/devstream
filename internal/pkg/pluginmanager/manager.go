package pluginmanager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/internal/pkg/version"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/md5"
)

const defaultReleaseUrl = "https://download.devstream.io"

func DownloadPlugins(tools []configmanager.Tool, pluginDir, os, arch string) error {
	return downloadPlugins(defaultReleaseUrl, tools, pluginDir, os, arch, version.Version)
}

func downloadPlugins(baseURL string, tools []configmanager.Tool, pluginDir, osName, arch, version string) error {
	if pluginDir == "" {
		return fmt.Errorf(`plugins directory should not be ""`)
	}

	log.Infof("Using dir <%s> to store plugins.", pluginDir)

	// download all plugins that don't exist locally
	dc := NewPluginDownloadClient(baseURL)

	for _, tool := range tools {
		pluginName := tool.GetPluginNameWithOSAndArch(osName, arch)
		pluginFileName := tool.GetPluginFileNameWithOSAndArch(osName, arch)
		pluginMD5FileName := tool.GetPluginMD5FileNameWithOSAndArch(osName, arch)

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
			if err := dc.download(pluginDir, pluginFileName, version); err != nil {
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
			if err := dc.download(pluginDir, pluginMD5FileName, version); err != nil {
				return err
			}
			log.Successf("[%s] download succeeded.", pluginMD5FileName)
		}

		// check if the plugin matches with .md5
		isMD5Match, err := ifPluginAndMD5Match(pluginDir, pluginFileName, pluginMD5FileName)
		if err != nil {
			return err
		}

		if !isMD5Match {
			// if existing .so doesn't match with .md5, re-download
			log.Infof("Plugin: [%s] doesn't match with .md5 and will be downloaded.", pluginFileName)
			if err = dc.reDownload(pluginDir, pluginFileName, pluginMD5FileName, version); err != nil {
				return err
			}
			// check if the downloaded plugin md5 matches with .md5
			if isMD5Match, err = ifPluginAndMD5Match(pluginDir, pluginFileName, pluginMD5FileName); err != nil {
				return err
			}
			if !isMD5Match {
				return fmt.Errorf("plugin %s doesn't match with .md5", tool.Name)
			}
		}

		log.Infof("Initialize [%s] finished.", pluginName)
		log.Separatorf(pluginName)
	}
	return nil
}

// CheckLocalPlugins checks if the local plugins exists, and matches with md5 value.
func CheckLocalPlugins(tools []configmanager.Tool) error {
	pluginDir := viper.GetString("plugin-dir")
	if pluginDir == "" {
		return fmt.Errorf(`plugins directory should not be ""`)
	}

	log.Infof("Using dir <%s> to store plugins.", pluginDir)

	for _, tool := range tools {
		pluginFileName := tool.GetPluginFileName()
		pluginMD5FileName := tool.GetPluginMD5FileName()
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
		matched, err := ifPluginAndMD5Match(pluginDir, pluginFileName, pluginMD5FileName)
		if err != nil {
			return err
		}
		if !matched {
			return fmt.Errorf("plugin %s doesn't match with .md5", tool.Name)
		}
	}
	return nil
}

// pluginAndMD5Matches checks if the plugins match with .md5
// it returns true if the plugin matches with .md5
// it returns error if the plugin doesn't exist or .md5 doesn't exist
func ifPluginAndMD5Match(pluginDir, soFileName, md5FileName string) (bool, error) {
	soFilePath := filepath.Join(pluginDir, soFileName)
	md5FilePath := filepath.Join(pluginDir, md5FileName)
	isMD5Match, err := md5.FileMatchesMD5(soFilePath, md5FilePath)
	if err != nil {
		return false, err
	}
	return isMD5Match, nil
}
