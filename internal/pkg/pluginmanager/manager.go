package pluginmanager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
	"github.com/devstream-io/devstream/internal/pkg/version"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/md5"
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

	// record download status, tool.key -> status
	downloadMap := make(map[string]string)

	const (
		unknown = ""
		fail    = "fail"
		success = "success"
	)

	// judge if the download status of all dependencies of given plugin are fail
	ifDependsFiled := func(tool *configloader.Tool) (failed bool, failedDepend string) {
		for _, dep := range tool.TrimmedDependsOn() {
			if downloadMap[dep] == fail {
				return true, dep
			}
		}

		return false, ""
	}

	// judge if the download status of all dependencies of given plugin are specific
	ifDependsSpecified := func(tool *configloader.Tool) bool {
		specific := true
		for _, dep := range tool.TrimmedDependsOn() {
			// if map key not set, it will return default string value ""(unknown)
			if downloadMap[dep] == unknown {
				specific = false
				break
			}
		}
		return specific
	}

	// get tools whose dependencies' download status are all specific (success or fail)
	getSpecificStatusTools := func() []configloader.Tool {
		var tools []configloader.Tool
		for _, tool := range conf.Tools {
			// skip if the status of status is specified
			if downloadMap[tool.Key()] == success || downloadMap[tool.Key()] == fail {
				continue
			}

			// select all plugins whose dependencies' download status are specific
			if ifDependsSpecified(&tool) {
				tools = append(tools, tool)
			}
		}
		for _, t := range tools {
			log.Debugf("next download plugins patch: %v", t.Key())
		}

		return tools
	}

	// download one tool
	downloadOneTool := func(tool *configloader.Tool) error {
		// if the download of the dependent plugin fails, the plugin will not be downloaded either
		if failed, failedDepend := ifDependsFiled(tool); failed {
			log.Warnf("The download of plugin <%s> was cancelled because the download of the plugin <%s> it depends on failed.", tool.Key(), failedDepend)
			// plugins that depend on this plugin will also not be downloaded
			downloadMap[tool.Key()] = fail
			return nil
		}

		pluginFileName := configloader.GetPluginFileName(tool)
		pluginMD5FileName := configloader.GetPluginMD5FileName(tool)
		// plugin does not exist
		if _, err := os.Stat(filepath.Join(pluginDir, pluginFileName)); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
			// download .so file
			if err := dc.download(pluginDir, pluginFileName, version.Version); err != nil {
				return err
			}
			log.Successf("[%s] download succeeded.", pluginFileName)
		}
		// .md5 does not exist
		if _, err := os.Stat(filepath.Join(pluginDir, pluginMD5FileName)); err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return err
			}
			// download .md5 file
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
		// if .so matches with .md5, continue
		if isMD5Match {
			log.Infof("Plugin: %s already exists, no need to download.", pluginFileName)
			downloadMap[tool.Key()] = success
			return nil
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
		downloadMap[tool.Key()] = success
		return nil
	}

	for {
		// traverse until all tools' download status are specific.
		// you can think of this as a simplified version of finding all nodes with entry degree 0 of a directed acyclic graph
		// returning a batch of nodes with 0 entry at a time instead of just one is to support the parallel download feature in the future.
		tools := getSpecificStatusTools()
		if len(tools) == 0 {
			break
		}
		log.Debug("start to download next patches")

		var eg errgroup.Group
		for _, tool := range tools {
			// define a new var to avoid each loop use the same tool
			tool := tool
			// download one tool
			eg.Go(func() error {
				return downloadOneTool(&tool)
			})
		}
		if err := eg.Wait(); err != nil {
			return err
		}
	}

	// check if all plugins are traversed
	traverseCount := 0
	for _, status := range downloadMap {
		if status != unknown {
			traverseCount++
		}
	}
	if traverseCount != len(conf.Tools) {
		return fmt.Errorf("the download of some plugins failed, please check the log for more details")
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
