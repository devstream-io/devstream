package template

import (
	"fmt"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Getters(localFile, content, url)", func() {
	var (
		getter ContentGetter

		src                string
		rendered, expected []byte
		err                error
	)

	const (
		str = "abc123\ndevstream\ntest"
	)

	BeforeEach(func() {
		expected = []byte(str)
		src = str
	})

	JustAfterEach(func() {
		rendered, err = getter()
		Expect(err).To(Succeed())
		Expect(rendered).To(Equal(expected))
	})

	When("template content is local file", func() {
		It("should return rendered template", func() {
			file, err := os.CreateTemp("", "test")
			Expect(err).To(Succeed())
			_, err = file.WriteString(src)
			Expect(err).To(Succeed())

			getter = FromLocalFile(file.Name())
		})
	})

	When("template content is content", func() {
		It("should return rendered template", func() {
			getter = FromContent(src)
		})
	})

	When("template content is url", func() {
		var (
			server   *ghttp.Server
			testPath string
		)

		BeforeEach(func() {
			testPath = "/testPath"

			server = ghttp.NewServer()
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", testPath),
					ghttp.RespondWith(http.StatusOK, str),
				),
			)
		})

		It("should return rendered template", func() {
			getter = FromURL(server.URL() + testPath)
		})

		AfterEach(func() {
			server.Close()
		})
	})
})

var _ = Describe("getContentFromURL", func() {
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
			_, err := getContentFromURL(reqURL)
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
			content, err := getContentFromURL(reqURL)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(string(content)).Should(Equal(remoteContent))
		})
	})

	AfterEach(func() {
		server.Close()
	})
})
