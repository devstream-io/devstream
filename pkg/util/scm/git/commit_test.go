package git_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("CommitInfo struct", func() {
	var (
		fileContent      []byte
		dstPath, tempDir string
		filePaths        []*git.GitFilePathInfo
	)

	Context("GetFileContent method", func() {
		When("file not exist", func() {
			BeforeEach(func() {
				filePaths = []*git.GitFilePathInfo{
					{
						SourcePath:      "not_exist_file",
						DestinationPath: dstPath,
					}}
			})
			It("should return empty map", func() {
				fileMap := git.GetFileContent(filePaths)
				Expect(len(fileMap)).Should(Equal(0))
			})
		})
		When("file exist", func() {
			BeforeEach(func() {
				fileContent = []byte("this file content")
				tempDir = GinkgoT().TempDir()
				dstPath = filepath.Join(tempDir, "dstFile")
				testFile, err := os.CreateTemp(tempDir, "test")
				Expect(err).Error().ShouldNot(HaveOccurred())
				err = os.WriteFile(testFile.Name(), fileContent, 0755)
				Expect(err).Error().ShouldNot(HaveOccurred())
				filePaths = []*git.GitFilePathInfo{
					{
						SourcePath:      testFile.Name(),
						DestinationPath: dstPath,
					}}
			})
			It("should return empty map", func() {
				fileMap := git.GetFileContent(filePaths)
				Expect(fileMap).ShouldNot(BeEmpty())
				result, ok := fileMap[dstPath]
				Expect(ok).Should(BeTrue())
				Expect(result).Should(Equal(fileContent))
			})
		})
	})
})

var _ = Describe("GenerateGitFileInfo func", func() {
	var (
		tempDir, tempFileLoc, gitDir string
		testContent                  []byte
	)
	BeforeEach(func() {
		gitDir = "git"
		testContent = []byte("this is test content")
		tempDir = GinkgoT().TempDir()
		tempFile, err := os.CreateTemp(tempDir, "test")
		tempFileLoc = tempFile.Name()
		Expect(err).Error().ShouldNot(HaveOccurred())
		err = os.WriteFile(tempFileLoc, testContent, 0755)
		Expect(err).Error().ShouldNot(HaveOccurred())
	})
	When("path not exist", func() {
		It("should return err", func() {
			_, err := git.GenerateGitFileInfo([]string{"not_exist"}, "")
			Expect(err).Should(HaveOccurred())
		})
	})
	When("filePath type is file", func() {
		It("should return file", func() {
			fileInfo, err := git.GenerateGitFileInfo([]string{tempFileLoc}, gitDir)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(fileInfo)).Should(Equal(1))
			Expect(fileInfo[0].DestinationPath).Should(Equal(filepath.Join(gitDir, tempFileLoc)))
		})
	})
	When("filePath type is dir", func() {
		It("should return dirFiles", func() {
			fileInfo, err := git.GenerateGitFileInfo([]string{tempDir}, gitDir)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(len(fileInfo)).Should(Equal(1))
			fileName := filepath.Base(tempFileLoc)
			Expect(fileInfo[0].DestinationPath).Should(Equal(filepath.Join(gitDir, fileName)))
		})
	})
})
