package version

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidatePlugInMD5Exist file exists,  md5 info is also in dtm MD5string
func TestValidatePlugInMD5Exist(t *testing.T) {
	clearMD5string()
	fileName1 := "test1.so"
	fileName2 := "test2.so"

	md5_1, md5_2, err := createFile(fileName1, fileName2)
	assert.NoError(t, err)

	setMD5string(md5_1 + ":" + md5_2)
	exist, err := ValidatePlugInMD5(fileName1)
	assert.NoError(t, err)
	assert.True(t, exist)
	exist, err = ValidatePlugInMD5(fileName2)
	assert.NoError(t, err)
	assert.True(t, exist)

}

// TestValidatePlugInMD5NotExist file exists, but md5 info is not in dtm MD5string
func TestValidatePlugInMD5NotExist(t *testing.T) {
	clearMD5string()
	fileName1 := "test1_not_in_md5string.so"
	fileName2 := "test2_not_in_md5string.so"

	_, _, err := createFile(fileName1, fileName2)
	assert.NoError(t, err)

	exist, err := ValidatePlugInMD5(fileName1)
	assert.NoError(t, err)
	assert.False(t, exist)
	exist, err = ValidatePlugInMD5(fileName2)
	assert.NoError(t, err)
	assert.False(t, exist)
}

func createFile(fileName1, fileName2 string) (string, string, error) {
	f1, err := os.Create(fileName1)
	if err != nil {
		return "", "", err
	}
	f2, err := os.Create(fileName2)
	if err != nil {
		return "", "", err
	}
	defer f1.Close()
	defer f2.Close()

	md5_1, err := CalcFileMD5(fileName1)
	if err != nil {
		return "", "", err
	}
	md5_2, err := CalcFileMD5(fileName2)
	if err != nil {
		return "", "", err
	}
	return md5_1, md5_2, nil
}

func setMD5string(md5 string) {
	MD5String = md5
}

func clearMD5string() {
	MD5String = ""
}
