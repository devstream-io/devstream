package configloader_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
)

var _ = Describe("Tool Config", func() {
	Context("Validate method", func() {
		It("should return empty if tool is valid", func() {
			tool := configloader.Tool{Name: "test", InstanceID: "test"}
			err := tool.Validate()
			Expect(err).Should(BeEmpty())
		})

		It("should return error if tool is not valid", func() {
			tool := configloader.Tool{Name: "！test", InstanceID: "！test"}
			err := tool.Validate()
			Expect(err).ShouldNot(BeEmpty())
		})
	})

	Context("Key method", func() {
		It("should return joined key", func() {
			toolName := "test_tool"
			toolInstance := "0"
			tool := configloader.Tool{Name: toolName, InstanceID: toolInstance}
			key := tool.Key()
			Expect(key).Should(Equal(fmt.Sprintf("%s.%s", toolName, toolInstance)))
		})
	})

	Context("DeepCopy method", func() {
		It("should return with different address", func() {
			tool := configloader.Tool{Name: "test"}
			copyTool := tool.DeepCopy()
			Expect(copyTool.Name).Should(Equal(tool.Name))
			Expect(copyTool).ShouldNot(Equal(tool))
		})
	})

	Context("NewToolWithToolConfigBytesAndVarsConfigBytes func", func() {
		toolConfigBytes := []byte(`
tools:
  - name: github-repo-scaffolding-golang
    instanceID: 0
    options:
      owner: [[ owner ]]
      repo: go-webapp-devstream-demo
      branch: main
      image_repo: YOUR_DOCKER_USERNAME/go-webapp-devstream-demo`)

		It("should return error if yaml not valid", func() {
			varBytes := []byte(`ownertest_owner`)
			_, err := configloader.NewToolWithToolConfigBytesAndVarsConfigBytes(toolConfigBytes, varBytes)
			Expect(err).Error().Should(HaveOccurred())
		})

		It("should set variable in tool option", func() {
			varBytes := []byte(`
owner: test_owner`)
			tools, err := configloader.NewToolWithToolConfigBytesAndVarsConfigBytes(toolConfigBytes, varBytes)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(len(tools)).Should(Equal(1))
			owner, exist := tools[0].Options["owner"]
			Expect(exist).Should(BeTrue())
			Expect(fmt.Sprintf("%v", owner)).Should(ContainSubstring("test_owner"))
		})
	})
})
