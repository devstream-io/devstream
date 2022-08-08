package file

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// getFileFromURL will download file from url and return file path
func getFileFromURL(url string) (string, error) {
	// 1. create temp file for save content
	tempFile, err := os.CreateTemp("", defaultTempName)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// 2. download content to file
	log.Debugf("Download URL: %s to filename %s", url, tempFile.Name())
	resp, err := http.Get(url)

	// 3. check response error
	if err != nil {
		log.Debugf("Download file from url failed: %s", err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Download file from url failed: %+v", resp)
	}

	// 4. copy response body to template file
	_, err = io.Copy(tempFile, resp.Body)
	if err != nil {
		log.Debugf("Download file copy content failed: %s", err)
		return "", err
	}
	return tempFile.Name(), nil
}
