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

// if filename is "", use the remote filename at local.
func Download(url, filename, targetDir string) (int64, error) {
	log.Debugf("Target dir: %s.", targetDir)
	if url == "" {
		return 0, fmt.Errorf("url must not be empty: %s", url)
	}
	if filename == "" {
		// when url is empty filepath.Base(url) will return "."
		filename = filepath.Base(url)
	}
	if filename == "." {
		return 0, fmt.Errorf("failed to get the filename from url: %s", url)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return 0, err
	}

	f, err := os.Create(filepath.Join(targetDir, filename))
	if err != nil {
		return 0, err
	}
	defer func() {
		err := f.Close()
		if err != nil {
			log.Debugf("download create file failed: %s", err)
		}
	}()
	return FetchContentFromURLAndSetUpProgressBar(url, f)
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

func FetchContentFromURLAndSetUpProgressBar(url string, dst io.Writer) (int64, error) {
	resp, err := http.Get(url)

	// check response error
	if err != nil {
		log.Debugf("Download from url failed: %s", err)
		return 0, err
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			log.Errorf("Close response body failed: %s", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Download file from url failed: %+v", resp)
	}
	return SetUpProgressBar(resp, dst)
}
