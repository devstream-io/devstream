package configmanager

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kr/pretty"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("LoadConfig Most Complex", func() {
	const (
		mainConfigFile = "config.yaml"
		varFile        = "var.yaml"
		toolFile       = "tool.yaml"
		appFile        = "app.yaml"
		templateFile   = "template.yaml"
		pluginDir      = "./plugins"
	)

	var tmpDir string
	var mainConfig = fmt.Sprintf(`# main config
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
	const appConfig = `apps:
  - name: service-A
    spec:
      language: python
      framework: django
    repo:
      scmType: github
      owner: devstream-io
      org: devstream-io # choose between owner and org
      name: # optional, default is the same as app name
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
              type: "[[ var3 ]]" # while overridden, use global variables
        vars: # optional, use to render vars in template（valid only if the ci.type is template）
          dockerUser: dockerUser1
          app: service-A
    cd:
      - type: cd-pipeline-custom # if the type is not "template", it means plugins
        options: # options to the plugins
          key: "[[ var1 ]]" # use global vars
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
      key1: "[[ var1 ]]"
  - name: plugin2
    instanceID: ins2
    options:
      key1: value1
      key2: "[[ var2 ]]"
`
	const varConfig = `var2: value-of-var2
var3: dockerhub-overwrite
argocdNamespace: argocd
`
	const templateConfig = `pipelineTemplates:
  - name: ci-pipeline-1
    type: githubactions # corresponding to a plugin
    options:
      branch: main # optional, default is main
      docker:
        registry:
          type: dockerhub
          username: "[[ dockerUser ]]"
          repository: "[[ app ]]"
  - name: cd-pipeline-1
    type: argocdapp
    options:
      app:
        namespace: "[[ argocdNamespace ]]" # you can use global vars in templates
      destination:
        server: https://kubernetes.default.svc
        namespace: default
      source:
        valuefile: values.yaml
        path: "helm/[[ app ]]"
        repoURL: ${{repo-scaffolding.myapp.outputs.repoURL}}
`
	var (
		tools = Tools{
			{
				Name:       "plugin1",
				InstanceID: "default",
				DependsOn:  []string{},
				Options: RawOptions{
					"key1": "value-of-var1",
				},
			},
			{
				Name:       "plugin2",
				InstanceID: "ins2",
				DependsOn:  []string{},
				Options: RawOptions{
					"key1": "value1",
					"key2": "value-of-var2",
				},
			},
		}
		state = &State{
			Backend: "local",
			Options: StateConfigOptions{
				StateFile: "devstream.state",
			},
		}
		app = App{
			Name: "service-A",
			Spec: RawOptions{
				"language":  "python",
				"framework": "django",
			},
			Repo: &Repo{
				RepoInfo: &RepoInfo{
					ScmType: "github",
					Owner:   "devstream-io",
					Org:     "devstream-io",
					Name:    "service-A",
					URL:     "github.com/devstream-io/service-A",
				},
				ApiURL: "gitlab.com/some/path/to/your/api",
			},
			RepoTemplate: &RepoTemplate{
				RepoInfo: &RepoInfo{
					ScmType: "github",
					Owner:   "devstream-io",
					Org:     "devstream-io",
					Name:    "dtm-scaffolding-golang",
					URL:     "github.com/devstream-io/dtm-scaffolding-golang",
				},
			},
			CIPipelines: []PipelineTemplate{
				{
					Name: "ci-pipeline-1",
					Type: "githubactions",
					Options: RawOptions{
						"branch": "main",
						"docker": RawOptions{
							"registry": RawOptions{
								"type":       "dockerhub-overwrite",
								"username":   "dockerUser1",
								"repository": "service-A",
							},
						},
					},
				},
			},
			CDPipelines: []PipelineTemplate{
				{
					Name: "cd-pipeline-custom",
					Type: "cd-pipeline-custom",
					Options: RawOptions{
						"key": "value-of-var1",
					},
				},
				{
					Name: "cd-pipeline-1",
					Type: "argocdapp",
					Options: RawOptions{
						"app": RawOptions{
							"namespace": "argocd",
						},
						"destination": RawOptions{
							"server":    "https://kubernetes.default.svc",
							"namespace": "devstream-io",
						},
						"source": RawOptions{
							"valuefile": "values.yaml",
							"path":      "helm/service-A",
							"repoURL":   "${{repo-scaffolding.myapp.outputs.repoURL}}",
						},
					},
				},
			},
		}
		apps = Apps{app}

		expectedConfig = &Config{
			PluginDir: pluginDir,
			Tools:     tools,
			Apps:      apps,
			State:     state,
		}
	)
	BeforeEach(func() {
		// create files
		tmpDir = GinkgoT().TempDir()
		filesMap := map[string]string{
			mainConfigFile: mainConfig,
			varFile:        varConfig,
			toolFile:       toolConfig,
			appFile:        appConfig,
			templateFile:   templateConfig,
		}
		fmt.Println(tmpDir)
		for file, content := range filesMap {
			err := os.WriteFile(filepath.Join(tmpDir, file), []byte(content), 0644)
			Expect(err).Should(Succeed())
		}
	})
	It("should load config correctly", func() {
		manager := NewManager(filepath.Join(tmpDir, mainConfigFile))
		config, err := manager.LoadConfig()
		Expect(err).Should(Succeed())
		Expect(config).ShouldNot(BeNil())
		pretty.Println(config)

		Expect(config).Should(Equal(expectedConfig))
	})
})
