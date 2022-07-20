package downloader

import (
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
	return DownloadToFile(url, f)
}

func DownloadToFile(url string, f *os.File) (int64, error) {
	log.Debugf("URL: %s.", url)
	log.Debugf("Filename: %s.", f.Name())

	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	return io.Copy(f, resp.Body)
}
