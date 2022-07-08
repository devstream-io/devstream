package pluginmanager

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

func mockPlugGetter(reqClient *resty.Client, url, plugName string) error {
	return nil
}

func mockPlugNotFoundGetter(reqClient *resty.Client, url, plugName string) error {
	return errors.New("Plug Not Exist")
}

func TestDownloadSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	plugName := "argocdapp_0.0.1-rc1.so"
	version := "0.0.1-ut-do-not-delete"
	c := NewDownloadClient()
	tmpFilePath := filepath.Join(tmpDir, fmt.Sprintf("%s.tmp", plugName))
	os.Create(tmpFilePath)
	c.pluginGetter = mockPlugGetter
	err := c.download(tmpDir, plugName, version)
	assert.Nil(t, err)
	// check plug file renamed
	_, err = os.Stat(filepath.Join(tmpDir, plugName))
	assert.Nil(t, err)
}

func TestDownloadFileNotDownloadSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	c := NewDownloadClient()
	c.pluginGetter = mockPlugGetter
	err := c.download(tmpDir, "argocdapp_0.0.1-rc1.so", "0.0.1-ut-do-not-delete")
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestDownloadNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	c := NewDownloadClient()
	c.pluginGetter = mockPlugNotFoundGetter
	err := c.download(tmpDir, "doesntexist", "0.0.1-ut-do-not-delete")
	assert.Contains(t, err.Error(), "Plug Not Exist")
}
