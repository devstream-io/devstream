// 1. get plugins *.so
// 2. show progress bar on console

package pluginmanager

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/log"
)

type PluginDownloadClient struct {
	*http.Client
	baseURL string
}

func NewPluginDownloadClient(baseURL string) *PluginDownloadClient {
	dClient := PluginDownloadClient{}
	dClient.Client = http.DefaultClient
	dClient.Client.Timeout = time.Second * 60 * 60
	dClient.baseURL = baseURL
	return &dClient
}

// download from release assets
func (pd *PluginDownloadClient) download(pluginsDir, pluginOrMD5Filename, version string) error {
	downloadURL := fmt.Sprintf("%s/v%s/%s", pd.baseURL, version, pluginOrMD5Filename)
	dc := downloader.New().WithProgressBar().WithClient(pd.Client)
	_, err := dc.Download(downloadURL, pluginOrMD5Filename, pluginsDir)
	return err
}

// reDownload plugins from remote
func (pd *PluginDownloadClient) reDownload(pluginDir, pluginFileName, pluginMD5FileName, version string) error {
	if err := os.Remove(filepath.Join(pluginDir, pluginFileName)); err != nil {
		return err
	}
	if err := os.Remove(filepath.Join(pluginDir, pluginMD5FileName)); err != nil {
		return err
	}
	// download .so file
	if err := pd.download(pluginDir, pluginFileName, version); err != nil {
		return err
	}
	log.Successf("[%s] download succeeded.", pluginFileName)
	// download .md5 file
	if err := pd.download(pluginDir, pluginMD5FileName, version); err != nil {
		return err
	}
	log.Successf("[%s] download succeeded.", pluginMD5FileName)
	return nil
}
