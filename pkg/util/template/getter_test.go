package template

import (
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Getters(localFile, content, url)", func() {
	var (
		getter ContentGetter

		src                string // content of source, each getter will use this content
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
		// here is the core of the test
		// because each getter is the same type and the expected result is the same,
		// we can use "JustAfterEach" to test them all at end of each test case
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
			// then "JustAfterEach" will test the result
		})
	})

	When("template content is content", func() {
		It("should return rendered template", func() {
			getter = FromContent(src)
			// then "JustAfterEach" will test the result
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
			// then "JustAfterEach" will test the result
		})

		AfterEach(func() {
			server.Close()
		})
	})
})
