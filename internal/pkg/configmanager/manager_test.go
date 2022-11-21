package configmanager_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
)

var _ = Describe("LoadConfig", func() {
	const (
		mainConfigFile = "config.yaml"
		varFile        = "var.yaml"
		toolFile       = "tool.yaml"
		appFile        = "app.yaml"
		templateFile   = "template.yaml"
		pluginDir      = "./plugins"
	)

	var tmpDir string
	var mainConfig string
	const appConfig = `apps:
  - name: service-A
    spec:
      language: python
      framework: django
    repo:
      scmType: github
      owner: devstream-io
      org: devstream-io # choose between owner and org
      url: github.com/devstream-io/service-A # optional，if url is specified，we can infer scm/owner/org/name from url
      apiURL: gitlab.com/some/path/to/your/api # optional, if you want to create a repo from repo template
    # if repoTemplate is not empty，we could help user to create repo from scaffoldingRepo
    repoTemplate: # optional
      scmType: github
      owner: devstream-io
      org: devstream-io # choose between owner and org
      name: dtm-scaffolding-golang
      url: github.com/devstream-io/dtm-scaffolding-golang # optional，if url is specified，we can infer scm/owner/org/name from url
    ci:
      - type: template
        templateName: ci-pipeline-1
        options: # overwrite options in pipelineTemplates
          docker:
            registry:
              type: [[ var3 ]] # while overridden, use global variables
        vars: # optional, use to render vars in template（valid only if the ci.type is template）
          dockerUser: dockerUser1
          app: service-A
    cd:
      - type: template
        templateName: cd-pipeline-1
        options: # overwrite options in pipelineTemplates
          destination:
            namespace: devstream-io
        vars: # optional, use to render vars in template（valid only if the cd.type is template）
          app: service-A
`
	const toolConfig = `tools:
  - name: plugin1
    instanceID: default
    options:
      key1: [[ var1 ]]
  - name: plugin2
    instanceID: ins2
    options:
      key1: value1
      key2: [[ var2 ]]
`
	const varConfig = `var2: value-of-var2
var3: dockerhub-overwrite
argocdNamespace: argocd
`
	const templateConfig = `pipelineTemplates:
  - name: ci-pipeline-1
    type: github-actions # corresponding to a plugin
    options:
      branch: main # optional, default is main
      docker:
        registry:
          type: dockerhub
          username: [[ dockerUser ]]
          repository: [[ app ]]
  - name: cd-pipeline-1
    type: argocdapp
    options:
      app:
        namespace: [[ argocdNamespace ]] # you can use global vars in templates
      destination:
        server: https://kubernetes.default.svc
        namespace: default
      source:
        valuefile: values.yaml
        path: helm/[[ app ]]
        repoURL: ${{repo-scaffolding.myapp.outputs.repoURL}}
`
	var (
		state = &configmanager.State{
			Backend: "local",
			Options: configmanager.StateConfigOptions{
				StateFile: "devstream.state",
			},
		}

		expectedConfig = &configmanager.Config{
			PluginDir: pluginDir,
			State:     state,
		}
	)
	BeforeEach(func() {
		// create files
		tmpDir = GinkgoT().TempDir()
		filesMap := map[string]string{
			varFile:      varConfig,
			toolFile:     toolConfig,
			appFile:      appConfig,
			templateFile: templateConfig,
		}
		for file, content := range filesMap {
			err := os.WriteFile(filepath.Join(tmpDir, file), []byte(content), 0644)
			Expect(err).Should(Succeed())
		}
	})
	When("only with config file", func() {
		BeforeEach(func() {
			mainConfig = fmt.Sprintf(`# main config
varFile: %s
toolFile: %s
appFile: %s
templateFile: %s
pluginDir: %s
state: # state config, backend can be local or s3
  backend: local
  options:
    stateFile: devstream.state

var1: value-of-var1 # var1 is a global var

`, varFile, toolFile, appFile, templateFile, pluginDir)
			err := os.WriteFile(filepath.Join(tmpDir, mainConfigFile), []byte(mainConfig), 0644)
			Expect(err).Should(Succeed())
		})
		It("should load config correctly", func() {
			expectedTools1 := configmanager.Tool{
				Name:       "plugin1",
				InstanceID: "default",
				DependsOn:  nil,
				Options: configmanager.RawOptions{
					"key1": "value-of-var1",
				},
			}
			expectedTools2 := configmanager.Tool{
				Name:       "plugin2",
				InstanceID: "ins2",
				DependsOn:  nil,
				Options: configmanager.RawOptions{
					"key1": "value1",
					"key2": "value-of-var2",
				},
			}
			expectedTools3 := configmanager.Tool{
				Name:       "github-actions",
				InstanceID: "service-A",
				DependsOn: []string{
					"repo-scaffolding.service-A",
				},
				Options: configmanager.RawOptions{
					"pipeline": configmanager.RawOptions{
						"docker": configmanager.RawOptions{
							"registry": configmanager.RawOptions{
								"repository": "service-A",
								"type":       "dockerhub-overwrite",
								"username":   "dockerUser1",
							},
						},
						"branch":         "main",
						"configLocation": "git@github.com:devstream-io/ci-template.git//github-actions",
					},
					"scm": configmanager.RawOptions{
						"apiURL":  "gitlab.com/some/path/to/your/api",
						"owner":   "devstream-io",
						"org":     "devstream-io",
						"name":    "service-A",
						"scmType": "github",
						"url":     "https://github.com/devstream-io/service-A",
					},
					"instanceID": "service-A",
				},
			}
			expectedTools4 := configmanager.Tool{
				Name:       "argocdapp",
				InstanceID: "service-A",
				DependsOn: []string{
					"repo-scaffolding.service-A",
				},
				Options: configmanager.RawOptions{
					"pipeline": configmanager.RawOptions{
						"destination": configmanager.RawOptions{
							"namespace": "devstream-io",
							"server":    "https://kubernetes.default.svc",
						},
						"app": configmanager.RawOptions{
							"namespace": "argocd",
						},
						"source": configmanager.RawOptions{
							"valuefile": "values.yaml",
							"path":      "helm/service-A",
							"repoURL":   "${{repo-scaffolding.myapp.outputs.repoURL}}",
						},
						"configLocation": "",
					},
					"scm": configmanager.RawOptions{
						"url":     "https://github.com/devstream-io/service-A",
						"apiURL":  "gitlab.com/some/path/to/your/api",
						"owner":   "devstream-io",
						"org":     "devstream-io",
						"name":    "service-A",
						"scmType": "github",
					},
					"instanceID": "service-A",
				},
			}
			expectedTools5 := configmanager.Tool{
				Name:       "repo-scaffolding",
				InstanceID: "service-A",
				DependsOn:  nil,
				Options: configmanager.RawOptions{
					"destinationRepo": configmanager.RawOptions{
						"needAuth": true,
						"org":      "devstream-io",
						"repo":     "service-A",
						"branch":   "main",
						"repoType": "github",
						"url":      "github.com/devstream-io/service-A",
					},
					"sourceRepo": configmanager.RawOptions{
						"repoType": "github",
						"url":      "github.com/devstream-io/dtm-scaffolding-golang",
						"needAuth": true,
						"org":      "devstream-io",
						"repo":     "dtm-scaffolding-golang",
						"branch":   "main",
					},
					"vars": configmanager.RawOptions{
						"var2":            "value-of-var2",
						"framework":       "django",
						"language":        "python",
						"var3":            "dockerhub-overwrite",
						"argocdNamespace": "argocd",
						"var1":            "value-of-var1",
					},
					"instanceID": "service-A",
				},
			}

			manager := configmanager.NewManager(filepath.Join(tmpDir, mainConfigFile))
			config, err := manager.LoadConfig()
			Expect(err).Should(Succeed())
			Expect(config).ShouldNot(BeNil())
			Expect(config.PluginDir).Should(Equal(expectedConfig.PluginDir))
			Expect(config.State).Should(Equal(expectedConfig.State))
			Expect(len(config.Tools)).Should(Equal(5))
			Expect(config.Tools[0]).Should(Equal(expectedTools1))
			Expect(config.Tools[1]).Should(Equal(expectedTools2))
			Expect(config.Tools[2]).Should(Equal(expectedTools3))
			Expect(config.Tools[3]).Should(Equal(expectedTools4))
			Expect(config.Tools[4]).Should(Equal(expectedTools5))
		})
	})
	When("global file and configFile all has config", func() {
		BeforeEach(func() {
			mainConfig = fmt.Sprintf(`# main config
toolFile: %s
pluginDir: %s
var1: test
var2: test2
state: # state config, backend can be local or s3
  backend: local
  options:
    stateFile: devstream.state
tools:
  - name: plugin1
    instanceID: default
    options:
      config: app
  - name: plugin3
    instanceID: ins2
    options:
      key1: value1
      key2: [[ var2 ]]
`, toolFile, pluginDir)
			err := os.WriteFile(filepath.Join(tmpDir, mainConfigFile), []byte(mainConfig), 0644)
			Expect(err).Should(Succeed())
		})
		It("should merge config correctly", func() {
			manager := configmanager.NewManager(filepath.Join(tmpDir, mainConfigFile))
			config, err := manager.LoadConfig()
			Expect(err).Should(Succeed())
			Expect(config).ShouldNot(BeNil())
			Expect(config.PluginDir).Should(Equal(expectedConfig.PluginDir))
			Expect(config.State).Should(Equal(expectedConfig.State))
			Expect(len(config.Tools)).Should(Equal(4))
		})
	})
	When("with global config", func() {
		BeforeEach(func() {
			mainConfig = `
---
varFile:
toolFile:
pluginDir: ""
state:
  backend: local
  options:
    stateFile: devstream.state

---
# variables config
defaultBranch: main
githubUsername: testUser
repoName: dtm-test-go
jiraID: merico
jiraUserEmail: test@email
jiraProjectKey: DT
dockerhubUsername: exploitht
argocdNameSpace: test
argocdDeployTimeout: 10m

---
# plugins config
tools:
  - name: repo-scaffolding
    instanceID: golang-github
    options:
      destinationRepo:
        owner: [[ githubUsername ]]
        org: ""
        repo: [[ repoName ]]
        branch: [[ defaultBranch ]]
        repoType: github
      sourceRepo:
        org: devstream-io
        repo: dtm-scaffolding-golang
        repoType: github
      vars:
        ImageRepo: "[[ dockerhubUsername ]]/[[ repoName ]]"
  - name: jira-github-integ
    instanceID: default
    dependsOn: [ "repo-scaffolding.golang-github" ]
    options:
      owner: [[ githubUsername ]]
      repo: [[ repoName ]]
      jiraBaseUrl: https://[[ jiraID ]].atlassian.net
      jiraUserEmail: [[ jiraUserEmail ]]
      jiraProjectKey: [[ jiraProjectKey ]]
      branch: main
  - name: githubactions-golang
    instanceID: default
    dependsOn: [ "repo-scaffolding.golang-github" ]
    options:
      owner: ${{repo-scaffolding.golang-github.outputs.owner}}
      org: ""
      repo: ${{repo-scaffolding.golang-github.outputs.repo}}
      language:
        name: go
        version: "1.18"
      branch: [[ defaultBranch ]]
      build:
        enable: True
        command: "go build ./..."
      test:
        enable: True
        command: "go test ./..."
        coverage:
          enable: True
          profile: "-race -covermode=atomic"
          output: "coverage.out"
      docker:
        enable: True
        registry:
          type: dockerhub
          username: [[ dockerhubUsername ]]
          repository: ${{repo-scaffolding.golang-github.outputs.repo}}`
			err := os.WriteFile(filepath.Join(tmpDir, mainConfigFile), []byte(mainConfig), 0644)
			Expect(err).Should(Succeed())
		})
		It("should merge config correctly", func() {
			manager := configmanager.NewManager(filepath.Join(tmpDir, mainConfigFile))
			config, err := manager.LoadConfig()
			Expect(err).Should(Succeed())
			Expect(config).ShouldNot(BeNil())
			Expect(config.PluginDir).Should(Equal(""))
			Expect(len(config.Tools)).Should(Equal(3))
		})
	})

})

