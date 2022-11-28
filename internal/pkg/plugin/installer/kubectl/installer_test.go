package kubectl

import (
	"fmt"
	"io"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	utilKubectl "github.com/devstream-io/devstream/pkg/util/kubectl"
)

var _ = Describe("renderKubectlContent", func() {
	var (
		content string
		options configmanager.RawOptions
	)

	BeforeEach(func() {
		content = `metadata:
  name: "[[ .app.name ]]"
  namespace: "[[ .app.namespace ]]"
  finalizers:
    - resources-finalizer.argocd.argoproj.io
`
		options = map[string]interface{}{
			"app": map[string]interface{}{
				"name":      "app-name",
				"namespace": "app-namespace",
			},
		}
	})

	It("should render kubectl content", func() {

		contentExpected := `metadata:
  name: "app-name"
  namespace: "app-namespace"
  finalizers:
    - resources-finalizer.argocd.argoproj.io
`
		reader, err := renderKubectlContent(content, options)
		Expect(err).To(Succeed())

		bytes, err := io.ReadAll(reader)

		Expect(err).To(Succeed())
		Expect(string(bytes)).To(Equal(contentExpected))
	})

})

var _ = Describe("ProcessByContent", Ordered, func() {
	var options configmanager.RawOptions
	var s *ghttp.Server

	BeforeAll(func() {
		s = ghttp.NewServer()
		s.RouteToHandler("GET", "/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "ok")
		})
	})
	AfterAll(func() {
		s.Close()
	})
	It("should return error if content is empty", func() {
		op := ProcessByContent(utilKubectl.Apply, "")
		err := op(options)
		Expect(err).To(HaveOccurred())
	})
	It("action is kubectl apply", func() {
		op := ProcessByURL(utilKubectl.Apply, s.URL())
		err := op(options)
		Expect(err).To(HaveOccurred())
	})
	It("action is kubectl create", func() {
		op := ProcessByURL(utilKubectl.Create, s.URL())
		err := op(options)
		Expect(err).To(HaveOccurred())
	})
	It("action is not support", func() {
		op := ProcessByURL("", s.URL())
		err := op(options)
		Expect(err).To(HaveOccurred())
	})
})
