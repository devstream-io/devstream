// 1. init download
// 2. get assets *.so

package pluginmanager

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"

	"github.com/devstream-io/devstream/pkg/util/log"
)

const (
	defaultRetryCount = 3
	defaultReleaseUrl = "https://download.devstream.io"
)

type plugDownloader func(reqClient *resty.Client, url, plugName string) error

type DownloadClient struct {
	*resty.Client
	pluginGetter plugDownloader
}

func downloadPlugin(reqClient *resty.Client, url, plugName string) error {
	response, err := reqClient.R().
		SetOutput(plugName).
		SetHeader("Accept", "application/octet-stream").
		Get(url)
	if err != nil {
		return err
	}
	if response.StatusCode() != http.StatusOK {
		if err = os.Remove(filepath.Join(url, plugName)); err != nil {
			return err
		}
		err = fmt.Errorf("downloading plugin %s from %s status code %d", plugName, url, response.StatusCode())
		log.Error(err)
		return err
	}
	return nil
}

func NewDownloadClient() *DownloadClient {
	dClient := DownloadClient{
		pluginGetter: downloadPlugin,
	}
	dClient.Client = resty.New()
	dClient.SetRetryCount(defaultRetryCount)
	return &dClient
}

func (dc *DownloadClient) download(pluginDir, pluginFilename, version string) error {
	dc.SetOutputDirectory(pluginDir)

	// download plug file
	downloadURL := fmt.Sprintf("%s/v%s/%s", defaultReleaseUrl, version, pluginFilename)
	tmpName := pluginFilename + ".tmp"
	err := dc.pluginGetter(dc.Client, downloadURL, tmpName)
	if err != nil {
		return err
	}
	// rename, tmp file to real file
	err = os.Rename(
		filepath.Join(pluginDir, tmpName),
		filepath.Join(pluginDir, pluginFilename))
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
