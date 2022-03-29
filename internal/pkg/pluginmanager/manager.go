package pluginmanager

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"

	"github.com/merico-dev/stream/internal/pkg/configloader"
	"github.com/merico-dev/stream/internal/pkg/version"
	"github.com/merico-dev/stream/pkg/util/log"
	"github.com/merico-dev/stream/pkg/util/md5"
)

const DTMRemoteMD5Dir = ".remote"

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
		pluginMD5FileName := configloader.GetPluginMD5FileName(&tool)
		// plugin does not exist
		if _, err := os.Stat(filepath.Join(pluginDir, pluginFileName)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// download .so file
				if err := dc.download(pluginDir, pluginFileName, version.Version); err != nil {
					return err
				}
				log.Successf("[%s] download succeeded.", pluginFileName)
			}
			return err
		}
		// .md5 does not exist
		if _, err := os.Stat(filepath.Join(pluginDir, pluginMD5FileName)); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				// download .md5 file
				if err := dc.download(pluginDir, pluginMD5FileName, version.Version); err != nil {
					return err
				}
				log.Successf("[%s] download succeeded.", pluginMD5FileName)
			}
			return err
		}
		// check if the plugin matches with .md5
		isMD5Match, err := md5.FileMatchesMD5(filepath.Join(pluginDir, pluginFileName), filepath.Join(pluginDir, pluginMD5FileName))
		if err != nil {
			return err
		}
		// if .so matches with .md5, continue
		if isMD5Match {
			log.Infof("Plugin: %s already exists, no need to download.", pluginFileName)
			continue
		}
		// if existing .so doesn't matches with .md5, re-download
		log.Infof("Plugin: %s doesn't match with .md5 and will be downloaded.", pluginFileName)
		if err := redownloadPlugins(dc, pluginDir, pluginFileName, pluginMD5FileName, version.Version); err != nil {
			return err
		}

		// check if the downloaded plugin md5 matches with .md5
		if err := pluginAndMD5Matches(pluginDir, pluginFileName, pluginMD5FileName, tool.Name); err != nil {
			return err
		}
	}
	return nil
}

// CheckLocalPlugins checks if the local plugins match with .md5
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

// DownloadDTMsMD5 download remote dtm .md5 for compare with local .md5
func DownloadDTMsMD5(remoteMD5Dir, dtmFileName string) error {
	dc := NewPbDownloadClient()
	if err := dc.download(remoteMD5Dir, dtmFileName, version.Version); err != nil {
		return err
	}
	return nil
}

// CompareDtmMD5 compare dtm itself and remote md5, if does not match, cannot  download plugins from remote
func CompareDtmMD5() error {
	// get remote dtm .md5
	dtmMD5FileName := configloader.GetDtmMD5FileName()
	if err := DownloadDTMsMD5(DTMRemoteMD5Dir, dtmMD5FileName); err != nil {
		return err
	}

	selfDtmName, err := getDtmSelf()
	if err != nil {
		return err
	}

	isMD5Match, err := md5.FileMatchesMD5(selfDtmName, filepath.Join(DTMRemoteMD5Dir, dtmMD5FileName))
	if err != nil {
		return err
	}
	if !isMD5Match {
		return fmt.Errorf("DTM %s doesn't match with remote .md5, cannot download Plugins from remote release version", selfDtmName)
	}
	log.Info("DTM md5 check passed, continue...")
	return nil
}

// getDtmSelf get executing dtm itself
func getDtmSelf() (string, error) {
	path, err := exec.LookPath(os.Args[0])
	if err != nil {
		return "", err
	}
	return path, nil
}
