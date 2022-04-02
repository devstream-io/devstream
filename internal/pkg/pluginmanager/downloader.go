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
	defaultReleaseUrl = "https://github.com/devstream-io/devstream/releases/download"
)

type DownloadClient struct {
	*resty.Client
}

func NewDownloadClient() *DownloadClient {
	dClient := DownloadClient{}
	dClient.Client = resty.New()
	dClient.SetRetryCount(defaultRetryCount)
	return &dClient
}

func (dc *DownloadClient) download(pluginDir, pluginFilename, version string) error {
	dc.SetOutputDirectory(pluginDir)

	downloadURL := fmt.Sprintf("%s/v%s/%s", defaultReleaseUrl, version, pluginFilename)
	tmpName := pluginFilename + ".tmp"

	response, err := dc.R().
		SetOutput(tmpName).
		SetHeader("Accept", "application/octet-stream").
		Get(downloadURL)
	if err != nil {
		return err
	}
	if response.StatusCode() != http.StatusOK {
		if err = os.Remove(filepath.Join(pluginDir, tmpName)); err != nil {
			return err
		}
		err = fmt.Errorf("downloading plugin %s from %s status code %d", pluginFilename, downloadURL, response.StatusCode())
		log.Error(err)
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
