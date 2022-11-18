package configmanager

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Config struct", func() {
	var c *Config
	Context("validate method", func() {
		When("when state is nil", func() {
			BeforeEach(func() {
				c = &Config{State: nil}
			})
			It("should return err", func() {
				e := c.validate()
				Expect(e).Error().Should(HaveOccurred())
			})
		})
	})
})

var _ = Describe("newConfigRaw func", func() {
	var (
		fLoc    string
		baseDir string
	)
	BeforeEach(func() {
		baseDir = GinkgoT().TempDir()
		f, err := os.CreateTemp(baseDir, "test")
		Expect(err).Error().ShouldNot(HaveOccurred())
		fLoc = f.Name()
	})
	When("file not exist", func() {
		BeforeEach(func() {
			fLoc = "not_exist"
		})
		It("should return err", func() {
			_, err := newConfigRaw(fLoc)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("no such file or directory"))
		})
	})
	When("file content is not valid yaml", func() {
		BeforeEach(func() {
			err := os.WriteFile(fLoc, []byte("not_Valid_Yaml{{}}"), 0666)
			Expect(err).Error().ShouldNot(HaveOccurred())
		})
		It("should return err", func() {
			_, err := newConfigRaw(fLoc)
			Expect(err).Error().Should(HaveOccurred())
			Expect(err.Error()).Should(ContainSubstring("cannot unmarshal"))
		})
	})
})

var _ = Describe("ConfigRaw struct", func() {
	var (
		r          *ConfigRaw
		baseDir    string
		fLoc       string
		globalVars map[string]any
	)
	BeforeEach(func() {
		r = &ConfigRaw{}
		baseDir = GinkgoT().TempDir()
		f, err := os.CreateTemp(baseDir, "test")
		Expect(err).Error().ShouldNot(HaveOccurred())
		fLoc = f.Name()
		globalVars = map[string]any{}
	})
	Context("getGlobalVars method", func() {
		When("varFile get content failed", func() {
			BeforeEach(func() {
				r.VarFile = "not_exist"
			})
			It("should return err", func() {
				_, err := r.getGlobalVars()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("not exists"))
			})
		})
		When("varFiles is not valid", func() {
			BeforeEach(func() {
				r.VarFile = configFileLoc(fLoc)
				err := os.WriteFile(fLoc, []byte("not_Valid_Yaml{{}}"), 0666)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
			It("should return err", func() {
				_, err := r.getGlobalVars()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("cannot unmarshal"))
			})
		})
	})

	Context("getApps method", func() {
		When("appFiles get content failed", func() {
			BeforeEach(func() {
				r.AppFile = "not_exist"
			})
			It("should return err", func() {
				_, err := r.getAppsTools(globalVars)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("not exists"))
			})
		})
		When("appFiles is not valid", func() {
			BeforeEach(func() {
				err := os.WriteFile(fLoc, []byte("not_Valid_Yaml{{}}"), 0666)
				Expect(err).Error().ShouldNot(HaveOccurred())
				r.AppFile = configFileLoc(fLoc)
			})
			It("should return err", func() {
				_, err := r.getAppsTools(globalVars)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("yaml parse path[$.apps[*]] failed"))
			})
		})
	})

	Context("getTools method", func() {
		When("toolsFile get content failed", func() {
			BeforeEach(func() {
				r.ToolFile = "not_exist"
			})
			It("should return err", func() {
				_, err := r.getTools(globalVars)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("not exists"))
			})
		})
		When("toolFile is not valid", func() {
			BeforeEach(func() {
				err := os.WriteFile(fLoc, []byte("not_Valid_Yaml{{}}"), 0666)
				Expect(err).Error().ShouldNot(HaveOccurred())
				r.ToolFile = configFileLoc(fLoc)
			})
			It("should return err", func() {
				_, err := r.getTools(globalVars)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("yaml parse path[$.tools[*]] failed"))
			})
		})
		When("render failed", func() {
			BeforeEach(func() {
				r.ToolFile = ""
				r.globalBytes = []byte(`
tools:
  - name: plugin1
    instanceID: default
    options:
      key1: [[ var1 ]]`)
			})
			It("should return err", func() {
				_, err := r.getTools(globalVars)
				Expect(err).Error().Should(HaveOccurred())
			})
		})
		When("yaml render failed", func() {
			BeforeEach(func() {
				r.ToolFile = ""
				r.globalBytes = []byte(`
tools:
  - name: plugin1
    instanceID: default
    options:
      key1: {{}}`)
			})
			It("should return err", func() {
				_, err := r.getTools(globalVars)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("unexpected mapping key"))
			})
		})
		When("tool validate failed", func() {
			BeforeEach(func() {
				r.ToolFile = ""
				r.globalBytes = []byte(`
tools:
- name: plugin1
  instanceID: default
  dependsOn: [ "not_exist" ]
  options:
    key1: [[ var1 ]]`)
				globalVars = map[string]any{
					"var1": "global",
				}
			})
			It("should return err", func() {
				_, err := r.getTools(globalVars)
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("tool default's dependency not_exist doesn't exist"))
			})
		})

	})

	Context("getPipelineTemplatesMap method", func() {
		When("templateFile get content failed", func() {
			BeforeEach(func() {
				r.TemplateFile = "not_exist"
			})
			It("should return err", func() {
				_, err := r.getPipelineTemplatesMap()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("not exists"))
			})
		})
		When("getMergedNodeConfig failed", func() {
			BeforeEach(func() {
				r.TemplateFile = ""
				r.globalBytes = []byte(`

pipelineTemplates:
  - name: ci-pipeline-1
    type: githubactions
    options:
      branch: main
      docker:
        registry:
          type: dockerhub
          username: {{}}
          repository: [[ app ]]`)
			})
			It("should return err", func() {
				_, err := r.getPipelineTemplatesMap()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("unexpected mapping key"))
			})
		})
	})
})

var _ = Describe("getMergedNodeConfig func", func() {
	When("all input is null", func() {
		It("should return nil", func() {
			c, a, err := getMergedNodeConfig(nil, nil, "$.test")
			Expect(err).Error().ShouldNot(HaveOccurred())
			Expect(c).Should(Equal(""))
			Expect(a).Should(BeNil())
		})
	})
})
