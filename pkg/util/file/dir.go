package file

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/devstream-io/devstream/pkg/util/log"
)

// DirFileFilterFunc is used to filter files when walk directory
// if this func return false, this file's content will not return
type DirFileFilterFunc func(filePath string, isDir bool) bool

// DirFileContentFunc is used to get file content then return this content
type DirFileContentFunc func(filePath string) ([]byte, error)

// DirFileNameFunc is used to make current filename to become map key
// srcPath is the basePath of file, for sometimes you may want to get relativePath of filePath
type DirFileNameFunc func(filePath, srcPath string) string

// GetFileMapByWalkDir will walk in directory return contentMap
func GetFileMapByWalkDir(
	srcPath string, filterFunc DirFileFilterFunc, fileNameFunc DirFileNameFunc, processFunc DirFileContentFunc,
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
		// if process file failed, return error
		content, err := processFunc(path)
		if err != nil {
			return err
		}
		contentMap[dstFileName] = content
		return nil
	}); err != nil {
		log.Debugf("Walk Dir %s failed: %+v", srcPath, err)
		return nil, err
	}
	return contentMap, nil
}

// GetFileMap return map of fileName and content
// if srcPath is a directory, it will invoke GetFileMapByWalkDir to get content map
// if srcPath is a file, it will use fileNameFunc and fileContentFunc to create a map
func GetFileMap(
	srcPath string, filterFunc DirFileFilterFunc, fileNameFunc DirFileNameFunc, fileContentFunc DirFileContentFunc,
) (map[string][]byte, error) {
	pathInfo, err := os.Stat(srcPath)
	if err != nil {
		log.Debugf("dir: get path info failed: %+v", err)
		return nil, err
	}
	if pathInfo.IsDir() {
		return GetFileMapByWalkDir(
			srcPath, filterFunc,
			fileNameFunc, fileContentFunc,
		)
	}
	content, err := fileContentFunc(srcPath)
	if err != nil {
		log.Debugf("dir: process file content failed: %+v", err)
		return nil, err
	}
	fileName := fileNameFunc(srcPath, filepath.Dir(srcPath))
	return map[string][]byte{
		fileName: []byte(content),
	}, nil
}

func CreateTempDir(dirPattern string) (string, error) {
	tempDir, err := os.MkdirTemp("", dirPattern)
	if err != nil {
		log.Debugf("create tempDir %s failed: %+v", dirPattern, err)
		return "", err
	}
	return tempDir, err
}

// DirFileFilterDefaultFunc is used for GetFileMap
// it will return false if isDir is true
func DirFileFilterDefaultFunc(filePath string, isDir bool) bool {
	return !isDir
}
