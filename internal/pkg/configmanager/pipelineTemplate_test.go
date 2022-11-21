package configmanager

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm"
)

var _ = Describe("pipelineRaw struct", func() {
	var (
		r                                *pipelineRaw
		s                                *scm.SCMInfo
		templateMap                      map[string]string
		opt, globalVars                  map[string]any
		typeInfo, templateName, cloneURL string
	)
	BeforeEach(func() {
		cloneURL = "http://test.com"
		typeInfo = "github-actions"
		templateName = "testTemplate"
		templateMap = map[string]string{}
		globalVars = map[string]any{}
		opt = map[string]any{}
		s = &scm.SCMInfo{
			CloneURL: cloneURL,
		}
	})
	Context("newPipeline method", func() {
		When("type is not template", func() {
			BeforeEach(func() {
				opt = RawOptions{
					"toolconfig": "here",
				}
				r = &pipelineRaw{
					Type:    typeInfo,
					Options: opt,
				}
			})
			It("should return template", func() {
				expectedInfo := RawOptions{
					"pipeline": RawOptions(opt),
					"scm": RawOptions{
						"url": cloneURL,
					},
				}
				t, err := r.newPipeline(s, templateMap, globalVars)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(t.Type).Should(Equal(typeInfo))
				Expect(t.Name).Should(Equal(typeInfo))
				Expect(t.Options).Should(Equal(expectedInfo))
			})
		})
		When("templateName is empty", func() {
			BeforeEach(func() {
				r = &pipelineRaw{
					Type:         "template",
					TemplateName: "",
				}
			})
			It("should return err", func() {
				_, err := r.newPipeline(s, templateMap, globalVars)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("templateName is required"))
			})
		})
		When("template not exist in templateMap", func() {
			BeforeEach(func() {
				r = &pipelineRaw{
					Type:         "template",
					TemplateName: "not_exit",
				}
				templateMap = map[string]string{}
			})
			It("should return err", func() {
				_, err := r.newPipeline(s, templateMap, globalVars)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("not found in pipelineTemplates"))
			})
		})
		When("render template failed", func() {
			BeforeEach(func() {
				r = &pipelineRaw{
					Type:         "template",
					TemplateName: templateName,
				}
				templateMap = map[string]string{
					templateName: "[[ not_exist ]]",
				}
			})
			It("should return err", func() {
				_, err := r.newPipeline(s, templateMap, globalVars)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("render pipelineTemplate failed"))
			})
		})
		When("template is not valid yaml format", func() {
			BeforeEach(func() {
				r = &pipelineRaw{
					Type:         "template",
					TemplateName: templateName,
				}
				templateMap = map[string]string{
					templateName: "test{}{}",
				}
			})
			It("should return err", func() {
				_, err := r.newPipeline(s, templateMap, globalVars)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("parse pipelineTemplate yaml failed"))
			})
		})
		When("app and global has same value", func() {
			BeforeEach(func() {
				r = &pipelineRaw{
					Type:         "template",
					TemplateName: templateName,
					Vars: map[string]any{
						"var1": "cover",
					},
				}
				templateMap = map[string]string{
					templateName: fmt.Sprintf(`
name: %s
type: github-actions
options:
  branch: main
  docker:
  registry:
    username: [[ var1 ]]`, templateName),
				}
				globalVars = map[string]any{
					"var1": "global",
				}
			})
			It("should render with app vars", func() {
				t, err := r.newPipeline(s, templateMap, globalVars)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(t.Options).Should(Equal(RawOptions{
					"pipeline": RawOptions{
						"branch": "main",
						"registry": RawOptions{
							"username": "cover",
						},
						"configLocation": validPipelineConfigMap["github-actions"],
					},
					"scm": RawOptions{
						"url": "http://test.com",
					},
				}))
			})
		})
		When("app and template has options", func() {
			BeforeEach(func() {
				r = &pipelineRaw{
					Type:         "template",
					TemplateName: templateName,
					Options: RawOptions{
						"app":    "test",
						"option": "app",
					},
					Vars: map[string]any{
						"var1": "cover",
					},
				}
				templateMap = map[string]string{
					templateName: fmt.Sprintf(`
name: %s
type: githubactions
options:
  branch: main
  app: template
  registry:
    username: [[ var1 ]]`, templateName),
				}
				globalVars = map[string]any{
					"var1": "global",
				}
			})
			It("should render with app vars", func() {
				t, err := r.newPipeline(s, templateMap, globalVars)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(t.Options).Should(Equal(RawOptions{
					"pipeline": RawOptions{
						"option": "app",
						"branch": "main",
						"registry": RawOptions{
							"username": "cover",
						},
						"app":            "test",
						"configLocation": "",
					},
					"scm": RawOptions{
						"url": "http://test.com",
					},
				}))
			})
		})
	})
})

var _ = Describe("PipelineTemplate struct", func() {
	var (
		t       *pipelineTemplate
		opts    map[string]any
		appName string
	)
	BeforeEach(func() {
		t = &pipelineTemplate{}
	})
	Context("getPipelineTool method", func() {
		When("pipeline type is not valid", func() {
			BeforeEach(func() {
				t.Type = "not_exist"
			})
			It("should return err", func() {
				_, err := t.getPipelineTool("test")
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("pipeline type not_exist not supported for now"))
			})
		})
		When("pipeline type is valid", func() {
			BeforeEach(func() {
				appName = "test_app"
				opts = map[string]any{
					"test": "testV",
				}
				t.Type = "github-actions"
				t.Options = opts
			})
			It("should return tool", func() {
				tool, err := t.getPipelineTool(appName)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(tool).Should(Equal(&Tool{
					Name:       t.Type,
					InstanceID: appName,
					Options:    opts,
				}))
			})
		})
	})
})
