package file

import (
	"bytes"
	"io"
	"os"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// getFileFromContent will create a temp file based on content
func getFileFromContent(content string) (string, error) {
	// 1. create temp file for save content
	tempFile, err := os.CreateTemp("", defaultTempName)
	if err != nil {
		return "", err
	}
	defer tempFile.Close()

	// 2. save content to file
	_, err = io.Copy(tempFile, bytes.NewBufferString(content))
	if err != nil {
		log.Debugf("Download file copy content failed: %s", err)
		return "", err
	}
	return tempFile.Name(), nil
}
