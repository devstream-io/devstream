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
				Expect(err.Error()).Should(ContainSubstring("field app.repo is required"))
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
				Expect(err.Error()).Should(ContainSubstring("field app.name is required"))
				Expect(err.Error()).Should(ContainSubstring("field app.spec is required"))
			})
		})
		When("app repo template is empty", func() {
			BeforeEach(func() {
				a = &app{
					Name: appName,
					Repo: &git.RepoInfo{
						CloneURL: "http://test.com/test/test_app",
					},
					Spec: &appSpec{Language: "go", FrameWork: "gin"},
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
				Expect(err.Error()).Should(ContainSubstring("field pipelineTemplate.type must be one of"))
			})
		})
	})
})
