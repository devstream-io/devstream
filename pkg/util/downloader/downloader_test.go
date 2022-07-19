package downloader_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/downloader"
)

func CreateFile(dir, filename string) (*os.File, error) {
	return os.Create(filepath.Join(dir, filename))
}

var _ = Describe("Downloader", func() {
	Context("Downloader test", func() {
		var url = "https://github.com/devstream-io/devstream/releases/download/v0.0.1/argocdapp_0.0.1.so"
		var error_url = "://github.com/devstream-io/devstream/releases/download/v0.0.1/argocdapp_0.0.1.so"
		var targetDir = "tmp"

		It("returns an error when url is empty", func() {
			size, err := downloader.Download("", ".", targetDir)
			Expect(err).To(HaveOccurred())
			Expect(size).To(Equal(int64(0)))
		})

		It("returns an error when filename is [.]", func() {
			size, err := downloader.Download(url, ".", targetDir)
			Expect(err).To(HaveOccurred())
			Expect(size).To(Equal(int64(0)))
		})

		It("returns an error when filename is empty and download properly", func() {
			size, err := downloader.Download(url, "", targetDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(size).NotTo(Equal(int64(0)))
		})

		It("returns an error when filename is dir", func() {
			size, err := downloader.Download(url, "/", targetDir)
			Expect(err).To(HaveOccurred())
			Expect(size).To(Equal(int64(0)))
		})

		It("returns an error when the targetDir is empty", func() {
			size, err := downloader.Download(url, "download.txt", "")
			Expect(err).To(HaveOccurred())
			Expect(size).To(Equal(int64(0)))
		})

		It("returns an error when the url is not right", func() {
			size, err := downloader.Download(error_url, "download.txt", targetDir)
			Expect(err).To(HaveOccurred())
			Expect(size).To(Equal(int64(0)))
		})

		AfterEach(func() {
			err := os.RemoveAll(targetDir)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
