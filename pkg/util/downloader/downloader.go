package downloader

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

type Downloader struct {
	EnableProgressBar bool
	client            *http.Client
}

func New() *Downloader {
	return &Downloader{
		client: &http.Client{},
	}
}

func (d *Downloader) WithProgressBar() *Downloader {
	d.EnableProgressBar = true
	return d
}

func (d *Downloader) WithClient(client *http.Client) *Downloader {
	d.client = client
	return d
}

// Download a file from the URL to the target path
// if filename is "", use the remote filename at local.
func (d *Downloader) Download(url, filename, targetDir string) (size int64, err error) {
	// handle filename and target dir
	filename, err = parseFilenameAndCreateTargetDir(url, filename, targetDir)
	if err != nil {
		return 0, err
	}

	// get http response
	resp, err := d.getHttpResponse(url)
	if err != nil {
		return 0, err
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Errorf("Close response body failed: %s", err)
		}
	}(resp.Body)

	pluginTmpLocation := filepath.Join(targetDir, filename+".tmp")
	pluginLocation := filepath.Join(targetDir, filename)

	// download to tmp file
	log.Infof("Downloading: [%s] ...", filename)
	downFile, err := os.Create(pluginTmpLocation)
	if err != nil {
		return 0, err
	}
	defer func() {
		err := downFile.Close()
		if err != nil {
			log.Debugf("download create file failed: %s", err)
		}
		err = removeFileIfExists(pluginTmpLocation)
		if err != nil {
			log.Debugf("download create file failed: %s", err)
		}
	}()

	if d.EnableProgressBar {
		// create progress bar when reading response body
		size, err = SetUpProgressBar(resp, downFile)
	} else {
		// just copy response body to file
		size, err = io.Copy(downFile, resp.Body)
	}

	if err != nil {
		return 0, err
	}

	// rename, tmp file to real file
	if err = os.Rename(pluginTmpLocation, pluginLocation); err != nil {
		return 0, err
	}
	return size, nil
}

func parseFilenameAndCreateTargetDir(url, filename, targetDir string) (finalFilename string, err error) {
	log.Debugf("Downloading url is: %s.", url)
	log.Debugf("Target dir: %s.", targetDir)
	if url == "" {
		return "", fmt.Errorf("url must not be empty: %s", url)
	}
	// get filename from local or remote
	if filename != "" {
		// use local filename
		// when filename is "/", ".", "..", the filepath.Base will return "/", ".", ".."
		finalFilename = filepath.Base(filename)
		if finalFilename == "/" || finalFilename == "." || finalFilename == ".." {
			return "", fmt.Errorf("filename must not be dir")
		}
	} else {
		// use remote filename
		// when url is empty filepath.Base(url) will return "."
		finalFilename = filepath.Base(url)
		if finalFilename == "." {
			return "", fmt.Errorf("failed to get the filename from url: %s", url)
		}
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", err
	}

	return finalFilename, nil
}

func (d *Downloader) getHttpResponse(url string) (*http.Response, error) {
	resp, err := d.client.Get(url)

	// check response error
	if err != nil {
		log.Debugf("Download from url failed: %s", err)
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Download file from url failed: %+v", resp)
	}

	return resp, nil
}

// setUpProgressBar create bar and setup
func SetUpProgressBar(resp *http.Response, downFile io.Writer) (int64, error) {
	//get size
	i, _ := strconv.Atoi(resp.Header.Get("Content-Length"))
	sourceSiz := int64(i)
	source := resp.Body

	//create a bar and set param
	bar := pb.New(int(sourceSiz)).SetRefreshRate(time.Millisecond * 10).SetUnits(pb.U_BYTES).SetWidth(100)
	bar.ShowSpeed = true
	bar.ShowTimeLeft = true
	bar.ShowFinalTime = true
	bar.SetWidth(80)
	bar.Start()
	defer bar.Finish()
	writer := io.MultiWriter(downFile, bar)
	return io.Copy(writer, source)
}

func FetchContentFromURL(url string) ([]byte, error) {
	resp, err := http.Get(url)

	// check response error
	if err != nil {
		log.Debugf("Download from url failed: %s", err)
		return nil, err
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Errorf("Close response body failed: %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Download file from url failed: %+v", resp)
	}

	return io.ReadAll(resp.Body)
}

func removeFileIfExists(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filename)
}
