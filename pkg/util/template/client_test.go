package template

import (
	"fmt"
	"os"
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

var _ = Describe("template default funcs", func() {
	Context("getEnvInTemplate func", func() {
		var (
			tokenKey string
			existVal string
		)
		BeforeEach(func() {
			tokenKey = "TEMPLATE_ENV_TEST"
			existVal = os.Getenv(tokenKey)
			if existVal != "" {
				err := os.Unsetenv(tokenKey)
				Expect(err).ShouldNot(HaveOccurred())
			}
		})
		When("env not exist", func() {
			BeforeEach(func() {
				err := os.Unsetenv(tokenKey)
				Expect(err).ShouldNot(HaveOccurred())
			})
			It("should return err", func() {
				_, err := getEnvInTemplate(tokenKey)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("template can't get environment variable TEMPLATE_ENV_TEST, maybe you should set this environment first"))
			})
		})
		When("env exist", func() {
			BeforeEach(func() {
				err := os.Setenv(tokenKey, "test")
				Expect(err).ShouldNot(HaveOccurred())
			})
			It("should return err", func() {
				data, err := getEnvInTemplate(tokenKey)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(data).Should(Equal("test"))
			})
		})
		AfterEach(func() {
			if existVal != "" {
				os.Setenv(tokenKey, existVal)
			}
		})
	})
})
