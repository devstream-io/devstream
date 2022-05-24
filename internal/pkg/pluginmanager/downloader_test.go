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
	err := c.download(".", "argocdapp_0.0.1-rc1.so", "0.0.1-ut-do-not-delete")
	if err != nil {
		t.Fatal("downloaded error")
	}

	os.Remove(filepath.Join(".", "argocdapp_0.0.1-rc1.so"))
}

func TestDownloadNotFound(t *testing.T) {
	c := NewDownloadClient()
	err := c.download(".", "doesntexist", "0.0.1-ut-do-not-delete")
	// Since the right granted to public users on aws does not include listing bucket
	// AWS returns 403 instead of 404 when acquiring an object where bucket does not exist: there is no list right.
	assert.Contains(t, err.Error(), "403")
}
