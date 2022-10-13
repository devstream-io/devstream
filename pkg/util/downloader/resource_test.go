package downloader_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/downloader"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("ResourceClient struct", func() {
	var (
		source string
		c      *downloader.ResourceClient
	)
	BeforeEach(func() {
		c = &downloader.ResourceClient{}
	})
	Context("GetWithGoGetter method", func() {
		When("source is local files or directory", func() {
			BeforeEach(func() {
				source = GinkgoT().TempDir()
				c.Source = source
			})
			It("should return local source path directly", func() {
				dstPath, err := c.GetWithGoGetter()
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
					c.Source = reqPath
				})
				It("should return err", func() {
					_, err := c.GetWithGoGetter()
					Expect(err).Error().Should(HaveOccurred())
					Expect(err.Error()).Should(ContainSubstring("http: no Host in request URL"))
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
					c.Source = reqPath
				})
				When("destination is not setted", func() {
					It("should create temp dir and download file", func() {
						dstPath, err := c.GetWithGoGetter()
						Expect(err).Error().ShouldNot(HaveOccurred())
						files, err := ioutil.ReadDir(dstPath)
						Expect(err).Error().ShouldNot(HaveOccurred())
						Expect(len(files)).Should(Equal(1))
						file := files[0]
						Expect(file.Name()).Should(Equal("exist_resource"))
						filePath := path.Join(dstPath, file.Name())
						content, err := os.ReadFile(filePath)
						Expect(err).Error().ShouldNot(HaveOccurred())
						Expect(string(content)).Should(Equal(returnContent))
						c.CleanUp()
					})
				})
			})
			AfterEach(func() {
				s.Close()
			})
		})
	})
})
