package template

import (
	"fmt"
	"strings"
	"text/template"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("renderClient", func() {
	var (
		tplName, tplWithFunc string
		vars                 map[string]interface{}
		client               *renderClient
	)
	When("all func is used", func() {
		BeforeEach(func() {
			tplName = "template"
			tplWithFunc = `metadata:
  name: "[[ .App.Name ]]"
  namespace: "[[ .App.NameSpace ]]"
  funcMap:
    len: [[ "abc" | len ]]
    equal: [[ eq "test" .App.Name ]]
    upper: "[[ .App.Name | upper ]]"`
			vars = map[string]interface{}{
				"App": map[string]interface{}{
					"Name":      "test",
					"NameSpace": "test_namespace",
				},
			}

		})

		It("should render template successfully", func() {
			upper := func(s string) string {
				return strings.ToUpper(s)
			}
			funcMap := template.FuncMap{"upper": upper}

			addNewLineProcessor := func(b []byte) []byte {
				return append(b, '\n')
			}

			content, err := NewRenderClient(
				&TemplateOption{
					Name:    tplName,
					FuncMap: funcMap,
				}, ContentGetter, addNewLineProcessor, addNewLineProcessor,
			).Render(tplWithFunc, vars)

			expected := `metadata:
  name: "test"
  namespace: "test_namespace"
  funcMap:
    len: 3
    equal: true
    upper: "TEST"

`
			Expect(err).NotTo(HaveOccurred())
			Expect(content).To(Equal(expected))
		})
	})
	When("getter error", func() {
		BeforeEach(func() {
			getterErrorFunc := func(inputStr string) ([]byte, error) {
				return nil, fmt.Errorf("test error")
			}
			client = &renderClient{
				Getter: getterErrorFunc,
			}
		})
		It("should return error", func() {
			_, err := client.Render("test", nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(Equal("test error"))
		})
	})
	When("parse template failed", func() {
		BeforeEach(func() {
			client = &renderClient{
				Getter: ContentGetter,
			}
		})
		It("should return error", func() {
			invalidTemplateData := `[[ .... not invalid ]] [[ notvalid data ]]`
			_, err := client.Render(invalidTemplateData, nil)
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("parse template: default_template"))
		})
	})

	When("render template failed", func() {
		BeforeEach(func() {
			client = &renderClient{
				Getter: ContentGetter,
			}
		})
		It("should return error", func() {
			invalidTemplateData := `[[ .variable ]] `
			_, err := client.Render(invalidTemplateData, map[string]interface{}{})
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("render template: default_template"))
		})
	})

})
