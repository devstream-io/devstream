package version

import (
	"io/ioutil"

	"k8s.io/utils/strings"
)

const md5_length = 32

// ValidateFileMatchMD5 check current PlugIn MD5 matches with .md5 file
func ValidateFileMatchMD5(fileName, md5FileName string) (bool, error) {
	currentPlugInMD5, err := CalcFileMD5(fileName)
	if err != nil {
		return false, err
	}

	md5ContentBytes, err := ioutil.ReadFile(md5FileName)
	if err != nil {
		return false, err
	}
	// intercept string, md5 code is from 0 to 31
	md5Content := strings.ShortenString(string(md5ContentBytes), md5_length)

	return currentPlugInMD5 == md5Content, nil
}
