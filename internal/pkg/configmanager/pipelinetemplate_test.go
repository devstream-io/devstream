package configmanager

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/pkg/util/scm/git"
)

var _ = Describe("pipelineRaw struct", func() {
	var (
		r                      *pipelineRaw
		templateMap            map[string]string
		opt, globalVars        map[string]any
		typeInfo, templateName string
	)
	BeforeEach(func() {
		typeInfo = "github-actions"
		templateName = "testTemplate"
		templateMap = map[string]string{}
		globalVars = map[string]any{}
		opt = map[string]any{}
	})
	Context("getPipelineTemplate method", func() {
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
					"toolconfig": "here",
				}
				t, err := r.getPipelineTemplate(templateMap, globalVars)
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
				_, err := r.getPipelineTemplate(templateMap, globalVars)
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
				_, err := r.getPipelineTemplate(templateMap, globalVars)
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
				_, err := r.getPipelineTemplate(templateMap, globalVars)
				Expect(err).Error().Should(HaveOccurred())
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
				_, err := r.getPipelineTemplate(templateMap, globalVars)
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
				t, err := r.getPipelineTemplate(templateMap, globalVars)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(t.Options).Should(Equal(RawOptions{
					"branch": "main",
					"docker": nil,
					"registry": RawOptions{
						"username": "cover",
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
type: github-actions
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
				t, err := r.getPipelineTemplate(templateMap, globalVars)
				Expect(err).ShouldNot(HaveOccurred())
				Expect(t.Options).Should(Equal(RawOptions{
					"branch": "main",
					"registry": RawOptions{
						"username": "cover",
					},
					"app":    "test",
					"option": "app",
				}))
			})
		})
	})
})

var _ = Describe("PipelineTemplate struct", func() {
	var (
		t                 *pipelineTemplate
		s                 *git.RepoInfo
		opts              map[string]any
		appName, cloneURL string
		globalOption      *pipelineGlobalOption
	)
	BeforeEach(func() {
		appName = "test_app"
		cloneURL = "http://test.com"
		t = &pipelineTemplate{}
		s = &git.RepoInfo{
			CloneURL: git.ScmURL(cloneURL),
		}
		globalOption = &pipelineGlobalOption{
			Repo:    s,
			AppName: appName,
		}
	})
	Context("generatePipelineTool method", func() {
		When("pipeline type is not valid", func() {
			BeforeEach(func() {
				t.Type = "not_exist"
			})
			It("should return err", func() {
				_, err := t.generatePipelineTool(globalOption)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("pipeline type [not_exist] not supported for now"))
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
				tool, err := t.generatePipelineTool(globalOption)
				Expect(err).Error().ShouldNot(HaveOccurred())
				Expect(tool).Should(Equal(&Tool{
					Name:       t.Type,
					InstanceID: appName,
					DependsOn:  []string{},
					Options: RawOptions{
						"pipeline": RawOptions{
							"test":           "testV",
							"configLocation": "https://raw.githubusercontent.com/devstream-io/dtm-pipeline-templates/main/github-actions/workflows/main.yml",
						},
						"scm": RawOptions{
							"url": git.ScmURL(cloneURL),
						},
					},
				}))
			})
		})
	})

	Context("updatePipelineVars method", func() {
		var pipelineOpt *pipelineGlobalOption

		BeforeEach(func() {
			t.Options = RawOptions{
				"imageRepo": RawOptions{
					"user": "test_user",
				},
			}
			pipelineOpt = &pipelineGlobalOption{}
		})
		It("should update global option", func() {
			t.updatePipelineVars(pipelineOpt)
			Expect(pipelineOpt.ImageRepo).ShouldNot(BeNil())
			Expect(pipelineOpt.ImageRepo).Should(Equal(RawOptions{
				"user": "test_user",
			}))
		})
	})
})
