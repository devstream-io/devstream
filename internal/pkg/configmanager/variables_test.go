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

var _ = Describe("getVarsFromConfigFile", func() {
	const configFile = `---
foo: bar
vars:
  foo1: bar1
  foo2: 123
foo3: bar3
`
	When("get vars from config file", func() {
		It("should works fine", func() {
			varMap, err := getVarsFromConfigFile([]byte(configFile))
			Expect(err).NotTo(HaveOccurred())
			Expect(varMap).NotTo(BeNil())
			Expect(len(varMap)).To(Equal(2))
			Expect(varMap["foo1"]).To(Equal(interface{}("bar1")))
			Expect(varMap["foo2"]).To(Equal(interface{}(123)))
		})
	})
})
