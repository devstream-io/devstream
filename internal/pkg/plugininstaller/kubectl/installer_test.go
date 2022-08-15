package kubectl_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/kubectl"
	utilKubectl "github.com/devstream-io/devstream/pkg/util/kubectl"
)

var _ = Describe("ProcessByContent", Ordered, func() {
	var options plugininstaller.RawOptions
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
		op := kubectl.ProcessByContent(utilKubectl.Apply, "")
		err := op(options)
		Expect(err).To(HaveOccurred())
	})
	It("action is kubectl apply", func() {
		op := kubectl.ProcessByURL(utilKubectl.Apply, s.URL())
		err := op(options)
		Expect(err).To(HaveOccurred())
	})
	It("action is kubectl create", func() {
		op := kubectl.ProcessByURL(utilKubectl.Create, s.URL())
		err := op(options)
		Expect(err).To(HaveOccurred())
	})
	It("action is kubectl delete", func() {
		op := kubectl.ProcessByURL(utilKubectl.Delete, s.URL())
		err := op(options)
		Expect(err).To(HaveOccurred())
	})
	It("action is not support", func() {
		op := kubectl.ProcessByURL("", s.URL())
		err := op(options)
		Expect(err).To(HaveOccurred())
	})
})
