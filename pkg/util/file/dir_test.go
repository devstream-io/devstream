package file_test

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/pkg/util/file"
)

var _ = Describe("WalkDir func", func() {
	var (
		tempDir, exepectKey    string
		testContent            []byte
		mockFilterSuccessFunc  file.DirFIleFilterFunc
		mockFilterFailedFunc   file.DirFIleFilterFunc
		mockGetFileNameFunc    file.DirFileNameFunc
		mockProcessFailedFunc  file.DirFileProcessFunc
		mockProcessSuccessFunc file.DirFileProcessFunc
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
			_, err := file.WalkDir("not_exist_dir", mockFilterSuccessFunc, mockGetFileNameFunc, mockProcessSuccessFunc)
			Expect(err).Error().Should(HaveOccurred())
		})
	})
	When("filter return false", func() {
		It("should return empty", func() {
			result, err := file.WalkDir(tempDir, mockFilterFailedFunc, mockGetFileNameFunc, mockProcessSuccessFunc)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(len(result)).Should(Equal(0))
		})
	})
	When("processFunc return false", func() {
		It("should return empty", func() {
			result, err := file.WalkDir(tempDir, mockFilterSuccessFunc, mockGetFileNameFunc, mockProcessFailedFunc)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(len(result)).Should(Equal(0))
		})
	})
	When("all func work normal", func() {
		It("should return map with fileName and content", func() {
			result, err := file.WalkDir(tempDir, mockFilterSuccessFunc, mockGetFileNameFunc, mockProcessSuccessFunc)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(len(result)).Should(Equal(1))
			val, ok := result[exepectKey]
			Expect(ok).Should(BeTrue())
			Expect(val).Should(Equal(testContent))
		})
	})
})

var _ = Describe("unzip func", func() {
	var (
		zipFileName, tempDir, zipLocation, tempFile string
	)
	BeforeEach(func() {
		zipLocation = "test"
		tempFile = "testfile"
		tempDir = GinkgoT().TempDir()
		zipFile, err := os.CreateTemp(tempDir, "*.zip")
		zipFileName = zipFile.Name()
		Expect(err).Error().ShouldNot(HaveOccurred())
		defer zipFile.Close()
		writer := zip.NewWriter(zipFile)
		defer writer.Close()

		newFile, err := os.CreateTemp(tempDir, tempFile)
		Expect(err).Error().ShouldNot(HaveOccurred())
		defer newFile.Close()
		w1, _ := writer.Create(fmt.Sprintf("%s/%s", zipLocation, tempFile))
		_, err = io.Copy(w1, newFile)
		Expect(err).Error().ShouldNot(HaveOccurred())
	})

	It("should work", func() {
		dstPath, err := file.Unzip(zipFileName)
		Expect(err).Error().ShouldNot(HaveOccurred())
		dirFiles, err := os.ReadDir(dstPath)
		Expect(err).Error().ShouldNot(HaveOccurred())
		Expect(len(dirFiles)).Should(Equal(1))
		Expect(dirFiles[0].Name()).Should(Equal(zipLocation))
		zipDirFiles, err := os.ReadDir(filepath.Join(dstPath, zipLocation))
		Expect(err).Error().ShouldNot(HaveOccurred())
		Expect(len(dirFiles)).Should(Equal(1))
		Expect(zipDirFiles[0].Name()).Should(Equal(tempFile))
	})
})

var _ = Describe("DownloadAndUnzipFile func", func() {
	var (
		server   *ghttp.Server
		testPath string
	)

	BeforeEach(func() {
		testPath = "/testPath"
		server = ghttp.NewServer()
	})

	When("server return error code", func() {
		BeforeEach(func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", testPath),
					ghttp.RespondWith(http.StatusNotFound, ""),
				),
			)

		})
		It("should return err", func() {
			reqURL := fmt.Sprintf("%s%s", server.URL(), testPath)
			_, err := file.DownloadAndUnzipFile(reqURL)
			Expect(err).Error().Should(HaveOccurred())
		})
	})
	//TODO(steinliber): Add more Unzip func tests
	AfterEach(func() {
		server.Close()
	})
})
