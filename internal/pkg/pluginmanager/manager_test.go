package pluginmanager

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/md5"
)

var _ = Describe("downloadPlugins func", func() {

	var (
		server  *mockPluginServer
		tools   configmanager.Tools
		tempDir string
	)

	const version = "0.1.0"

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		server = newMockPluginServer()
	})

	When("plugins are right", func() {
		const (
			argocd  = "argocd"
			jenkins = "jenkins"
		)
		BeforeEach(func() {
			tools = configmanager.Tools{
				{Name: argocd},
				{Name: jenkins},
			}
			server.registerPluginOK(argocd, argocd+"content", version, runtime.GOOS, runtime.GOARCH)
			server.registerPluginOK(jenkins, jenkins+"content", version, runtime.GOOS, runtime.GOARCH)
		})

		It("should download plugins successfully", func() {
			err := downloadPlugins(server.URL(), tools, tempDir, runtime.GOOS, runtime.GOARCH, version)
			Expect(err).ShouldNot(HaveOccurred())
			viper.Set("plugin-dir", tempDir)
			err = CheckLocalPlugins(tools)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	When("pluginDir is Empty", func() {
		It("should return err", func() {
			err := downloadPlugins(server.URL(), tools, "", runtime.GOOS, runtime.GOARCH, version)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("plugins directory should not be "))
		})
	})

	When("plugin is not exist", func() {
		BeforeEach(func() {
			const invalidPlugin = "invalidPlugin"
			tools = configmanager.Tools{
				{Name: invalidPlugin},
			}
			server.registerPluginNotFound(invalidPlugin, version, runtime.GOOS, runtime.GOARCH)
		})
		It("should return err", func() {

			err := downloadPlugins(server.URL(), tools, tempDir, runtime.GOOS, runtime.GOARCH, version)
			Expect(err).Should(HaveOccurred())
		})
	})
})

var _ = Describe("MD5", func() {
	var (
		err                                           error
		config                                        *configmanager.Config
		tempDir, file, fileMD5, filePath, fileMD5Path string
		tools                                         configmanager.Tools
	)

	createNewFile := func(fileName string) error {
		f1, err := os.Create(fileName)
		if err != nil {
			return err
		}
		defer f1.Close()
		return nil
	}

	addMD5File := func(fileName, md5FileName string) error {
		md5, err := md5.CalcFileMD5(fileName)
		if err != nil {
			return err
		}
		md5File, err := os.Create(md5FileName)
		if err != nil {
			return err
		}
		_, err = md5File.Write([]byte(md5))
		if err != nil {
			return err
		}
		return nil
	}

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		viper.Set("plugin-dir", tempDir)
		tools = configmanager.Tools{
			{InstanceID: "a", Name: "a"},
		}
		config = &configmanager.Config{Tools: tools}

		file = tools[0].GetPluginFileName()
		filePath = filepath.Join(tempDir, file)
		fileMD5 = tools[0].GetPluginMD5FileName()
		fileMD5Path = filepath.Join(tempDir, fileMD5)
		err := createNewFile(filePath)
		Expect(err).NotTo(HaveOccurred())

		err = addMD5File(filePath, fileMD5Path)
		Expect(err).NotTo(HaveOccurred())

	})

	Describe("CheckLocalPlugins func", func() {
		It("should match .md5 file content", func() {
			err = CheckLocalPlugins(config.Tools)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should mismatch .md5 file content", func() {
			err = createNewFile(fileMD5Path)
			Expect(err).NotTo(HaveOccurred())
			err = CheckLocalPlugins(config.Tools)
			expectErrMsg := fmt.Sprintf("plugin %s doesn't match with .md5", tools[0].InstanceID)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).To(Equal(expectErrMsg))
		})
	})

	Describe("ifPluginAndMD5Match func", func() {
		It("should match .md5 file content", func() {
			matched, err := ifPluginAndMD5Match(tempDir, file, fileMD5)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(matched).To(BeTrue())
		})

		It("should mismatch .md5 file content", func() {
			err = createNewFile(fileMD5Path)
			Expect(err).NotTo(HaveOccurred())
			matched, err := ifPluginAndMD5Match(tempDir, file, fileMD5)
			Expect(err).ToNot(HaveOccurred())
			Expect(matched).To(BeFalse())
		})
	})

	Describe("reDownload func", func() {
		var (
			pbDownloadClient *PluginDownloadClient
			pluginVersion    string
		)

		BeforeEach(func() {
			pbDownloadClient = NewPluginDownloadClient("not_exist_url")
		})

		When("pluginFile not exist", func() {
			It("should return error", func() {
				notExistPluginFile := "not_exist_plugin"
				err = pbDownloadClient.reDownload(tempDir, notExistPluginFile, fileMD5, pluginVersion)
				Expect(err).Error().Should(HaveOccurred())
				err = pbDownloadClient.reDownload(tempDir, file, notExistPluginFile, pluginVersion)
				Expect(err).Error().Should(HaveOccurred())
			})
		})

		When("old plugin exist", func() {
			var (
				server     *ghttp.Server
				rspContent string
			)

			BeforeEach(func() {
				pluginVersion = "1.0"
				reqPath := fmt.Sprintf("/v%s/%s", pluginVersion, file)
				md5Path := fmt.Sprintf("/v%s/%s", pluginVersion, fileMD5)
				server = ghttp.NewServer()
				pbDownloadClient = NewPluginDownloadClient(server.URL())
				rspContent = "reDownload Content"
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", reqPath),
						ghttp.RespondWith(http.StatusOK, rspContent),
					),
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", md5Path),
						ghttp.RespondWith(http.StatusOK, rspContent),
					),
				)
			})

			It("should re download success if download success", func() {
				err = pbDownloadClient.reDownload(tempDir, file, fileMD5, pluginVersion)
				Expect(err).Error().ShouldNot(HaveOccurred())
				newFileContent, err := os.ReadFile(filePath)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(string(newFileContent)).Should(Equal(rspContent))
			})

			AfterEach(func() {
				server.Close()
			})
		})
	})
})
