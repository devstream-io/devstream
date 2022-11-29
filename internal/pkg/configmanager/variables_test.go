package configmanager

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("renderConfigWithVariables", func() {
	var vars = map[string]interface{}{
		"foo1": "bar1",
		"foo2": "bar2",
	}
	When("render a config with variables", func() {
		It("should works fine", func() {
			result1, err1 := renderConfigWithVariables("[[ foo1 ]]/[[ foo2]]", vars)
			result2, err2 := renderConfigWithVariables("a[[ foo1 ]]/[[ foo2]]b", vars)
			result3, err3 := renderConfigWithVariables(" [[ foo1 ]] [[ foo2]] ", vars)
			Expect(err1).NotTo(HaveOccurred())
			Expect(err2).NotTo(HaveOccurred())
			Expect(err3).NotTo(HaveOccurred())
			Expect(result1).To(Equal([]byte("bar1/bar2")))
			Expect(result2).To(Equal([]byte("abar1/bar2b")))
			Expect(result3).To(Equal([]byte(" bar1 bar2 ")))
		})
	})
})
