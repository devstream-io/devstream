package configloader

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbsCustomPathExists(t *testing.T) {
	customPath := "tools.yaml"
	configFileName := "config.yaml"
	_, err := os.Create(customPath)
	assert.NoError(t, err)
	defer func() {
		err := os.Remove(customPath)
		assert.NoError(t, err)
	}()

	absPath, err := filepath.Abs(customPath)
	assert.NoError(t, err)

	path, err := parseCustomPath(configFileName, absPath)
	assert.NoError(t, err)
	assert.Equal(t, path, absPath)

}

func TestAbsCustomPathNotExists(t *testing.T) {
	customPath := "tools.yaml"
	configFileName := "config.yaml"

	absPath, err := filepath.Abs(customPath)
	assert.NoError(t, err)

	_, err = parseCustomPath(configFileName, absPath)
	assert.EqualError(t, err, "stat "+absPath+": no such file or directory")
}

func TestRelativeCustomPathExists(t *testing.T) {
	customPath := "tools.yaml"
	configFileName := "config.yaml"
	_, err := os.Create(customPath)
	assert.NoError(t, err)
	_, err = os.Create(configFileName)
	assert.NoError(t, err)

	defer func() {
		err := os.Remove(customPath)
		assert.NoError(t, err)
		err = os.Remove(configFileName)
		assert.NoError(t, err)
	}()

	customPathAbsPath, err := parseCustomPath(configFileName, customPath)
	assert.NoError(t, err)

	configFileAbsPath, err := filepath.Abs(configFileName)
	assert.NoError(t, err)
	assert.Equal(t, customPathAbsPath, filepath.Join(filepath.Dir(configFileAbsPath), customPath))
}

func TestRelativeCustomPathNotExists(t *testing.T) {
	customPath := "tools.yaml"
	configFileName := "config.yaml"
	_, err := os.Create(configFileName)
	assert.NoError(t, err)

	defer func() {
		err = os.Remove(configFileName)
		assert.NoError(t, err)
	}()

	configFileAbsPath, err := filepath.Abs(configFileName)
	assert.NoError(t, err)

	_, err = parseCustomPath(configFileName, customPath)
	assert.EqualError(t, err, "stat "+filepath.Join(filepath.Dir(configFileAbsPath), customPath)+": no such file or directory")
}
