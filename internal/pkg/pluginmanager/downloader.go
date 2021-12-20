// 1. init download
// 2. get assets *.so

package pluginmanager

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-resty/resty/v2"
)

const (
	defaultReleaseUrl = "https://github.com/merico-dev/stream/releases/download"
)

type DownloadClient struct {
	Filepath   string
	AssetName  string
	Version    string
	ReleaseUrl string
	client     *resty.Client
	RetryCount int
}

func NewDownloadClient() *DownloadClient {
	return &DownloadClient{
		ReleaseUrl: defaultReleaseUrl,
		client:     resty.New(),
		RetryCount: 3,
	}
}

func (dc *DownloadClient) download(pluginsDir, pluginFilename, version string) error {
	dc.client.SetOutputDirectory(pluginsDir)

	downloadURL := fmt.Sprintf("%s/v%s/%s", defaultReleaseUrl, version, pluginFilename)
	tmpName := pluginFilename + ".tmp"

	response, err := dc.client.R().SetOutput(tmpName).
		SetHeader("Accept", "application/octet-stream").
		Get(downloadURL)
	if response.StatusCode() != http.StatusOK {
		os.Remove(filepath.Join(pluginsDir, tmpName))
		errMsg := fmt.Sprintf("Downloading plugin %s from %s status code %d", pluginFilename, downloadURL, response.StatusCode())
		log.Print(errMsg)
		return errors.New(errMsg)
	}
	if err != nil {
		log.Print(err)
		return err
	}

	// rename, tmp file to real file
	err = os.Rename(
		filepath.Join(pluginsDir, tmpName),
		filepath.Join(pluginsDir, pluginFilename))
	if err != nil {
		log.Print(err)
		return err
	}

	return nil
}
