package local_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/backend/local"
)

var _ = Describe("NewLocal func", func() {
	var tFile, tempDir string

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		tFile = "test_state_file"
	})

	When("input filename", func() {
		It("should create state file", func() {
			_, err := local.NewLocal(tempDir, tFile)
			Expect(err).Should(Succeed())
			_, err = os.Stat(filepath.Join(tempDir, tFile))
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})
})

var _ = Describe("Local struct", func() {
	var (
		tFile, tFileLoc, tempDir string
		tLocal                   *local.Local
		err                      error
	)

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		tFile = "test_state_file"
		tLocal, err = local.NewLocal(tempDir, tFile)
		Expect(err).ShouldNot(HaveOccurred())
	})

	Describe("Read method", func() {
		var testData []byte

		BeforeEach(func() {
			testData = []byte("this is test data")
			err := os.WriteFile(tFileLoc, testData, 0644)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})

		It("should return file content", func() {
			fileData, err := tLocal.Read()
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(fileData).Should(Equal(testData))
		})
	})

	Describe("Write method", func() {
		var writeData []byte

		BeforeEach(func() {
			writeData = []byte("this is write test")
		})

		It("should write  data to file", func() {
			err := tLocal.Write(writeData)
			Expect(err).Error().ShouldNot(HaveOccurred())
			fileData, err := os.ReadFile(tFileLoc)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(fileData).Should(Equal(writeData))
		})
	})

	// After each test, clean file content
	AfterEach(func() {
		err := os.Truncate(tFileLoc, 0)
		Expect(err).Error().ShouldNot(HaveOccurred())
	})
})
