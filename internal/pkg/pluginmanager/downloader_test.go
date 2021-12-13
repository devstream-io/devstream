package pluginmanager

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDownload(t *testing.T) {
	c := NewDownloadClient()

	os.Remove(filepath.Join(".", "test.so"))

	err := c.download(".", "argocdapp", "0.0.1")
	if err != nil {
		t.Fatal("downloaded error")
	}

	os.Remove(filepath.Join(".", "argocdapp"))
}

func TestDownloadNotFound(t *testing.T) {
	c := NewDownloadClient()
	err := c.download(".", "doesntexist", "0.0.1")
	assert.Contains(t, err.Error(), "404")
}
