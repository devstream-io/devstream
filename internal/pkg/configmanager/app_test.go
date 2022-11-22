package configmanager

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm"
)

var _ = Describe("getToolsFromApp func", func() {
	var (
		appStr      string
		globalVars  map[string]any
		templateMap map[string]string
	)
	BeforeEach(func() {
		templateMap = map[string]string{}
		globalVars = map[string]any{}
	})
	When("render app failed", func() {
		BeforeEach(func() {
			appStr = "[[ field ]]"
			globalVars = map[string]any{
				"not_exist": "failed",
			}
		})
		It("should return err", func() {
			_, err := getToolsFromApp(appStr, globalVars, templateMap)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("app render globalVars failed"))
		})
	})
	When("unmarshal appData to yaml failed", func() {
		BeforeEach(func() {
			appStr = "fsdfsd{123213}"
		})
		It("should return error", func() {
			_, err := getToolsFromApp(appStr, globalVars, templateMap)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("app parse yaml failed"))
		})
	})
	When("build ciTemplate failed", func() {
		BeforeEach(func() {
			appStr = `
name: test
repo:
  name: test
ci:
- type: template
  templateName:`
		})
		It("should return error", func() {
			_, err := getToolsFromApp(appStr, globalVars, templateMap)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("app[test] get pipeline tools failed"))
		})
	})
	When("repo not exist", func() {
		BeforeEach(func() {
			appStr = `
name: test
cd:
- type: template
  templateName:`
		})
		It("should return error", func() {
			_, err := getToolsFromApp(appStr, globalVars, templateMap)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("app.repo field can't be empty"))
		})
	})
	When("build cdTemplate failed", func() {
		BeforeEach(func() {
			appStr = `
name: test
repo:
  name: test
cd:
- type: template
  templateName:`
		})
		It("should return error", func() {
			_, err := getToolsFromApp(appStr, globalVars, templateMap)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("app[test] get pipeline tools failed"))
		})
	})
	When("app data is valid", func() {
		BeforeEach(func() {
			appStr = `
name: service-A
spec:
  language: python
  framework: django
repo:
  url: github.com/devstream-io/service-A
repoTemplate: # optional
  scmType: github
  owner: devstream-io
  org: devstream-io
  name: dtm-scaffolding-golang
ci:
  - type: argocdapp
    options:
          key: test
cd:
  - type: template
    templateName: testTemplate`
			globalVars = map[string]any{
				"var1": "test",
			}
			templateMap = map[string]string{
				"testTemplate": `
type: github-actions
name: testTemplate
options:
  keyFromApp: test`,
			}
		})
		It("should return tools", func() {
			expectA := Tool{
				Name:       "argocdapp",
				InstanceID: "service-A",
				DependsOn: []string{
					"repo-scaffolding.service-A",
				},
				Options: RawOptions{
					"pipeline": RawOptions{
						"key":            "test",
						"configLocation": validPipelineConfigMap["argocdapp"],
					},
					"scm": RawOptions{
						"name": "service-A",
						"url":  "https://github.com/devstream-io/service-A",
					},
					"instanceID": "service-A",
				},
			}
			expectG := Tool{
				Name:       "github-actions",
				InstanceID: "service-A",
				DependsOn: []string{
					"repo-scaffolding.service-A",
				},
				Options: RawOptions{
					"pipeline": RawOptions{
						"keyFromApp":     "test",
						"configLocation": validPipelineConfigMap["github-actions"],
					},
					"scm": RawOptions{
						"url":  "https://github.com/devstream-io/service-A",
						"name": "service-A",
					},
					"instanceID": "service-A",
				},
			}
			expectR := Tool{
				Name:       "repo-scaffolding",
				InstanceID: "service-A",
				DependsOn:  nil,
				Options: RawOptions{
					"destinationRepo": RawOptions{
						"needAuth": true,
						"owner":    "devstream-io",
						"repo":     "service-A",
						"branch":   "main",
						"repoType": "github",
						"url":      "github.com/devstream-io/service-A",
					},
					"sourceRepo": RawOptions{
						"repoType": "github",
						"needAuth": true,
						"owner":    "devstream-io",
						"org":      "devstream-io",
						"repo":     "dtm-scaffolding-golang",
						"branch":   "main",
					},
					"vars": RawOptions{
						"var1":      "test",
						"language":  "python",
						"framework": "django",
					},
					"instanceID": "service-A",
				},
			}
			tools, err := getToolsFromApp(appStr, globalVars, templateMap)
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(tools).ShouldNot(BeNil())
			Expect(len(tools)).Should(Equal(3))
			Expect(tools[0]).Should(Equal(expectA))
			Expect(tools[1]).Should(Equal(expectG))
			Expect(tools[2]).Should(Equal(expectR))
		})
	})
})

var _ = Describe("rawApp struct", func() {
	var (
		a           *rawApp
		appName     string
		rawConfig   []pipelineRaw
		templateMap map[string]string
		globalVars  map[string]any
	)
	BeforeEach(func() {
		templateMap = map[string]string{}
		globalVars = map[string]any{}
	})
	Context("setDefault method", func() {
		When("repoInfo is not config", func() {
			BeforeEach(func() {
				appName = "test"
				a = &rawApp{
					Repo: &scm.SCMInfo{},
					Name: appName,
				}
			})
			It("should update repoinfo by appName", func() {
				a.setDefault()
				Expect(a.Repo.Name).Should(Equal("test"))
			})
		})
	})
	Context("generateCICDToolsFromAppConfig method", func() {
		When("config is not valid", func() {
			BeforeEach(func() {
				rawConfig = []pipelineRaw{
					{
						Type:         "template",
						TemplateName: "",
					},
				}
				a.CIRawConfigs = rawConfig
			})
			It("should return error", func() {
				_, err := a.generateCICDToolsFromAppConfig(templateMap, globalVars)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring(".templateName is required"))
			})
		})
		When("all valid", func() {
			BeforeEach(func() {
				rawConfig = []pipelineRaw{
					{
						Type: "github-actions",
					},
					{
						Type: "jenkins-pipeline",
					},
				}
				a.CIRawConfigs = rawConfig
			})
			It("should return pipeline array", func() {
				p, err := a.generateCICDToolsFromAppConfig(templateMap, globalVars)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(len(p)).Should(Equal(2))
				Expect(p[0].Name).Should(Equal("github-actions"))
				Expect(p[1].Name).Should(Equal("jenkins-pipeline"))
			})
		})
	})
	Context("getRepoTemplateTool method", func() {
		When("repo is null", func() {
			BeforeEach(func() {
				a.Repo = nil
			})
			It("should return err", func() {
				_, err := a.getRepoTemplateTool(globalVars)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("app.repo field can't be empty"))
			})
			When("template is nil", func() {
				BeforeEach(func() {
					a.Repo = &scm.SCMInfo{
						Name: "test",
					}
					a.RepoTemplate = nil
				})
				It("should return nil", func() {
					t, err := a.getRepoTemplateTool(globalVars)
					Expect(err).Error().ShouldNot(HaveOccurred())
					Expect(t).Should(BeNil())
				})
			})
		})
	})
})
