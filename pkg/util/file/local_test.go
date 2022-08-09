package file

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("getFileFromLocal func", func() {
	var (
		tempDir      string
		notExistFile string
		existFile    string
	)
	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		notExistFile = filepath.Join(tempDir, "not_exist")
		existFile = filepath.Join(tempDir, "exist")
		f, err := os.Create(existFile)
		Expect(err).Error().ShouldNot(HaveOccurred())
		defer f.Close()
	})

	When("file not exist", func() {
		It("should contain error", func() {
			_, err := getFileFromLocal(notExistFile)
			Expect(err).Error().Should(HaveOccurred())
		})
	})

	When("file exist", func() {
		It("should config normal", func() {
			loc, err := getFileFromLocal(existFile)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(loc).Should(Equal(existFile))
		})
	})
})

var _ = Describe("getFileFromContent func", func() {
	var (
		testContent string
	)
	BeforeEach(func() {
		testContent = "This is a test Content"

	})

	It("should return a file for content", func() {
		fileName, _ := getFileFromContent(testContent)
		// check file exist
		content, err := os.ReadFile(fileName)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(string(content)).Should(Equal(testContent))
	})
})
