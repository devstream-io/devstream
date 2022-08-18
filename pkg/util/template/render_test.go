package template

import (
	"strings"
	"text/template"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Render", func() {
	var (
		tpl, content string
		vars         map[string]interface{}
		funcMaps     []template.FuncMap
		err          error
	)

	var (
		wrongVars = map[string]interface{}{
			"App": map[string]interface{}{
				"name": "test",
			},
		}

		rightVars = map[string]interface{}{
			"App": map[string]interface{}{
				"Name":      "test",
				"NameSpace": "test_namespace",
			},
		}
	)

	const (
		wrongTpl = `metadata:[[`
		rightTpl = `metadata:
  name: "[[ .App.Name ]]"
  namespace: "[[ .App.NameSpace ]]"`
		tplWithFunc = `metadata:
  name: "[[ .App.Name ]]"
  namespace: "[[ .App.NameSpace ]]"
  funcMap:
    len: [[ "abc" | len ]]
    equal: [[ eq "test" .App.Name ]]
    upper: "[[ .App.Name | upper ]]"`
		templateName = "template"
	)

	JustBeforeEach(func() {
		content, err = Render(templateName, tpl, vars, funcMaps...)
	})

	When("template content wrong", func() {
		BeforeEach(func() {
			tpl = wrongTpl
		})

		It("should return err", func() {
			Expect(err).To(HaveOccurred())
		})
	})

	When("template content right", func() {
		BeforeEach(func() {
			tpl = rightTpl
		})

		When("vars is wrong", func() {
			BeforeEach(func() {
				vars = wrongVars
			})

			It("should return err", func() {
				Expect(err).To(HaveOccurred())
			})
		})

		When("vars is right", func() {
			BeforeEach(func() {
				vars = rightVars
			})

			When("funcMaps is nil", func() {
				BeforeEach(func() {
					funcMaps = nil
				})

				It("should return content", func() {

					expected := `metadata:
  name: "test"
  namespace: "test_namespace"`
					Expect(err).To(Succeed())
					Expect(content).To(Equal(expected))
				})
			})

			When("funcMaps is not nil", func() {
				BeforeEach(func() {
					tpl = tplWithFunc

					upper := func(s string) string {
						return strings.ToUpper(s)
					}
					funcMaps = []template.FuncMap{{"upper": upper}}
				})

				It("should return content", func() {
					expected := `metadata:
  name: "test"
  namespace: "test_namespace"
  funcMap:
    len: 3
    equal: true
    upper: "TEST"`
					Expect(err).To(Succeed())
					Expect(content).To(Equal(expected))
				})
			})
		})

	})
})
