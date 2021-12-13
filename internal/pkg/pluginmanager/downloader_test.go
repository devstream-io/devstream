package pluginmanager

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetAssetWithRetry(t *testing.T) {
	c := NewDownloadClient()

	os.Remove(filepath.Join(".", "test.so"))

	err := c.download(".", "argocdapp", "0.0.1")
	if err != nil {
		t.Fatal("downloaded error")
	}

	os.Remove(filepath.Join(".", "argocdapp"))
}
