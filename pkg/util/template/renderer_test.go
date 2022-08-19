package template

import (
	"strings"
	"text/template"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Renderer", func() {
	var (
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
	)

	It("should render template successfully", func() {
		upper := func(s string) string {
			return strings.ToUpper(s)
		}
		funcMaps := []template.FuncMap{{"upper": upper}}

		addNewLineProcessor := func(b []byte) ([]byte, error) {
			return append(b, '\n'), nil
		}

		content, err := New().
			FromContent(tplWithFunc).
			AddProcessor(addNewLineProcessor).
			AddProcessor(addNewLineProcessor).
			SetDefaultRender(tplName, vars, funcMaps...).
			Render()

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
