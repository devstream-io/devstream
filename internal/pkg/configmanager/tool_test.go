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

var _ = Describe("Tool struct", func() {
	var tools Tools

	const (
		toolName   = "test_tool"
		instanceID = "test_instance"
	)
	Context("renderInstanceIDtoOptions method", func() {
		BeforeEach(func() {
			tools = Tools{
				{Name: toolName, InstanceID: instanceID},
			}
		})
		When("tool option is null", func() {
			It("should set nil to RawOptions", func() {
				tools.renderInstanceIDtoOptions()
				Expect(len(tools)).Should(Equal(1))
				tool := tools[0]
				Expect(tool.Options).Should(Equal(RawOptions{
					"instanceID": instanceID,
				}))
			})
		})
	})
})

var _ = Describe("duplicatedCheck", func() {
	var (
		errs                   []error
		tools                  Tools
		toolsWithoutDuplicated = Tools{
			{Name: "test_tool", InstanceID: "0"},
			{Name: "test_tool", InstanceID: "1"},
			{Name: "test_tool", InstanceID: "2"},
		}

		toolsWithDuplicated = Tools{
			{Name: "test_tool", InstanceID: "0"},
			{Name: "test_tool", InstanceID: "1"},
			{Name: "test_tool", InstanceID: "0"},
		}
	)

	JustBeforeEach(func() {
		errs = tools.duplicatedCheck()
	})

	When("tools has duplicated name and instanceID", func() {
		BeforeEach(func() {
			tools = toolsWithDuplicated
		})
		It("should return error", func() {
			Expect(errs).Should(HaveLen(1))
		})
	})

	When("tools don't have duplicated name and instanceID", func() {
		BeforeEach(func() {
			tools = toolsWithoutDuplicated
		})
		It("should return nil", func() {
			Expect(errs).Should(BeEmpty())
		})
	})
})
