package reposcaffolding

import (
	"fmt"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("utils for walk dir func", func() {
	var (
		tempDir, tempFilePath string
	)
	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
	})
	Context("filterGitFiles func", func() {
		When("input path isDir", func() {
			It("should return false", func() {
				result := filterGitFiles(tempDir, true)
				Expect(result).Should(BeFalse())
			})
		})
		When("input file contains invalid name", func() {
			It("should return false", func() {
				invalidGitFilePath := ".git/test"
				result := filterGitFiles(invalidGitFilePath, false)
				Expect(result).Should(BeFalse())
				invalidReadmeFile := "test/README.md"
				result = filterGitFiles(invalidReadmeFile, false)
				Expect(result).Should(BeFalse())
			})
		})
		When("input name is valid", func() {
			It("should return true", func() {
				validFileName := "test.go"
				result := filterGitFiles(validFileName, false)
				Expect(result).Should(BeTrue())
			})
		})
	})

	Context("getRepoFileNameFunc func", func() {
		It("should work as expected", func() {
			testRepoName := "test_repo"
			testAppName := "test_app"
			srcPath := "test"
			filePath := fmt.Sprintf("%s/%s/_app_name_.txt", srcPath, testRepoName)
			result := getRepoFileNameFunc(testAppName, testRepoName)(filePath, srcPath)
			Expect(result).Should(Equal(fmt.Sprintf("%s.txt", testAppName)))
		})
	})

	Context("processRepoFileFunc func", func() {
		var (
			testContent []byte
			testAppName string
		)
		BeforeEach(func() {
			testAppName = "test_app"
		})
		When("file is not tpl", func() {
			BeforeEach(func() {
				testContent = []byte("testContent")
				tempFile, err := os.CreateTemp(tempDir, "testFile")
				Expect(err).Error().ShouldNot(HaveOccurred())
				tempFilePath = tempFile.Name()
				err = os.WriteFile(tempFilePath, testContent, 0755)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
			It("should return file content", func() {
				result, err := processRepoFileFunc(testAppName, map[string]interface{}{})(tempFilePath)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(result).Should(Equal(testContent))
			})
		})
		When("file is tpl", func() {
			BeforeEach(func() {
				tempFile, err := os.CreateTemp(tempDir, "testFile.tpl")
				Expect(err).Error().ShouldNot(HaveOccurred())
				tempFilePath = tempFile.Name()
				testContent = []byte("[[ .Name ]]_test")
				err = os.WriteFile(tempFilePath, testContent, 0755)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
			It("should work as expected", func() {
				result, err := processRepoFileFunc(testAppName, map[string]interface{}{
					"Name": "devstream",
				})(tempFilePath)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(result).Should(Equal([]byte("devstream_test")))
			})
		})
	})
})
