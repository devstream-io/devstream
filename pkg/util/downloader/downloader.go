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
	log.Debugf("URL: %s.", url)
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	var fileName = filename
	if fileName == "" {
		fileName = filepath.Base(url)
	}
	log.Debugf("Filename: %s.", fileName)
	if fileName == "." {
		return 0, fmt.Errorf("failed to get the filename from url: %s", url)
	}

	err = os.MkdirAll(targetDir, 0755)
	if err != nil {
		return 0, err
	}

	f, err := os.Create(filepath.Join(targetDir, fileName))
	if err != nil {
		return 0, err
	}

	return io.Copy(f, resp.Body)
}
