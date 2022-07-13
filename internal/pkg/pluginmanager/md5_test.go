package pluginmanager

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("manager", func() {
	var tempDir, fileName, fileContent string

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		fileName = "testMd5.temp"
		fileContent = "testContent"
	})

	Describe("LocalContentMD5 func", func() {
		It("should return md5 content", func() {
			filePath := filepath.Join(tempDir, fileName)
			f, err := os.Create(filePath)
			Expect(err).Error().NotTo(HaveOccurred())
			_, err = f.Write([]byte(fileContent))
			Expect(err).Error().NotTo(HaveOccurred())
			result, err := LocalContentMD5(filePath)
			Expect(err).Error().NotTo(HaveOccurred())
			// testContent md5 val is uulB4NHN9Ct11tDva9fSWg==
			Expect(result).To(Equal("uulB4NHN9Ct11tDva9fSWg=="))
		})
	})
})
