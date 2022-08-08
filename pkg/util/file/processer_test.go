package file

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("unZipFileProcesser func", func() {
	var (
		zipFileName, tempDir, zipLocation, tempFile string
	)
	BeforeEach(func() {
		zipLocation = "test"
		tempFile = "testfile"
		tempDir = GinkgoT().TempDir()
		zipFile, err := os.CreateTemp(tempDir, "*.zip")
		zipFileName = zipFile.Name()
		Expect(err).Error().ShouldNot(HaveOccurred())
		defer zipFile.Close()
		writer := zip.NewWriter(zipFile)
		defer writer.Close()

		newFile, err := os.CreateTemp(tempDir, tempFile)
		Expect(err).Error().ShouldNot(HaveOccurred())
		defer newFile.Close()
		w1, _ := writer.Create(fmt.Sprintf("%s/%s", zipLocation, tempFile))
		_, err = io.Copy(w1, newFile)
		Expect(err).Error().ShouldNot(HaveOccurred())
	})

	It("should work", func() {
		dstPath, err := unZipFileProcesser(zipFileName)
		Expect(err).Error().ShouldNot(HaveOccurred())
		dirFiles, err := ioutil.ReadDir(dstPath)
		Expect(err).Error().ShouldNot(HaveOccurred())
		Expect(len(dirFiles)).Should(Equal(1))
		Expect(dirFiles[0].Name()).Should(Equal(zipLocation))
		zipDirFiles, err := ioutil.ReadDir(filepath.Join(dstPath, zipLocation))
		Expect(err).Error().ShouldNot(HaveOccurred())
		Expect(len(dirFiles)).Should(Equal(1))
		Expect(zipDirFiles[0].Name()).Should(Equal(tempFile))
	})
})
