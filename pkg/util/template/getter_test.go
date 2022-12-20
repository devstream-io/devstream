package template_test

import (
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/pkg/util/template"
)

var _ = Describe("Getters(localFile, content, url)", func() {

	const src = "test_data"

	When("template content is local file", func() {
		It("should return rendered template", func() {
			file, err := os.CreateTemp("", "test")
			Expect(err).To(Succeed())
			_, err = file.WriteString(src)
			Expect(err).To(Succeed())

			data, err := template.LocalFileGetter(file.Name())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(src))
		})
	})

	When("template content is content", func() {
		It("should return rendered template", func() {
			data, err := template.ContentGetter(src)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(src))
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
					ghttp.RespondWith(http.StatusOK, src),
				),
			)
		})

		It("should return rendered template", func() {
			data, err := template.URLGetter(server.URL() + testPath)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(string(data)).Should(Equal(src))
		})

		AfterEach(func() {
			server.Close()
		})
	})
})
