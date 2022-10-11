package downloader_test

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/pkg/util/downloader"
)

func CreateFile(dir, filename string) (*os.File, error) {
	return os.Create(filepath.Join(dir, filename))
}

var _ = Describe("Downloader", func() {
	var (
		s                                  *ghttp.Server
		reqPath, failReqPath, tempDir, url string
		testContent                        []byte
	)
	BeforeEach(func() {
		tempDir = GinkgoT().TempDir()
		reqPath = "/test_plugin.so"
		failReqPath = "/not_exist.so"
		testContent = []byte("test Content")
		s = ghttp.NewServer()
		url = fmt.Sprintf("%s%s", s.URL(), reqPath)
	})
	AfterEach(func() {
		s.Close()
	})
	Context("Downloader test", func() {
		When("input params is wrong", func() {
			It("returns an error when url is empty", func() {
				size, err := downloader.New().Download("", ".", tempDir)
				Expect(err).To(HaveOccurred())
				Expect(size).To(Equal(int64(0)))
			})

			It("returns an error when filename is [.]", func() {
				size, err := downloader.New().Download(url, ".", tempDir)
				Expect(err).To(HaveOccurred())
				Expect(size).To(Equal(int64(0)))
			})

			It("returns an error when filename is dir", func() {
				size, err := downloader.New().Download(url, "/", tempDir)
				Expect(err).To(HaveOccurred())
				Expect(size).To(Equal(int64(0)))
			})

			It("returns an error when the targetDir is empty", func() {
				size, err := downloader.New().Download(url, "download.txt", "")
				Expect(err).To(HaveOccurred())
				Expect(size).To(Equal(int64(0)))
			})
		})
		When("remote request error", func() {
			BeforeEach(func() {
				s.RouteToHandler("GET", failReqPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", failReqPath),
					ghttp.RespondWithJSONEncoded(http.StatusNotFound, nil),
				))

			})
			It("returns an error when the url is not right", func() {
				errorURL := path.Join(s.URL(), failReqPath)
				size, err := downloader.New().Download(errorURL, "download.txt", tempDir)
				Expect(err).To(HaveOccurred())
				Expect(size).To(Equal(int64(0)))
			})
		})
		When("remote success", func() {
			BeforeEach(func() {
				s.RouteToHandler("GET", reqPath, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", reqPath),
					ghttp.RespondWith(http.StatusOK, string(testContent)),
				))
			})
			When("filename is empty", func() {
				It("should returns an error ", func() {
					size, err := downloader.New().Download(url, "", tempDir)
					Expect(err).NotTo(HaveOccurred())
					Expect(size).NotTo(Equal(int64(0)))
				})
			})
			When("fileName is right", func() {
				It("should get fileName with content", func() {
					fileName := "testFile"
					size, err := downloader.New().Download(url, fileName, tempDir)
					Expect(err).NotTo(HaveOccurred())
					Expect(size).NotTo(Equal(int64(0)))
					fileContent, err := os.ReadFile(filepath.Join(tempDir, fileName))
					Expect(err).Error().ShouldNot(HaveOccurred())
					Expect(fileContent).Should(Equal(testContent))
				})
			})
		})
	})
})

var _ = Describe("FetchContentFromURL func", func() {
	var (
		server                  *ghttp.Server
		testPath, remoteContent string
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
			_, err := downloader.FetchContentFromURL(reqURL)
			Expect(err).Error().Should(HaveOccurred())
		})
	})

	When("server return success", func() {
		BeforeEach(func() {
			remoteContent = "download content"
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", testPath),
					ghttp.RespondWith(http.StatusOK, remoteContent),
				),
			)
		})

		It("download the correct content", func() {
			reqURL := fmt.Sprintf("%s%s", server.URL(), testPath)
			content, err := downloader.FetchContentFromURL(reqURL)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(string(content)).Should(Equal(remoteContent))
		})
	})

	AfterEach(func() {
		server.Close()
	})
})
