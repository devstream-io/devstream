package configmanager

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("newRawConfigFromConfigText func", func() {
	var testText []byte
	const (
		validConfigTextWithoutKey = `
  state:
    backend: local
    options:
      stateFile: devstream.state`
		validVarsTextWithoutKey = `
  tests: argocd`
		validAppsTextWithoutKey = `
- name: service-a
  spec:
    language: python
    framework: [[ var1 ]]/[[ var2 ]]_gg`
	)
	When("text is valid", func() {
		BeforeEach(func() {
			testText = []byte(fmt.Sprintf(`---
config:%s
vars:%s
apps:%s`, validConfigTextWithoutKey, validVarsTextWithoutKey, validAppsTextWithoutKey))
		})
		It("should return config map", func() {
			rawMap, err := newRawConfigFromConfigBytes(testText)
			Expect(err).ShouldNot(HaveOccurred())
			removeLineFeed := func(s string) string {
				return strings.ReplaceAll(s, "\n", "")
			}
			Expect(removeLineFeed(string(rawMap.apps))).Should(Equal(removeLineFeed(validAppsTextWithoutKey)))
			Expect(removeLineFeed(string(rawMap.config))).Should(Equal(removeLineFeed(validConfigTextWithoutKey)))
			Expect(removeLineFeed(string(rawMap.vars))).Should(Equal(removeLineFeed(validVarsTextWithoutKey)))
		})
	})
	When("text is invalid", func() {
		When("invalid keys are used", func() {
			BeforeEach(func() {
				testText = []byte(fmt.Sprintf(`---
config:
%s
vars:
%s
app:
%s`, validConfigTextWithoutKey, validVarsTextWithoutKey, validAppsTextWithoutKey))
			})
			It("should return error", func() {
				_, err := newRawConfigFromConfigBytes(testText)
				Expect(err).Should(HaveOccurred())
			})
		})

		When("there are no enough keys", func() {
			BeforeEach(func() {
				testText = []byte(fmt.Sprintf(`
config:
%s
`, validConfigTextWithoutKey))
			})
			It("should return error", func() {
				_, err := newRawConfigFromConfigBytes(testText)
				Expect(err).Error().Should(HaveOccurred())
			})
		})
	})

})

var _ = Describe("rawConfig struct", func() {
	var r *rawConfig
	Context("checkValid method", func() {
		When("config is not exist", func() {
			BeforeEach(func() {
				r = &rawConfig{
					apps: []byte("apps"),
				}
			})
			It("should return error", func() {
				e := r.validate()
				Expect(e).Error().Should(HaveOccurred())
				Expect(e.Error()).Should(Equal("config not valid; check the [config] section of your config file"))
			})
		})
		When("apps and tools is not exist", func() {
			BeforeEach(func() {
				r = &rawConfig{
					config: []byte("config"),
				}
			})
			It("should return error", func() {
				e := r.validate()
				Expect(e).Error().Should(HaveOccurred())
				Expect(e.Error()).Should(Equal("config not valid; check the [tools and apps] section of your config file"))
			})
		})
	})

	Context("getTemplateMap method", func() {
		BeforeEach(func() {
			r = &rawConfig{
				pipelineTemplates: []byte(`---
- name: ci-pipeline-for-gh-actions
  type: github-actions # corresponding to a plugin
  options:
    docker:
      registry:
        username: [[ dockerUser ]]
        repository: [[ app ]]
- name: cd-pipeline-for-argocdapp
  type: argocdapp
  options:
    app:
      namespace: [[ argocdNamespace ]]/[[ test ]] # you can use global vars in templates
`)}
		})
		It("should get template", func() {
			m, err := r.getTemplatePipelineMap()
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(m).Should(Equal(map[string]string{
				"ci-pipeline-for-gh-actions": "name: ci-pipeline-for-gh-actions\ntype: github-actions\noptions:\n    docker:\n        registry:\n            repository: '[[ app ]]'\n            username: '[[ dockerUser ]]'\n",
				"cd-pipeline-for-argocdapp":  "name: cd-pipeline-for-argocdapp\ntype: argocdapp\noptions:\n    app:\n        namespace: '[[ argocdNamespace ]]/[[ test ]]'\n",
			}))
		})
	})

	Context("getApps method", func() {
		When("app file is not valid", func() {
			BeforeEach(func() {
				r = &rawConfig{
					apps: []byte(`
---
- name: app-1
  cd:
  - type: template
    vars:
      app: [[ appName ]]`)}
			})
			It("should get apps", func() {
				vars := map[string]any{}
				_, err := r.getAppsWithVars(vars)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("has no entry"))
			})
		})
	})
	Context("getConfig method", func() {
		BeforeEach(func() {
			r = &rawConfig{
				config: []byte(`
  state:
    backend: local
    options:
      stateFile: devstream.state`)}
		})

		When("get core config from config file", func() {
			It("should works fine", func() {
				cc, err := r.getConfig()
				Expect(err).NotTo(HaveOccurred())
				Expect(cc).NotTo(BeNil())
				Expect(cc.State.Backend).To(Equal("local"))
				Expect(cc.State.Options.StateFile).To(Equal("devstream.state"))
			})
		})
	})

	Context("getTools method", func() {
		BeforeEach(func() {
			r = &rawConfig{
				tools: []byte(`
- name: name1
  instanceID: ins-1
  dependsOn: []
  options:
    foo: [[ foo ]]`)}
		})

		When("get tools from config file", func() {
			It("should return config with vars", func() {
				tools, err := r.getToolsWithVars(map[string]any{"foo": interface{}("bar")})
				Expect(err).NotTo(HaveOccurred())
				Expect(tools).NotTo(BeNil())
				Expect(len(tools)).To(Equal(1))
				Expect(len(tools[0].Options)).To(Equal(1))
				Expect(tools[0].Options["foo"]).To(Equal(interface{}("bar")))
			})
		})
	})
	Context("getPipelineTemplateMap method", func() {
		BeforeEach(func() {
			r = &rawConfig{
				pipelineTemplates: []byte(`
- name: tpl1
  type: github-actions
  options:
    foo: bar`)}
		})
		When("get tools from config file", func() {
			It("should return config with vars", func() {
				pipelineTemplatesMap, err := r.getTemplatePipelineMap()
				Expect(err).NotTo(HaveOccurred())
				Expect(pipelineTemplatesMap).NotTo(BeNil())
				Expect(len(pipelineTemplatesMap)).To(Equal(1))
				Expect(pipelineTemplatesMap["tpl1"]).To(Equal("name: tpl1\ntype: github-actions\noptions:\n    foo: bar\n"))
			})
		})
	})

	Context("getVars method", func() {
		BeforeEach(func() {
			r = &rawConfig{
				vars: []byte(`---
foo1: bar1
foo2: 123
foo3: bar3`)}
		})
		It("should works fine", func() {
			varMap, err := r.getVars()
			Expect(err).NotTo(HaveOccurred())
			Expect(varMap).NotTo(BeNil())
			Expect(len(varMap)).To(Equal(3))
			Expect(varMap["foo1"]).To(Equal(interface{}("bar1")))
			Expect(varMap["foo2"]).To(Equal(interface{}(123)))
		})
	})
})
