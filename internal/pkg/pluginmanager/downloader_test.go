package pluginmanager

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("downloader", func() {
	var (
		server                                      *ghttp.Server
		downlodClient                               *PluginDownloadClient
		reqPath, pluginName, pluginVersion, tempDir string
	)

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		server = ghttp.NewServer()
		pluginName = "test_plugin"
		pluginVersion = "1.0"
		reqPath = fmt.Sprintf("/v%s/%s", pluginVersion, pluginName)
		downlodClient = NewPluginDownloadClient(server.URL())

	})

	Describe("download func", func() {
		When("server return err code", func() {
			BeforeEach(func() {
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", reqPath),
						ghttp.RespondWith(http.StatusNotFound, "test"),
					),
				)
			})

			It("should return err for download from url error", func() {
				err := downlodClient.download(tempDir, pluginName, pluginVersion)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("404"))
			})
		})

		When("response return success", func() {
			var testContent string

			BeforeEach(func() {
				testContent = "test"
				server.AppendHandlers(
					ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", reqPath),
						ghttp.RespondWith(http.StatusOK, testContent),
					),
				)
			})

			It("should return err if body error", func() {
				err := downlodClient.download(tempDir, pluginName, pluginVersion)
				Expect(err).Error().ShouldNot(HaveOccurred())
				downloadedFile := filepath.Join(tempDir, pluginName)
				// check plugin file is downloaded
				fileContent, err := os.ReadFile(downloadedFile)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(string(fileContent)).Should(Equal(testContent))
			})
		})
	})

	AfterEach(func() {
		server.Close()
	})
})
