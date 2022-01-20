package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func Download(url, targetDir string) (int64, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}

	fileName := filepath.Base(url)
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
