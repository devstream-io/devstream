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

type PbDownloadClient struct {
	*http.Client
	baseURL string
}

func NewPbDownloadClient(baseURL string) *PbDownloadClient {
	dClient := PbDownloadClient{}
	dClient.Client = http.DefaultClient
	dClient.Client.Timeout = time.Second * 60 * 60
	dClient.baseURL = baseURL
	return &dClient
}

// download from release assets
func (pd *PbDownloadClient) download(pluginsDir, pluginFilename, version string) error {

	err := createPathIfNotExists(pluginsDir)
	if err != nil {
		return err
	}

	downloadURL := fmt.Sprintf("%s/v%s/%s", pd.baseURL, version, pluginFilename)
	log.Debugf("Downloading url is: %s.", downloadURL)

	tmpName := pluginFilename + ".tmp"

	resp, err := pd.Get(downloadURL)
	if err != nil {
		return err
	}
	pluginTmpLocation := filepath.Join(pluginsDir, tmpName)
	pluginLocation := filepath.Join(pluginsDir, pluginFilename)

	if resp.StatusCode == http.StatusOK {
		log.Infof("Downloading: [%s] ...", pluginFilename)

		downFile, err := os.Create(pluginTmpLocation)
		if err != nil {
			return err
		}

		defer downFile.Close()
		// create progress bar when reading response body
		_, errSetup := downloader.SetUpProgressBar(resp, downFile)
		if errSetup != nil {
			log.Error(errSetup)
			return errSetup
		}
	} else {
		log.Errorf("[%s] download failed, %s.", pluginFilename, resp.Status)
		if err = os.Remove(pluginTmpLocation); err != nil {
			log.Errorf("Remove [%s] failed, %s.", pluginLocation, err)
		}
		return fmt.Errorf("downloading %s from %s status code %d", pluginFilename, downloadURL, resp.StatusCode)
	}

	// rename, tmp file to real file
	if err = os.Rename(pluginTmpLocation, pluginLocation); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func createPathIfNotExists(path string) error {
	_, err := os.Stat(path)
	if !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		return err
	}
	return nil
}
