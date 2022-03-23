package md5

import (
	"io/ioutil"

	"k8s.io/utils/strings"
)

const md5_length = 32

// FileMatchesMD5 checks current PlugIn MD5 matches with .md5 file
func FileMatchesMD5(fileName, md5FileName string) (bool, error) {
	currentPlugInMD5, err := CalcFileMD5(fileName)
	if err != nil {
		return false, err
	}

	md5ContentBytes, err := ioutil.ReadFile(md5FileName)
	if err != nil {
		return false, err
	}
	// intercept string, md5 code length is 32
	md5Content := strings.ShortenString(string(md5ContentBytes), md5_length)

	return currentPlugInMD5 == md5Content, nil
}
