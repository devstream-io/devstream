package file

import (
	"fmt"
	"io/ioutil"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("getFileFromURL", func() {
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
			_, err := getFileFromURL(reqURL)
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

		It("should create file with content", func() {
			reqURL := fmt.Sprintf("%s%s", server.URL(), testPath)
			fileName, err := getFileFromURL(reqURL)
			Expect(err).Error().ShouldNot(HaveOccurred())
			fileContent, err := ioutil.ReadFile(fileName)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(string(fileContent)).Should(Equal(remoteContent))
		})
	})

	AfterEach(func() {
		server.Close()
	})
})
