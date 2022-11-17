package configmanager

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Dependency", func() {
	Context("singe dep", func() {
		tools := Tools{
			{InstanceID: "argocd", Name: "argocd"},
			{InstanceID: "argocdapp", Name: "argocdapp", DependsOn: []string{"argocd.argocd"}},
		}
		errors := tools.validateDependency()
		Expect(len(errors)).To(Equal(0))
	})

	Context("dep not exist", func() {
		tools := Tools{
			{InstanceID: "argocdapp", Name: "argocdapp", DependsOn: []string{"argocd.argocd"}},
		}
		errors := tools.validateDependency()
		Expect(len(errors)).To(Equal(1))
	})

	Context("multi-dep", func() {
		tools := Tools{
			{InstanceID: "argocd", Name: "argocd"},
			{InstanceID: "repo", Name: "github"},
			{InstanceID: "argocdapp", Name: "argocdapp", DependsOn: []string{"argocd.argocd", "github.repo"}},
		}
		errors := tools.validateDependency()
		Expect(len(errors)).To(Equal(0))
	})

	Context("empty dep", func() {
		tools := Tools{
			{InstanceID: "argocd", Name: "argocd"},
			{InstanceID: "argocdapp", Name: "argocdapp", DependsOn: []string{}},
		}
		errors := tools.validateDependency()
		Expect(len(errors)).To(Equal(0))
	})
})

var _ = Describe("Tool Validation", func() {
	It("should return empty error array if tools all valid", func() {
		tools := Tools{
			{Name: "test_tool", InstanceID: "0", DependsOn: []string{}},
		}
		errors := tools.validate()
		Expect(errors).Should(BeEmpty())
	})
	It("should return error if tool not valid", func() {
		tools := Tools{
			{Name: "", InstanceID: "", DependsOn: []string{}},
		}
		errors := tools.validate()
		Expect(errors).ShouldNot(BeEmpty())
	})

})