var _ = Describe("Manager struct", func() {
	var (
		m       *configmanager.Manager
		fLoc    string
		baseDir string
	)
	BeforeEach(func() {
		m = &configmanager.Manager{}
		baseDir = GinkgoT().TempDir()
		f, err := os.CreateTemp(baseDir, "test")
		Expect(err).Error().ShouldNot(HaveOccurred())
		fLoc = f.Name()
		err = os.WriteFile(fLoc, []byte(`
tools:
  - name: plugin1
    instanceID: default
    options:
      key1: test
		`), 0666)
		Expect(err).Error().ShouldNot(HaveOccurred())
		m.ConfigFile = fLoc
	})

	Context("LoadConfig method", func() {
		When("get RawConfig failed", func() {
			BeforeEach(func() {
				m.ConfigFile = "not_exist"
			})
			It("should return error", func() {
				_, err := m.LoadConfig()
				Expect(err).Error().Should(HaveOccurred())
				Expect(err.Error()).Should(ContainSubstring("no such file or directory"))
			})
		})
		When("getGlobalVars failed", func() {
			BeforeEach(func() {
				err := os.WriteFile(fLoc, []byte(`
varFile: not_exist
state:
  backend: local
  options:
    stateFile: devstream.state`), 0666)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
			It("should return error", func() {
				_, err := m.LoadConfig()
				Expect(err).Error().Should(HaveOccurred())
			})
		})
		When("getTools failed", func() {
			BeforeEach(func() {
				err := os.WriteFile(fLoc, []byte(`
toolFile: not_exist
state:
  backend: local
  options:
    stateFile: devstream.state`), 0666)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
			It("should return error", func() {
				_, err := m.LoadConfig()
				Expect(err).Error().Should(HaveOccurred())
			})
		})
		When("getApps failed", func() {
			BeforeEach(func() {
				err := os.WriteFile(fLoc, []byte(`
appFile: not_exist
state:
  backend: local
  options:
    stateFile: devstream.state`), 0666)
				Expect(err).Error().ShouldNot(HaveOccurred())
			})
			It("should return error", func() {
				_, err := m.LoadConfig()
				Expect(err).Error().Should(HaveOccurred())
			})

		})
	})
})
