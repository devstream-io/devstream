package file

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("CopyFile func", func() {
	var (
		tempDir, srcPath, dstPath string
		testContent               []byte
	)

	BeforeEach(func() {
		testContent = []byte("test_content")
		tempDir = GinkgoT().TempDir()
		srcPath = filepath.Join(tempDir, "src")
		dstPath = filepath.Join(tempDir, "dst")
		f1, err := os.Create(srcPath)
		Expect(err).Error().ShouldNot(HaveOccurred())
		defer f1.Close()
		f2, err := os.Create(dstPath)
		Expect(err).Error().ShouldNot(HaveOccurred())
		defer f2.Close()
	})

	It("should copy content form src to dst", func() {
		err := os.WriteFile(srcPath, testContent, 0666)
		Expect(err).Error().ShouldNot(HaveOccurred())
		err = CopyFile(srcPath, dstPath)
		Expect(err).Error().ShouldNot(HaveOccurred())
		data, err := os.ReadFile(dstPath)
		Expect(err).Error().ShouldNot(HaveOccurred())
		Expect(data).Should(Equal(testContent))
	})
})
