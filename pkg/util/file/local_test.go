package file

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
