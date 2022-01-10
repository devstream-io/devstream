package pluginmanager

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	os.Remove(filepath.Join(".", "argocdapp_0.0.1-rc1.so"))

	c := NewDownloadClient()
	err := c.download(".", "argocdapp_0.0.1-rc1.so", "0.0.1-rc1")
	if err != nil {
		t.Fatal("downloaded error")
	}

	os.Remove(filepath.Join(".", "argocdapp_0.0.1-rc1.so"))
}

func TestDownloadNotFound(t *testing.T) {
	c := NewDownloadClient()
	err := c.download(".", "doesntexist", "0.0.1")
	assert.Contains(t, err.Error(), "404")
}
