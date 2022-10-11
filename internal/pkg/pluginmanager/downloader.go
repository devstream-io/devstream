// 1. get plugins *.so
// 2. show progress bar on console

package pluginmanager

import (
	"fmt"
	"net/http"
	"time"

	"github.com/devstream-io/devstream/pkg/util/downloader"
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
func (pd *PluginDownloadClient) download(pluginsDir, pluginFilename, version string) error {
	downloadURL := fmt.Sprintf("%s/v%s/%s", pd.baseURL, version, pluginFilename)
	dc := downloader.New().WithProgressBar().WithClient(pd.Client)
	_, err := dc.Download(downloadURL, pluginFilename, pluginsDir)
	return err
}
