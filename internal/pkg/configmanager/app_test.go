package configmanager

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("getToolsFromConfigFileWithVarsRendered", func() {
	const appsConfig = `---
apps:
- name: app-1
  cd:
  - type: template
    vars:
      app: [[ appName ]]
`
	When("get apps from config file", func() {
		It("should return config with vars", func() {
			apps, err := getAppsFromConfigFileWithVarsRendered([]byte(appsConfig), map[string]any{"appName": interface{}("app-1")})
			Expect(err).NotTo(HaveOccurred())
			Expect(apps).NotTo(BeNil())
			Expect(len(apps)).To(Equal(1))
			Expect(len(apps[0].CDRawConfigs)).To(Equal(1))
			Expect(apps[0].CDRawConfigs[0].Vars["app"]).To(Equal(interface{}("app-1")))
		})
	})
})
