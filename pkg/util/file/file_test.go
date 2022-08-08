package file

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("file struct", func() {
	var (
		templateConfig *TemplateConfig
	)
	BeforeEach(func() {
		templateConfig = NewTemplate()
	})

	Context("FromLocal method", func() {
		It("should set getter method and info", func() {
			testLocation := "test_location"
			Expect(templateConfig.getter).Should(BeNil())
			config := templateConfig.FromLocal(testLocation)
			Expect(config.getter).ShouldNot(BeNil())
			Expect(config.info).Should(ContainSubstring(testLocation))
		})
	})

	Context("FromContent method", func() {
		It("should set getter method and info", func() {
			testContent := "this is a test content"
			Expect(templateConfig.getter).Should(BeNil())
			config := templateConfig.FromContent(testContent)
			Expect(config.getter).ShouldNot(BeNil())
			Expect(config.info).Should(ContainSubstring(testContent))
		})
	})

	Context("FromRemote method", func() {
		It("should set url method and info", func() {
			Expect(templateConfig.getter).Should(BeNil())
			testURL := "http://www.test.com"
			config := templateConfig.FromRemote(testURL)
			Expect(config.getter).ShouldNot(BeNil())
			Expect(config.info).Should(Equal(testURL))
		})
	})

	Context("Run method", func() {
		var (
			inputFileName        string
			outputFileName       string
			mockErr              error
			mockFunc             fileProcesser
			mockFailedFunc       fileProcesser
			mockRenderFunc       renderProcesser
			mockRenderFailedFunc renderProcesser
		)
		BeforeEach(func() {
			inputFileName = "test"
			outputFileName = "changed_file_name"
			mockErr = errors.New("mock test")
			mockFunc = func(input string) (string, error) {
				return outputFileName, nil
			}
			mockFailedFunc = func(input string) (string, error) {
				return "", mockErr
			}
			mockRenderFunc = func(templateName, filePath string, vars map[string]interface{}) (string, error) {
				return filePath, nil
			}
			mockRenderFailedFunc = func(templateName, filePath string, vars map[string]interface{}) (string, error) {
				return "", mockErr
			}
		})

		When("getter method process error", func() {
			It("should return err", func() {
				templateConfig.info = inputFileName
				templateConfig.getter = mockFailedFunc
				_, err := templateConfig.Run()
				Expect(err).Error().Should(HaveOccurred())

			})
		})

		When("render method process error", func() {
			It("should return err", func() {
				templateConfig.info = inputFileName
				templateConfig.getter = mockFunc
				templateConfig.render = mockRenderFailedFunc
				_, err := templateConfig.Run()
				Expect(err).Error().Should(HaveOccurred())
			})
		})

		When("struct has error", func() {
			It("should return err driectly", func() {
				errMsg := "test_err"
				templateConfig.processErr = errors.New(errMsg)
				_, err := templateConfig.Run()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(Equal(errMsg))
			})
			It("should return empty error", func() {
				templateConfig.info = ""
				_, err := templateConfig.Run()
				Expect(err).Error().Should(HaveOccurred())
			})
			AfterEach(func() {
				templateConfig = NewTemplate()
			})
		})
		When("all config is right", func() {
			It("should return name and path", func() {
				templateConfig.info = inputFileName
				templateConfig.getter = mockFunc
				templateConfig.render = mockRenderFunc
				output, err := templateConfig.Run()
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(output).Should(Equal(outputFileName))
			})
		})
	})
})

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
		err := ioutil.WriteFile(srcPath, testContent, 0666)
		Expect(err).Error().ShouldNot(HaveOccurred())
		err = CopyFile(srcPath, dstPath)
		Expect(err).Error().ShouldNot(HaveOccurred())
		data, err := ioutil.ReadFile(dstPath)
		Expect(err).Error().ShouldNot(HaveOccurred())
		Expect(data).Should(Equal(testContent))
	})
})
