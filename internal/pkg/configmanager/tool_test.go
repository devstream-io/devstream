package configmanager

import (
	"fmt"
	"runtime"

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

var _ = Describe("Tool struct", func() {
	var (
		t                *Tool
		name, instanceID string
		opts             RawOptions
	)
	BeforeEach(func() {
		name = "test"
		instanceID = "test_instance"
		opts = RawOptions{"test": "gg"}
		t = &Tool{
			Name:       name,
			InstanceID: instanceID,
			Options:    opts,
		}
	})
	Context("DeepCopy method", func() {
		It("should copy struct", func() {
			x := t.DeepCopy()
			Expect(x).Should(Equal(t))
		})
	})
	Context("GetPluginName method", func() {
		It("should return valid", func() {
			Expect(t.GetPluginName()).Should(Equal(
				fmt.Sprintf("%s-%s-%s_%s", t.Name, runtime.GOOS, runtime.GOARCH, ""),
			))
		})
	})
	Context("GetPluginFileName method", func() {
		It("should return valid", func() {
			Expect(t.GetPluginFileName()).Should(
				Equal(fmt.Sprintf("%s-%s-%s_%s.so", t.Name, runtime.GOOS, runtime.GOARCH, "")),
			)
		})
	})
	Context("GetPluginFileNameWithOSAndArch method", func() {
		It("should return valid", func() {
			Expect(t.GetPluginFileNameWithOSAndArch("test", "plug")).Should(Equal("test-test-plug_.so"))
		})
	})
	Context("GetPluginMD5FileName method", func() {
		It("should return valid", func() {
			Expect(t.GetPluginMD5FileName()).Should(
				Equal(fmt.Sprintf("%s-%s-%s_%s.md5", t.Name, runtime.GOOS, runtime.GOARCH, "")),
			)
		})
	})
	Context("GetPluginMD5FileNameWithOSAndArch method", func() {
		It("should return valid", func() {
			Expect(t.GetPluginMD5FileNameWithOSAndArch("test", "plug")).Should(Equal("test-test-plug_.md5"))
		})
	})
})
