package configmanager

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

const configFileStr = `---
config:
  state:
    backend: local
    options:
      stateFile: devstream.state

vars:
  foo1: bar1
  foo2: 123
  appName: service-a
  registryType: dockerhub
  argocdNamespace: argocd

apps:
- name: service-a
  spec:
    language: python
    framework: django
  repo:
    scmType: github
    owner: devstream-io
    org: devstream-io # choose between owner and org
    url: github.com/devstream-io/service-a # optional，if url is specified，we can infer scm/owner/org/name from url
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
    templateName: ci-pipeline-for-gh-actions
    options: # overwrite options in pipelineTemplates
      docker:
        registry:
          type: [[ registryType ]] # while overridden, use global variables
    vars: # optional, use to render vars in template（valid only if the ci.type is template）
      dockerUser: dockerUser1
      app: service-a
  cd:
  - type: template
    templateName: cd-pipeline-for-argocdapp
    options: # overwrite options in pipelineTemplates
      destination:
        namespace: devstream-io
    vars: # optional, use to render vars in template（valid only if the cd.type is template）
      app: [[ appName ]]

tools:
- name: plugin1
  instanceID: default
  dependsOn: []
  options:
    foo1: [[ foo1 ]]
- name: plugin2
  instanceID: tluafed
  dependsOn: []
  options:
    foo: bar
    foo2: [[ foo2 ]]

pipelineTemplates:
- name: ci-pipeline-for-gh-actions
  type: github-actions # corresponding to a plugin
  options:
    branch: main # optional, default is main
    docker:
      registry:
        type: dockerhub
        username: [[ dockerUser ]]
        repository: [[ app ]]
- name: cd-pipeline-for-argocdapp
  type: argocdapp
  options:
    app:
      namespace: [[ argocdNamespace ]] # you can use global vars in templates
    destination:
      server: https://kubernetes.default.svc
      namespace: devstream-io
    source:
      valuefile: values.yaml
      path: helm/[[ app ]]
      repoURL: ${{repo-scaffolding.myapp.outputs.repoURL}}
`

var tmpWorkDir string

var _ = Describe("LoadConfig", func() {

	tool1 := &Tool{
		Name:       "plugin1",
		InstanceID: "default",
		DependsOn:  []string{},
		Options: RawOptions{
			"foo1":       "bar1",
			"instanceID": "default",
		},
	}

	tool2 := &Tool{
		Name:       "plugin2",
		InstanceID: "tluafed",
		DependsOn:  []string{},
		Options: RawOptions{
			"instanceID": "tluafed",
			"foo":        "bar",
			"foo2":       123,
		},
	}

	tool3 := &Tool{
		Name:       "github-actions",
		InstanceID: "service-a",
		DependsOn: []string{
			"repo-scaffolding.service-a",
			"plugin1.default",
			"plugin2.tluafed",
		},
		Options: RawOptions{
			"instanceID": "service-a",
			"pipeline": RawOptions{
				"language": RawOptions{
					"name":      "python",
					"framework": "django",
				},
				"docker": RawOptions{
					"registry": RawOptions{
						"repository": "service-a",
						"type":       "dockerhub",
						"username":   "dockerUser1",
					},
				},
				"branch":         "main",
				"configLocation": "git@github.com:devstream-io/ci-template.git//github-actions",
			},
			"scm": RawOptions{
				"apiURL":  "gitlab.com/some/path/to/your/api",
				"owner":   "devstream-io",
				"org":     "devstream-io",
				"name":    "service-a",
				"scmType": "github",
				"url":     "https://github.com/devstream-io/service-a",
			},
		},
	}

	tool4 := &Tool{
		Name:       "argocdapp",
		InstanceID: "service-a",
		DependsOn: []string{
			"repo-scaffolding.service-a",
			"plugin1.default",
			"plugin2.tluafed",
		},
		Options: RawOptions{
			"instanceID": "service-a",
			"destination": RawOptions{
				"namespace": "devstream-io",
				"server":    "https://kubernetes.default.svc",
			},
			"app": RawOptions{
				"namespace": "argocd",
			},
			"source": RawOptions{
				"valuefile": "values.yaml",
				"path":      "helm/service-a",
				"repoURL":   "${{repo-scaffolding.myapp.outputs.repoURL}}",
			},
		},
	}

	tool5 := &Tool{
		Name:       "repo-scaffolding",
		InstanceID: "service-a",
		DependsOn: []string{
			"plugin1.default",
			"plugin2.tluafed",
		},
		Options: RawOptions{
			"instanceID": "service-a",
			"destinationRepo": RawOptions{
				"needAuth": true,
				"org":      "devstream-io",
				"repo":     "service-a",
				"branch":   "main",
				"repoType": "github",
				"url":      "github.com/devstream-io/service-a",
			},
			"sourceRepo": RawOptions{
				"repoType": "github",
				"url":      "github.com/devstream-io/dtm-scaffolding-golang",
				"needAuth": true,
				"org":      "devstream-io",
				"repo":     "dtm-scaffolding-golang",
				"branch":   "main",
			},
			"vars": RawOptions{},
		},
	}

	BeforeEach(func() {
		tmpWorkDir = GinkgoT().TempDir()
		err := os.WriteFile(filepath.Join(tmpWorkDir, "config.yaml"), []byte(configFileStr), 0644)
		Expect(err).NotTo(HaveOccurred())
	})

	When("load a config file", func() {
		It("should return 5 tools", func() {
			mgr := NewManager(filepath.Join(tmpWorkDir, "config.yaml"))
			cfg, err := mgr.LoadConfig()
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg).NotTo(BeNil())

			GinkgoWriter.Printf("Config: %v", cfg)

			// config/state
			Expect(*cfg.Config.State).To(Equal(State{
				Backend: "local",
				Options: StateConfigOptions{
					StateFile: "devstream.state",
				},
			}))

			// vars
			Expect(len(cfg.Vars)).To(Equal(5))
			Expect(cfg.Vars["foo1"]).To(Equal("bar1"))

			// tools
			Expect(len(cfg.Tools)).To(Equal(5))
			for _, t := range cfg.Tools {
				switch t.Name {
				case "plugin1":
					Expect(t).Should(Equal(tool1))
				case "plugin2":
					Expect(t).Should(Equal(tool2))
				case "github-actions":
					Expect(t).Should(Equal(tool3))
				case "argocdapp":
					Expect(t).Should(Equal(tool4))
				case "repo-scaffolding":
					Expect(t).Should(Equal(tool5))
				default:
					Fail("Unexpected plugin name.")
				}
			}
		})
	})
})

var _ = Describe("getConfigFromFileWithGlobalVars", func() {
	BeforeEach(func() {
		tmpWorkDir = GinkgoT().TempDir()
		err := os.WriteFile(filepath.Join(tmpWorkDir, "config.yaml"), []byte(configFileStr), 0644)
		Expect(err).NotTo(HaveOccurred())
	})

	When("get config from file", func() {
		It("should return a config", func() {
			mgr := NewManager(filepath.Join(tmpWorkDir, "config.yaml"))
			cfg, err := mgr.getConfigFromFileWithGlobalVars()
			Expect(err).NotTo(HaveOccurred())
			Expect(cfg.Config.State.Backend).To(Equal("local"))
			Expect(cfg.Vars["foo1"]).To(Equal("bar1"))
			Expect(len(cfg.Apps)).To(Equal(1))
			Expect(cfg.Apps[0].Name).To(Equal("service-a"))
			Expect(len(cfg.Tools)).To(Equal(2))
			Expect(cfg.Tools[1].Name).To(Equal("plugin2"))
		})
	})
})
