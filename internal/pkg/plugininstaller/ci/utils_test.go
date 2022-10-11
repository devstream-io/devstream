package ci

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("processCIFilesFunc func", func() {
	var (
		tempFileLoc string
		testContent []byte
	)
	BeforeEach(func() {
		tempDir := GinkgoT().TempDir()
		tempFile, err := os.CreateTemp(tempDir, "testFile")
		Expect(err).Error().ShouldNot(HaveOccurred())
		tempFileLoc = tempFile.Name()
	})
	When("file is right", func() {
		BeforeEach(func() {
			testContent = []byte("[[ .Name ]]_test")
			err := os.WriteFile(tempFileLoc, testContent, 0755)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
		It("should work as expected", func() {
			result, err := processCIFilesFunc("test", map[string]interface{}{
				"Name": "devstream",
			})(tempFileLoc)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(result).Should(Equal([]byte("devstream_test")))
		})
	})
	When("file is not exist", func() {
		It("should return error", func() {
			_, err := processCIFilesFunc("test", map[string]interface{}{
				"Name": "devstream",
			})("not_exist_file")
			Expect(err).Error().Should(HaveOccurred())
		})
	})
})
