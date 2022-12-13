package configmanager

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("app struct", func() {
	var (
		a           *app
		appName     string
		vars        map[string]any
		templateMap map[string]string
	)
	BeforeEach(func() {
		appName = "test_app"
		vars = map[string]any{}
		templateMap = map[string]string{}
	})
	Context("getTools method", func() {
		When("repo is not valid", func() {
			BeforeEach(func() {
				a = &app{Name: appName}
			})
			It("should return error", func() {
				_, err := a.getTools(vars, templateMap)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("configmanager[app] is invalid, repo field must be configured"))
			})
		})
		When("ci/cd template is not valid", func() {
			BeforeEach(func() {
				a = &app{
					Repo: &git.RepoInfo{
						CloneURL: "http://test.com/test/test_app",
					},
					CIRawConfigs: []pipelineRaw{
						{
							Type:         "template",
							TemplateName: "not_exist",
						},
					},
				}
			})
			It("should return error", func() {
				_, err := a.getTools(vars, templateMap)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("not found in pipelineTemplates"))
			})
		})
		When("app repo template is empty", func() {
			BeforeEach(func() {
				a = &app{
					Name: appName,
					Repo: &git.RepoInfo{
						CloneURL: "http://test.com/test/test_app",
					},
				}
			})
			It("should return empty tools", func() {
				tools, err := a.getTools(vars, templateMap)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(len(tools)).Should(Equal(0))
			})
		})
	})

	Context("generateCICDTools method", func() {
		When("template type not exist", func() {
			BeforeEach(func() {
				a = &app{
					Repo: &git.RepoInfo{
						CloneURL: "http://test.com/test/test_app",
					},
					CIRawConfigs: []pipelineRaw{
						{
							Type:         "template",
							TemplateName: "not_valid",
						},
					},
				}
			})
			It("should return error", func() {
				templateMap = map[string]string{
					"not_valid": `
name: not_valid,
type: not_valid`}
				_, err := a.generateCICDTools(templateMap, vars)
				Expect(err).Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("pipeline type [not_valid] not supported for now"))
			})
		})
	})

	Context("generateLifecycleTool method", func() {
		When("lifeCycle config is empty", func() {
			BeforeEach(func() {
				a = &app{}
			})
			It("should return nil", func() {
				tool, err := a.generateLifecycleTool()
				Expect(err).Should(BeNil())
				Expect(tool).Should(BeNil())
			})
		})
		When("lifeCycle config is not valid", func() {
			BeforeEach(func() {
				a = &app{
					LifecycleManagement: &lifecycleManagement{
						Type: "not_exist",
					},
				}
			})
			It("should return nil", func() {
				_, err := a.generateLifecycleTool()
				Expect(err).ShouldNot(BeNil())
				Expect(err.Error()).Should(Equal("configmanager[app] lifecycle management only support jira"))
			})
		})
		When("lifecycle is valid", func() {
			BeforeEach(func() {
				a = &app{
					LifecycleManagement: &lifecycleManagement{
						Type: "jira-integ",
						Options: RawOptions{
							"test": "true",
						},
					},
					Repo: &git.RepoInfo{
						Repo: "test",
					},
				}
			})
			It("should return tool", func() {
				tool, err := a.generateLifecycleTool()
				Expect(err).Should(BeNil())
				Expect(tool).ShouldNot(BeNil())
				integOption, exist := tool.Options["integOptions"]
				Expect(exist).Should(BeTrue())
				Expect(integOption).Should(Equal(RawOptions{"test": "true"}))
				scmOptions, exist := tool.Options["scm"]
				Expect(exist).Should(BeTrue())
				Expect(scmOptions).Should(Equal(RawOptions{"name": "test"}))
			})
		})
	})
})
