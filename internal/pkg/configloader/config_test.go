package configloader_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configloader"
)

var _ = Describe("Config", func() {
	Describe("LoadConfig yaml", func() {
		configStateObj, err := configloader.LoadConfig("../../../examples/quickstart.yaml")
		Context("when the Yaml parses successfully", func() {
			Specify("should state filed correctly", func() {
				Expect(configStateObj.State.Backend).To(Or(Equal("local"), Equal("s3")))
				Expect(configStateObj.State.Options.StateFile).To(Equal("devstream.state"))
			})
			Specify("Tools Options cannot be empty", func() {
				Expect(len(configStateObj.Tools)).ShouldNot(BeNil())
				Expect(configStateObj.Tools[0].Name).ShouldNot(BeEmpty())
				Expect(configStateObj.Tools[0].InstanceID).ShouldNot(BeEmpty())
				Expect(configStateObj.Tools[0].Options).ShouldNot(BeEmpty())
				Expect(len(configStateObj.Tools[0].DependsOn)).ShouldNot(BeNil())
			})
		})
		Context("when the Yaml parses fails", func() {
			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
