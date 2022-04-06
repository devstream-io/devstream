package md5

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidatePlugInMD5Exist file exists,  md5 info matches with .md5
func TestValidateFileMatchMD5(t *testing.T) {
	fileName := "test1"

	md5FileName, err := createFileAndMD5File(fileName)
	assert.NoError(t, err)

	exist, err := FileMatchesMD5(fileName, md5FileName)
	assert.NoError(t, err)
	assert.True(t, exist)

	os.Remove(fileName)
	os.Remove(fileName + ".md5")
}

// TestValidatePlugInMD5NotExist file exists, but md5 info does not match with .md5
func TestValidateFileMatchMD5NotExist(t *testing.T) {
	fileNameMisMatch := "test1_not_in_md5string.so"
	fileName := "test1"

	_, err := os.Create(fileNameMisMatch)
	assert.NoError(t, err)

	md5FileName, err := createFileAndMD5File(fileName)
	assert.NoError(t, err)

	exist, err := FileMatchesMD5(fileName, md5FileName)
	assert.NoError(t, err)
	assert.True(t, exist)

	os.Remove(fileName)
	os.Remove(fileName + ".md5")
	os.Remove(fileNameMisMatch)
}

func createFileAndMD5File(fileName string) (string, error) {
	f1, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	defer f1.Close()

	md5, err := CalcFileMD5(fileName)
	if err != nil {
		return "", err
	}
	md5FileName := fmt.Sprint(fileName, ".md5")
	md5File, err := os.Create(md5FileName)
	if err != nil {
		return "", err
	}
	_, err = md5File.Write([]byte(md5))
	if err != nil {
		return "", err
	}
	return md5FileName, nil
}
