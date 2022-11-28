package configmanager

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("getCoreConfigFromConfigFile", func() {
	const configFile = `---
config:
  state:
    backend: local
    options:
      stateFile: devstream.state
`
	When("get core config from config file", func() {
		It("should works fine", func() {
			cc, err := getCoreConfigFromConfigFile([]byte(configFile))
			Expect(err).NotTo(HaveOccurred())
			Expect(cc).NotTo(BeNil())
			Expect(cc.State.Backend).To(Equal("local"))
			Expect(cc.State.Options.StateFile).To(Equal("devstream.state"))
		})
	})
})
