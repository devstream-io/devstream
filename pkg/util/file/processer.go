package file

import (
	"os"
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/zip"
)

// unZipFileProcesser will unzip zip file and return unzip files dir path
func unZipFileProcesser(zipFilePath string) (string, error) {
	// 1. create tempDir to save unzip files
	dirName := filepath.Dir(zipFilePath)
	tempDirName, err := os.MkdirTemp(dirName, defaultTempName)
	if err != nil {
		return "", err
	}
	err = zip.UnZip(zipFilePath, tempDirName)
	if err != nil {
		return "", err
	}
	return tempDirName, nil
}
