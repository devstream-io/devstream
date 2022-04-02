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
		It("Should get the file", func() {
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
