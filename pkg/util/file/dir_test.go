package file_test

import (
	"errors"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/file"
)

var _ = Describe("GetFileMapByWalkDir func", func() {
	var (
		tempDir, exepectKey    string
		testContent            []byte
		mockFilterSuccessFunc  file.DirFileFilterFunc
		mockFilterFailedFunc   file.DirFileFilterFunc
		mockGetFileNameFunc    file.DirFileNameFunc
		mockProcessFailedFunc  file.DirFileContentFunc
		mockProcessSuccessFunc file.DirFileContentFunc
	)
	BeforeEach(func() {
		// create temp file for test
		tempDir = GinkgoT().TempDir()
		testContent = []byte("test content")
		exepectKey = "exepectKey"
		tempFile, err := os.CreateTemp(tempDir, "testFile")
		Expect(err).Error().ShouldNot(HaveOccurred())
		err = os.WriteFile(tempFile.Name(), testContent, 0755)
		Expect(err).Error().ShouldNot(HaveOccurred())
		// mock func to run
		mockFilterFailedFunc = func(filePath string, isDir bool) bool {
			return false
		}
		mockFilterSuccessFunc = func(filePath string, isDir bool) bool {
			return true
		}
		mockGetFileNameFunc = func(workDir, filePath string) string {
			return exepectKey
		}
		mockProcessFailedFunc = func(filePath string) ([]byte, error) {
			return []byte{}, errors.New("test err")
		}
		mockProcessSuccessFunc = func(filePath string) ([]byte, error) {
			return testContent, nil
		}

	})
	When("dir not exist", func() {
		It("should return error", func() {
			_, err := file.GetFileMapByWalkDir("not_exist_dir", mockFilterSuccessFunc, mockGetFileNameFunc, mockProcessSuccessFunc)
			Expect(err).Error().Should(HaveOccurred())
		})
	})
	When("filter return false", func() {
		It("should return empty", func() {
			result, err := file.GetFileMapByWalkDir(tempDir, mockFilterFailedFunc, mockGetFileNameFunc, mockProcessSuccessFunc)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(len(result)).Should(Equal(0))
		})
	})
	When("processFunc return false", func() {
		It("should return empty", func() {
			_, err := file.GetFileMapByWalkDir(tempDir, mockFilterSuccessFunc, mockGetFileNameFunc, mockProcessFailedFunc)
			Expect(err).Error().Should(HaveOccurred())
		})
	})
	When("all func work normal", func() {
		It("should return map with fileName and content", func() {
			result, err := file.GetFileMapByWalkDir(tempDir, mockFilterSuccessFunc, mockGetFileNameFunc, mockProcessSuccessFunc)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(len(result)).Should(Equal(1))
			val, ok := result[exepectKey]
			Expect(ok).Should(BeTrue())
			Expect(val).Should(Equal(testContent))
		})
	})
})
