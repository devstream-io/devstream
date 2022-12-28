package md5

import (
	"os"

	"k8s.io/utils/strings"
)

const md5Length = 32

// FileMatchesMD5 checks current PlugIn MD5 matches with .md5 file
func FileMatchesMD5(fileName, md5FileName string) (bool, error) {
	currentPlugInMD5, err := CalcFileMD5(fileName)
	if err != nil {
		return false, err
	}

	md5ContentBytes, err := os.ReadFile(md5FileName)
	if err != nil {
		return false, err
	}
	// intercept string, md5 code length is 32
	md5Content := strings.ShortenString(string(md5ContentBytes), md5Length)

	return currentPlugInMD5 == md5Content, nil
}

func FilesMD5Equal(file1, file2 string) (bool, error) {
	file1MD5, err := CalcFileMD5(file1)
	if err != nil {
		return false, err
	}

	file2MD5, err := CalcFileMD5(file2)
	if err != nil {
		return false, err
	}

	return file1MD5 == file2MD5, nil
}
