package file

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("renderFile func", func() {
	var (
		tempFilePath, templateContent string
	)
	BeforeEach(func() {
		templateContent = `
			metadata:
			  name: "[[ .App.Name ]]"
			  namespace: "[[ .App.NameSpace ]]"`
		tempFile, err := os.CreateTemp("", "test_temp")
		Expect(err).Error().ShouldNot(HaveOccurred())
		defer tempFile.Close()
		tempFilePath = tempFile.Name()
		err = ioutil.WriteFile(tempFilePath, []byte(templateContent), 0666)
		Expect(err).Error().ShouldNot(HaveOccurred())
	})

	When("srcPath file is not exist", func() {
		It("should return err", func() {
			_, err := renderFile("test_app", "not_exist_path", map[string]interface{}{})
			Expect(err).Error().Should(HaveOccurred())
		})
	})

	When("input vars is empty", func() {
		It("should return err", func() {
			_, err := renderFile("test_app", tempFilePath, map[string]interface{}{})
			Expect(err).Error().Should(HaveOccurred())
		})
	})

	When("input vars and right srcPath", func() {
		var rightContent string
		BeforeEach(func() {
			rightContent = `
			metadata:
			  name: "test"
			  namespace: "test_namespace"`
		})

		It("should work normal", func() {
			dstPath, err := renderFile("test_app", tempFilePath, map[string]interface{}{
				"App": map[string]interface{}{
					"Name":      "test",
					"NameSpace": "test_namespace",
				},
			})
			Expect(err).Error().ShouldNot(HaveOccurred())
			content, err := ioutil.ReadFile(dstPath)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(string(content)).Should(Equal(rightContent))
		})
	})
})

var _ = Describe("replaceAppNameInPathStr func", func() {
	var (
		placeHolder string
		filePath    string
		appName     string
	)
	BeforeEach(func() {
		placeHolder = "__app__"
		appName = "test"
	})
	When("filePath not contains placeHolder", func() {
		BeforeEach(func() {
			filePath = "/app/dev"
		})
		It("should return same filePath", func() {
			newPath, err := replaceAppNameInPathStr(filePath, placeHolder, appName)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(newPath).Should(Equal(filePath))
		})
	})
	When("filPath contains placeHolder", func() {
		BeforeEach(func() {
			filePath = fmt.Sprintf("app/%s/dev", placeHolder)
		})
		It("should replace placeHolder with app name", func() {
			newPath, err := replaceAppNameInPathStr(filePath, placeHolder, appName)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(newPath).Should(Equal(fmt.Sprintf("app/%s/dev", appName)))
		})
	})
})

var _ = Describe("renderGitRepoDir func", func() {
	var (
		vars                                                        map[string]interface{}
		srcPath, contentDir, rawContent, tplContent, renderdContent string
	)

	createFile := func(filePath, content string) {
		f, err := os.Create(filePath)
		Expect(err).Error().ShouldNot(HaveOccurred())
		defer f.Close()
		err = ioutil.WriteFile(filePath, []byte(content), 0755)
		Expect(err).Error().ShouldNot(HaveOccurred())
	}
	createDir := func(dirPath string) {
		err := os.Mkdir(dirPath, 0755)
		Expect(err).Error().ShouldNot(HaveOccurred())
	}
	BeforeEach(func() {
		rawContent = "This is a file without template variable"
		tplContent = `
			metadata:
			  name: "[[ .App.Name ]]"
			  namespace: "[[ .App.NameSpace ]]"`
		renderdContent = `
			metadata:
			  name: "test"
			  namespace: "test_namespace"`

	})

	When("srcPath is not exist", func() {
		BeforeEach(func() {
			srcPath = "not_exist_path"
		})
		It("should return err", func() {
			_, err := renderGitRepoDir("test", srcPath, vars)
			Expect(err).Error().Should(HaveOccurred())
		})
	})
	When("all config is right", func() {
		BeforeEach(func() {
			contentDir = "content"
			vars = map[string]interface{}{
				"App": map[string]interface{}{
					"Name":      "test",
					"NameSpace": "test_namespace",
				},
			}
			srcPath = GinkgoT().TempDir()
			gitPath := filepath.Join(srcPath, ".git")
			createDir(gitPath)
			createFile(filepath.Join(gitPath, "gitFile"), tplContent)
			createFile(filepath.Join(srcPath, "README.md"), "")
			contentDirPath := filepath.Join(srcPath, contentDir)
			createDir(contentDirPath)
			createFile(filepath.Join(contentDirPath, "test.yaml.tpl"), tplContent)
			createFile(filepath.Join(contentDirPath, "raw.txt"), rawContent)
		})
		It("should render all dir", func() {
			dstPath, err := renderGitRepoDir("test", srcPath, vars)
			Expect(err).Error().ShouldNot(HaveOccurred())
			files, err := ioutil.ReadDir(dstPath)
			Expect(err).Error().ShouldNot(HaveOccurred())
			// test README.md dir is not copied
			Expect(len(files)).Should(Equal(2))
			// test git dir files should not copied
			gitDirFiles, err := ioutil.ReadDir(filepath.Join(dstPath, ".git"))
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(len(gitDirFiles)).Should(Equal(0))
			// test content dir files is copied
			contentDirLoc := filepath.Join(dstPath, contentDir)
			contentFiles, err := ioutil.ReadDir(contentDirLoc)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(len(contentFiles)).Should(Equal(2))
			// test file content
			tplFileContent, err := ioutil.ReadFile(filepath.Join(contentDirLoc, "test.yaml"))
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(string(tplFileContent)).Should(Equal(renderdContent))
			rawFileContent, err := ioutil.ReadFile(filepath.Join(contentDirLoc, "raw.txt"))
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(string(rawFileContent)).Should(Equal(rawContent))
		})
	})
})
