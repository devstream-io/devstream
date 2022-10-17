package pluginmanager

import (
	"os"
	"path/filepath"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

var _ = Describe("downloader", func() {
	var (
		server  *mockPluginServer
		pdc     *PluginDownloadClient
		tempDir string
	)

	const (
		pluginName    = "test_plugin"
		pluginVersion = "0.1.0"
	)

	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		server = newMockPluginServer()
		pdc = NewPluginDownloadClient(server.URL())
	})

	Describe("download func", func() {
		When("server return err code", func() {
			BeforeEach(func() {
				server.registerPluginNotFound(pluginName, pluginVersion, runtime.GOOS, runtime.GOARCH)
			})

			It("should return err for download from url error", func() {
				tool := &configmanager.Tool{Name: pluginName}
				pluginFileName := tool.GetPluginFileNameWithOSAndArch(runtime.GOOS, runtime.GOARCH)
				err := pdc.download(tempDir, pluginFileName, pluginVersion)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("404"))
			})
		})

		When("response return success", func() {
			var testContent string

			BeforeEach(func() {
				testContent = "test"
				server.registerPluginOK(pluginName, testContent, pluginVersion, runtime.GOOS, runtime.GOARCH)
			})

			It("should download plugins successfully", func() {
				// plugin file name and md5 file name
				tool := &configmanager.Tool{Name: pluginName}
				pluginFileName := tool.GetPluginFileNameWithOSAndArch(runtime.GOOS, runtime.GOARCH)
				pluginMD5FileName := tool.GetPluginMD5FileNameWithOSAndArch(runtime.GOOS, runtime.GOARCH)

				// download .so file
				err := pdc.download(tempDir, pluginFileName, pluginVersion)
				Expect(err).ShouldNot(HaveOccurred())
				// check plugin file is downloaded
				fileContent, err := os.ReadFile(filepath.Join(tempDir, pluginFileName))
				Expect(err).ShouldNot(HaveOccurred())
				Expect(string(fileContent)).Should(Equal(testContent))

				// download .md5 file
				err = pdc.download(tempDir, pluginMD5FileName, pluginVersion)
				Expect(err).ShouldNot(HaveOccurred())
				md5Matched, err := ifPluginAndMD5Match(tempDir, pluginFileName, pluginMD5FileName)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(md5Matched).Should(BeTrue())
			})
		})
	})

	AfterEach(func() {
		server.Close()
	})
})
