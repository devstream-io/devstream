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

var _ = Describe("MD5", func() {
	var (
		err                                           error
		config                                        *configmanager.Config
		tempDir, file, fileMD5, filePath, fileMD5Path string
		tools                                         []configmanager.Tool
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
		tools = []configmanager.Tool{
			{InstanceID: "a", Name: "a"},
		}
		config = &configmanager.Config{Tools: tools}

		file = configmanager.GetPluginFileName(&tools[0])
		filePath = filepath.Join(tempDir, file)
		fileMD5 = configmanager.GetPluginMD5FileName(&tools[0])
		fileMD5Path = filepath.Join(tempDir, fileMD5)
		err := createNewFile(filePath)
		Expect(err).NotTo(HaveOccurred())

		err = addMD5File(filePath, fileMD5Path)
		Expect(err).NotTo(HaveOccurred())

	})

	Describe("CheckLocalPlugins func", func() {
		It("should match .md5 file content", func() {
			err = CheckLocalPlugins(config)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should mismatch .md5 file content", func() {
			err = createNewFile(fileMD5Path)
			Expect(err).NotTo(HaveOccurred())
			err = CheckLocalPlugins(config)
			expectErrMsg := fmt.Sprintf("plugin %s doesn't match with .md5", tools[0].InstanceID)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).To(Equal(expectErrMsg))
		})
	})

	Describe("pluginAndMD5Matches func", func() {
		It("should match .md5 file content", func() {
			err = pluginAndMD5Matches(tempDir, file, fileMD5, tools[0].InstanceID)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})

		It("should mismatch .md5 file content", func() {
			err = createNewFile(fileMD5Path)
			Expect(err).NotTo(HaveOccurred())
			err = pluginAndMD5Matches(tempDir, file, fileMD5, tools[0].InstanceID)
			expectErrMsg := fmt.Sprintf("plugin %s doesn't match with .md5", tools[0].InstanceID)
			Expect(err.Error()).To(Equal(expectErrMsg))
		})
	})

	Describe("redownloadPlugins func", func() {
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
				err = redownloadPlugins(
					pbDownloadClient, tempDir, notExistPluginFile, fileMD5, pluginVersion,
				)
				Expect(err).Error().Should(HaveOccurred())
				err = redownloadPlugins(
					pbDownloadClient, tempDir, file, notExistPluginFile, pluginVersion,
				)
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
				err = redownloadPlugins(
					pbDownloadClient, tempDir, file, fileMD5, pluginVersion,
				)
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

	Describe("DownloadPlugins func", func() {
		When("pluginDir is Empty", func() {
			It("should return err", func() {
				viper.Set("plugin-dir", "")
				pluginDir := viper.GetString("plugin-dir")
				err = DownloadPlugins(config.Tools, pluginDir, runtime.GOOS, runtime.GOARCH)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("plugins directory should not be "))
			})
		})
	})
})
