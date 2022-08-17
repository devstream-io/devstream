package downloader

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

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
	content, err := FetchContentFromURL(url)
	if err != nil {
		return 0, err
	}
	return io.Copy(f, bytes.NewBuffer(content))
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
