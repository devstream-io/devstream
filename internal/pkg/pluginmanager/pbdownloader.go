// 1. get plugins *.so
// 2. show progress bar on console

package pluginmanager

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/cheggaaa/pb"

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
func (pd *PbDownloadClient) download(pluginsDir, pluginFilename, ver string) error {
	err := createPathIfNotExists(pluginsDir)
	if err != nil {
		return err
	}

	downloadURL := fmt.Sprintf("%s/v%s/%s", pd.baseURL, ver, pluginFilename)
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
		errSetup := setUpProgressBar(resp, downFile)
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

// setUpProgressBar create bar and setup
func setUpProgressBar(resp *http.Response, downFile *os.File) error {
	// get size
	i, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	sourceSiz := int64(i)
	source := resp.Body

	// create a bar and set param
	bar := pb.New(int(sourceSiz)).SetRefreshRate(time.Millisecond * 10).SetUnits(pb.U_BYTES).SetWidth(100)
	bar.ShowSpeed = true
	bar.ShowTimeLeft = true
	bar.ShowFinalTime = true
	bar.SetWidth(80)
	bar.Start()

	writer := io.MultiWriter(downFile, bar)
	_, err := io.Copy(writer, source)
	if err != nil {
		log.Error(err)
		return err
	}
	bar.Finish()
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
