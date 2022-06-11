package downloader_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/downloader"
)

var _ = Describe("Downloader", func() {
	Context("Downloader test", func() {
		var url = "https://github.com/devstream-io/devstream/releases/download/v0.0.1/argocdapp_0.0.1.so"
		var targetDir = "tmp"

		It("returns error if filename is empty and can't be parsed from the url", func() {
			size, err := downloader.Download(url, ".", targetDir)
			Expect(err).To(HaveOccurred())
			Expect(size).To(Equal(int64(0)))
		})

		It("returns an error when url and filename are empty", func() {
			size, err := downloader.Download("", "", targetDir)
			Expect(err).To(HaveOccurred())
			Expect(size).To(Equal(int64(0)))
		})

		It("should download the file properly", func() {
			size, err := downloader.Download(url, "", targetDir)
			Expect(err).NotTo(HaveOccurred())
			Expect(size).NotTo(Equal(int64(0)))
		})

		AfterEach(func() {
			err := os.RemoveAll(targetDir)
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
