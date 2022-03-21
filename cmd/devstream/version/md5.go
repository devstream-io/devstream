package version

import (
	"strings"
)

const delimiter = ":"

// MD5String is the .
// Assignment by the command:
// -o dtm-${GOOS}-${GOARCH} ./cmd/devstream/`
// See the Makefile for more info.
var MD5String string

// ValidateMD5 validate MD5 array contains current plugIn MD5
func validateMD5(plugInMD5 string, supportedMD5 []string) bool {
	for _, eachItem := range supportedMD5 {
		if eachItem == plugInMD5 {
			return true
		}
	}
	return false
}

// parseMD5String split MD5 string to MD5 array
func parseMD5String(md5String string) []string {
	return strings.Split(md5String, delimiter)
}

// ValidatePlugInMD5 check current PlugIn MD5 if exists
func ValidatePlugInMD5(fileName string) (bool, error) {
	currentPlugInMD5, err := calcFileMD5(fileName)
	if err != nil {
		return false, err
	}

	return validateMD5(currentPlugInMD5, parseMD5String(MD5String)), nil
}
