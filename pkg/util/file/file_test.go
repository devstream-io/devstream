package file_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/file"
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
		err = file.CopyFile(srcPath, dstPath)
		Expect(err).Error().ShouldNot(HaveOccurred())
		data, err := os.ReadFile(dstPath)
		Expect(err).Error().ShouldNot(HaveOccurred())
		Expect(data).Should(Equal(testContent))
	})
})

var _ = Describe("GenerateAbsFilePath func", func() {
	var baseDir, fileName string
	When("file not exist", func() {
		BeforeEach(func() {
			baseDir = "not_exist"
			fileName = "not_exist"
		})
		It("should return err", func() {
			_, err := file.GenerateAbsFilePath(baseDir, fileName)
			Expect(err).Error().Should(HaveOccurred())
		})
	})
	When("file exist", func() {
		BeforeEach(func() {
			baseDir = GinkgoT().TempDir()
			testFile, err := os.CreateTemp(baseDir, "test")
			Expect(err).Error().ShouldNot(HaveOccurred())
			fileName = filepath.Base(testFile.Name())
		})
		It("should return absPath", func() {
			path, err := file.GenerateAbsFilePath(baseDir, fileName)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(path).Should(Equal(filepath.Join(baseDir, fileName)))
		})
	})
})

var _ = Describe("GetFileAbsDirPath func", func() {
	var baseDir, fileName string
	When("file exist", func() {
		BeforeEach(func() {
			baseDir = GinkgoT().TempDir()
			testFile, err := os.CreateTemp(baseDir, "test")
			Expect(err).Error().ShouldNot(HaveOccurred())
			fileName = filepath.Base(testFile.Name())
		})
		It("should return absPath", func() {
			path, err := file.GetFileAbsDirPath(filepath.Join(baseDir, fileName))
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(path).Should(Equal(baseDir))
		})
	})
})

var _ = Describe("GetFileAbsDirPathOrDirItself func", func() {
	var baseDir, fileName string
	When("param if file", func() {
		BeforeEach(func() {
			baseDir = GinkgoT().TempDir()
			testFile, err := os.CreateTemp(baseDir, "test")
			Expect(err).Error().ShouldNot(HaveOccurred())
			fileName = filepath.Base(testFile.Name())
		})
		It("should return parent directory of file", func() {
			path, err := file.GetFileAbsDirPathOrDirItself(filepath.Join(baseDir, fileName))
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(path).Should(Equal(baseDir))
		})
	})
	When("param is dir", func() {
		BeforeEach(func() {
			baseDir = GinkgoT().TempDir()
		})
		It("should return dir itself", func() {
			path, err := file.GetFileAbsDirPathOrDirItself(baseDir)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(path).Should(Equal(baseDir))
		})
	})
})
