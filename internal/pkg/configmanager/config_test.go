package configmanager

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config struct", func() {
	var (
		c                    *Config
		toolName, instanceID string
	)
	BeforeEach(func() {
		c = &Config{}
		toolName = "test_tool"
		instanceID = "test_instance"
	})
	Context("renderInstanceIDtoOptions method", func() {
		When("tool option is null", func() {
			BeforeEach(func() {
				c.Tools = Tools{
					{Name: toolName, InstanceID: instanceID},
				}
			})
			It("should set nil to RawOptions", func() {
				c.renderInstanceIDtoOptions()
				Expect(len(c.Tools)).Should(Equal(1))
				tool := c.Tools[0]
				Expect(tool.Options).Should(Equal(RawOptions{
					"instanceID": instanceID,
				}))
			})
		})
	})
	Context("validate method", func() {
		When("config state is null", func() {
			BeforeEach(func() {
				c.Config.State = nil
			})
			It("should return err", func() {
				err := c.validate()
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(Equal("config.state is not defined"))
			})
		})
	})
})
