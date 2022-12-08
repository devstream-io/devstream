package configmanager

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config struct", func() {
	var c *Config

	BeforeEach(func() {
		c = &Config{}
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
