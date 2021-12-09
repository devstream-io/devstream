package download

import (
	"os"
	"testing"
)

//test download .so file
func TestGetAssetswithretry(t *testing.T) {

	downloader := NewDownloadClient()
	downloader.AssetName = "argocdapp"
	downloader.Version = ""
	downloader.Filepath = "plugins/test_1.so"

	os.Remove(downloader.Filepath)
	downloader.GetAssetswithretry()

	if FileExist(downloader.Filepath) {
		t.Fatal("no file downloaded")
	}
	os.Remove(downloader.Filepath)
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}
