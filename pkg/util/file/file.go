package file

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/log"
)

func CopyFile(srcFile, dstFile string) (err error) {
	// prepare source file
	sFile, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := sFile.Close(); closeErr != nil {
			log.Errorf("Failed to close file %s: %s", srcFile, closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}()

	// create destination file
	dFile, err := os.Create(dstFile)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := dFile.Close(); closeErr != nil {
			log.Errorf("Failed to close file %s: %s", dstFile, closeErr)
			if err == nil {
				err = closeErr
			}
		}
	}()

	// copy and sync
	if _, err = io.Copy(dFile, sFile); err != nil {
		return nil
	}
	return dFile.Sync()
}

// GenerateAbsFilePath return all the path with a given file name
func GenerateAbsFilePath(baseDir, file string) (string, error) {
	file = filepath.Join(baseDir, file)

	fileExist := func(path string) bool {
		if _, err := os.Stat(file); err != nil {
			log.Errorf("File %s not exists. Error: %s", file, err)
			return false
		}
		return true
	}

	absFilePath, err := filepath.Abs(file)
	if err != nil {
		log.Errorf(`Failed to get absolute path fo "%s".`, file)
		return "", err
	}
	log.Debugf("Abs path is %s.", absFilePath)
	if fileExist(absFilePath) {
		return absFilePath, nil
	} else {
		return "", fmt.Errorf("file %s not exists", absFilePath)
	}
}

// GetFileAbsDirPath will return abs dir path for file
func GetFileAbsDirPath(fileLoc string) (string, error) {
	fileAbs, err := filepath.Abs(fileLoc)
	if err != nil {
		return "", fmt.Errorf("%s not exists", fileLoc)
	}
	return filepath.Dir(fileAbs), nil
}

// GetFileAbsDirPathOrDirItself will return abs dir path,
// for file: return its parent directory
// for directory: return dir itself
func GetFileAbsDirPathOrDirItself(fileLoc string) (string, error) {
	fileAbs, err := filepath.Abs(fileLoc)
	if err != nil {
		return "", fmt.Errorf("%s not exists", fileLoc)
	}
	file, err := os.Stat(fileAbs)
	if err != nil {
		return "", fmt.Errorf("%s not exists", fileLoc)
	}
	// if it is a directory, return itself
	if file.IsDir() {
		return fileAbs, nil
	}
	// if it is a file, return its parent directory
	return filepath.Dir(fileAbs), nil
}
