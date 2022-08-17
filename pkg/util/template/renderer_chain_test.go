package template

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RenderChainExample", func() {
	It("should work", func() {
		res, err := example()

		Expect(res).To(Equal("xxx this is var"))
		Expect(err).To(Succeed())
	})

})

func example() (string, error) {
	contentSourceType := "xxx"
	src := "xxx [[.Var]]"

	t := New()
	var g WithGetter

	switch contentSourceType {
	case "url":
		g = t.FromURL(src)
	case "file":
		g = t.FromLocalFile(src)
	case "content":
		g = t.FromContent(src)
	default:
		g = t.FromContent(src)
	}

	renderMethod := "xxx"
	vars := struct {
		Var string
	}{
		Var: "this is var",
	}
	var r WithRender

	switch renderMethod {
	case "xxx":
		r = g.SetRender(DefaultRender("otherRender", vars))
	case "default":
		r = g.SetRender(DefaultRender("name", vars))
	}

	// res: xxx this is var
	return r.Render()
}
