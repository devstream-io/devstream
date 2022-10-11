package file

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/devstream-io/devstream/pkg/util/log"
	"github.com/devstream-io/devstream/pkg/util/zip"
)

const dirTempName = "dir_temp"

type DirFIleFilterFunc func(filePath string, isDir bool) bool
type DirFileProcessFunc func(filePath string) ([]byte, error)
type DirFileNameFunc func(filePath, srcPath string) string

func WalkDir(
	srcPath string, filterFunc DirFIleFilterFunc, fileNameFunc DirFileNameFunc, processFunc DirFileProcessFunc,
) (map[string][]byte, error) {
	contentMap := make(map[string][]byte)
	if err := filepath.Walk(srcPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			log.Debugf("Walk error: %s.", err)
			return err
		}
		if !filterFunc(path, info.IsDir()) {
			log.Debugf("Walk file is filtered: %s.", path)
			return nil
		}

		// if file ends-with tpl, render this file, else copy this file directly
		dstFileName := fileNameFunc(path, srcPath)
		content, err := processFunc(path)
		if err != nil {
			return nil
		}
		contentMap[dstFileName] = content
		return nil
	}); err != nil {
		return nil, err
	}
	return contentMap, nil
}

// DownloadAndUnzipAndRenderGitHubRepo will download/unzip/render git repo files and return unzip files dir path

func DownloadAndUnzipFile(url string) (string, error) {
	zipFileContent, err := downloader.FetchContentFromURL(url)
	if err != nil {
		log.Debugf("unzip download file copy content failed: %s", err)
		return "", err
	}
	zipFile, err := os.CreateTemp("", dirTempName)
	if err != nil {
		log.Debugf("unzip create temp file error: %s", err)
		return "", err
	}
	defer zipFile.Close()
	_, err = io.Copy(zipFile, bytes.NewBuffer(zipFileContent))
	if err != nil {
		log.Debugf("unzip copy file error: %s", err)
		return "", err
	}

	unZipDir, err := Unzip(zipFile.Name())
	if err != nil {
		return "", err
	}
	return unZipDir, nil
}

// Unzip will unzip zip file and return unzip files dir path
func Unzip(zipFilePath string) (string, error) {
	// 1. create tempDir to save unzip files
	dirName := filepath.Dir(zipFilePath)
	tempDirName, err := os.MkdirTemp(dirName, dirTempName)
	if err != nil {
		return "", err
	}
	err = zip.UnZip(zipFilePath, tempDirName)
	if err != nil {
		return "", err
	}
	return tempDirName, nil
}
