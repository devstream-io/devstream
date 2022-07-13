package local_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/devstream-io/devstream/internal/pkg/backend/local"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewLocal func", func() {
	var tFile, tFileLoc, tempDir string

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		tFile = "test_state_file"
		tFileLoc = filepath.Join(tempDir, tFile)
	})

	When("input filename", func() {
		It("should create state file", func() {
			local.NewLocal(tFileLoc)
			_, err := os.Stat(tFileLoc)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
	})
})

var _ = Describe("Local struct", func() {
	var (
		tFile, tFileLoc, tempDir string
		tLocal                   *local.Local
	)

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		tFile = "test_state_file"
		tFileLoc = filepath.Join(tempDir, tFile)
		tLocal = local.NewLocal(tFileLoc)
	})

	Describe("Read method", func() {
		var testData []byte
		BeforeEach(func() {
			testData = []byte("this is test data")
			err := ioutil.WriteFile(tFileLoc, testData, 0644)
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
			fileData, err := ioutil.ReadFile(tFileLoc)
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
