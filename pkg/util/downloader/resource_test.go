package downloader_test

import (
	"fmt"
	"net/http"
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/pkg/util/downloader"
)

var _ = Describe("ResourceClient struct", func() {
	var (
		source           string
		resourceLocation downloader.ResourceLocation
	)
	Context("GetWithGoGetter method", func() {
		When("source is local files or directory", func() {
			BeforeEach(func() {
				source = GinkgoT().TempDir()
				resourceLocation = downloader.ResourceLocation(source)
			})
			It("should return local source path directly", func() {
				dstPath, err := resourceLocation.Download()
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(dstPath).Should(Equal(source))
			})
		})
		When("download resource from web", func() {
			var (
				s       *ghttp.Server
				reqPath string
			)
			BeforeEach(func() {
				s = ghttp.NewServer()
			})
			When("resource return error", func() {
				BeforeEach(func() {
					reqPath = path.Join(s.URL(), "not_exist_resource")
					resourceLocation = downloader.ResourceLocation(reqPath)
				})
				It("should return err", func() {
					_, err := resourceLocation.Download()
					Expect(err).Error().Should(HaveOccurred())
					Expect(err.Error()).Should(ContainSubstring(fmt.Sprintf("get resource files %s failed", reqPath)))
				})
			})
			When("resource return normal", func() {
				var (
					returnContent string
				)
				BeforeEach(func() {
					returnContent = "test getter"
					pathName := "/exist_resource"
					reqPath = fmt.Sprintf("%s%s", s.URL(), pathName)

					s.RouteToHandler("HEAD", pathName, func(w http.ResponseWriter, r *http.Request) {
						fmt.Fprint(w, "ok")
					})
					s.RouteToHandler("GET", pathName, ghttp.CombineHandlers(
						ghttp.VerifyRequest("GET", pathName),
						ghttp.RespondWith(http.StatusOK, returnContent),
					))
					resourceLocation = downloader.ResourceLocation(reqPath)
				})
				When("destination is not setted", func() {
					It("should create temp dir and download file", func() {
						dstPath, err := resourceLocation.Download()
						Expect(err).Error().ShouldNot(HaveOccurred())
						files, err := os.ReadDir(dstPath)
						Expect(err).Error().ShouldNot(HaveOccurred())
						Expect(len(files)).Should(Equal(1))
						file := files[0]
						Expect(file.Name()).Should(Equal("exist_resource"))
						filePath := path.Join(dstPath, file.Name())
						content, err := os.ReadFile(filePath)
						Expect(err).Error().ShouldNot(HaveOccurred())
						Expect(string(content)).Should(Equal(returnContent))
					})
				})
			})
			AfterEach(func() {
				s.Close()
			})
		})
	})
})
