package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// if filename is "", use the remote filename at local.
func Download(url, filename, targetDir string) (int64, error) {
	if url == "" {
		return 0, fmt.Errorf("url must not be empty: %s", url)
	}
	if filename == "." {
		return 0, fmt.Errorf("filename must not be [%s]", filename)
	}
	if filename == "" {
		// when url is empty filepath.Base(url) will return "."
		filename = filepath.Base(url)
	}

	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return 0, err
	}

	f, err := os.Create(filepath.Join(targetDir, filename))
	if err != nil {
		return 0, err
	}
	defer f.Close()
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	return io.Copy(f, resp.Body)
}
