package configmanager

import (
	"fmt"
	"runtime"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var tools Tools

var _ = Describe("validateDependsOnConfig", func() {
	When("empty dep", func() {
		BeforeEach(func() {
			tools = Tools{
				{InstanceID: "ins-1", Name: "plugin1"},
				{InstanceID: "ins-2", Name: "plugin2"},
			}
		})
		It("should not have errors", func() {
			errs := tools.validateDependsOnConfig()
			Expect(len(errs)).To(Equal(0))
		})
	})

	When("singe dep", func() {
		BeforeEach(func() {
			tools = Tools{
				{InstanceID: "ins-1", Name: "plugin1"},
				{InstanceID: "ins-2", Name: "plugin2"},
			}
		})
		It("should not have errors", func() {
			tools[1].DependsOn = []string{"plugin1.ins-1"}
			errs := tools.validateDependsOnConfig()
			Expect(len(errs)).To(Equal(0))
		})
		It("should has some errors", func() {
			tools[1].DependsOn = []string{"plugin1.ins-2"}
			errs := tools.validateDependsOnConfig()
			Expect(len(errs)).To(Equal(1))
		})
	})

	When("multi-dep", func() {
		BeforeEach(func() {
			tools = Tools{
				{InstanceID: "ins-1", Name: "plugin1"},
				{InstanceID: "ins-2", Name: "plugin2"},
				{InstanceID: "ins-3", Name: "plugin3"},
			}
		})
		It("should not have errors", func() {
			tools[2].DependsOn = []string{"plugin1.ins-1"}
			tools[2].DependsOn = []string{"plugin2.ins-2"}
			tools[1].DependsOn = []string{"plugin1.ins-1"}
			errs := tools.validateDependsOnConfig()
			Expect(len(errs)).To(Equal(0))
		})
		It("should has some errors", func() {
			tools[1].DependsOn = []string{"plugin1.ins-3"}
			tools[2].DependsOn = []string{"plugin1.ins-2", "plugin2.ins-1"}
			errs := tools.validateDependsOnConfig()
			Expect(len(errs)).To(Equal(3))
		})
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

var _ = Describe("getToolsFromConfigFileWithVarsRendered", func() {
	const toolsConfig = `---
tools:
- name: name1
  instanceID: ins-1
  dependsOn: []
  options:
    foo: [[ foo ]]
`
	When("get tools from config file", func() {
		It("should return config with vars", func() {
			tools, err := getToolsFromConfigFileWithVarsRendered([]byte(toolsConfig), map[string]any{"foo": interface{}("bar")})
			Expect(err).NotTo(HaveOccurred())
			Expect(tools).NotTo(BeNil())
			Expect(len(tools)).To(Equal(1))
			Expect(len(tools[0].Options)).To(Equal(1))
			Expect(tools[0].Options["foo"]).To(Equal(interface{}("bar")))
		})
	})
})
