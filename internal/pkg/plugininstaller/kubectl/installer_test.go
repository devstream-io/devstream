package kubectl_test

import (
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"github.com/devstream-io/devstream/internal/pkg/plugininstaller"
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/kubectl"
	"github.com/devstream-io/devstream/pkg/util/file"
	utilKubectl "github.com/devstream-io/devstream/pkg/util/kubectl"
)

var _ = Describe("ProcessByContent", Ordered, func() {
	var templateConfig *file.TemplateConfig
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
	It("templateConfig with info is empty string", func() {
		templateConfig = file.NewTemplate()
		op := kubectl.ProcessByContent(utilKubectl.Apply, templateConfig)
		err := op(options)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("content is not setted"))
	})
	It("action is kubectl apply", func() {
		templateConfig = file.NewTemplate().FromRemote(s.URL())
		op := kubectl.ProcessByContent(utilKubectl.Apply, templateConfig)
		err := op(options)
		Expect(err).To(HaveOccurred())
	})
	It("action is kubectl create", func() {
		templateConfig = file.NewTemplate().FromRemote(s.URL())
		op := kubectl.ProcessByContent(utilKubectl.Create, templateConfig)
		err := op(options)
		Expect(err).To(HaveOccurred())
	})
	It("action is kubectl delete", func() {
		templateConfig = file.NewTemplate().FromRemote(s.URL())
		op := kubectl.ProcessByContent(utilKubectl.Delete, templateConfig)
		err := op(options)
		Expect(err).To(HaveOccurred())
	})
	It("action is not support", func() {
		templateConfig = file.NewTemplate().FromRemote(s.URL())
		op := kubectl.ProcessByContent("", templateConfig)
		err := op(options)
		Expect(err).To(HaveOccurred())
	})
})
