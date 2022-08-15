package file

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

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
		err = os.WriteFile(filePath, []byte(content), 0755)
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
			files, err := os.ReadDir(dstPath)
			Expect(err).Error().ShouldNot(HaveOccurred())
			// test README.md dir is not copied
			Expect(len(files)).Should(Equal(2))
			// test git dir files should not copied
			gitDirFiles, err := os.ReadDir(filepath.Join(dstPath, ".git"))
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(len(gitDirFiles)).Should(Equal(0))
			// test content dir files is copied
			contentDirLoc := filepath.Join(dstPath, contentDir)
			contentFiles, err := os.ReadDir(contentDirLoc)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(len(contentFiles)).Should(Equal(2))
			// test file content
			tplFileContent, err := os.ReadFile(filepath.Join(contentDirLoc, "test.yaml"))
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(string(tplFileContent)).Should(Equal(renderdContent))
			rawFileContent, err := os.ReadFile(filepath.Join(contentDirLoc, "raw.txt"))
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(string(rawFileContent)).Should(Equal(rawContent))
		})
	})
})
