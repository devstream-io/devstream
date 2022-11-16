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
      org: devstream-io # owner, org 二选一
      name: # 默认跟name一致
      url: github.com/devstream-io/service-A # 可选，如果有url，就不需要scm/owner/org/name了    # 如果 repoTemplate 为空，则默认 repo 已经建立好了，不需要再创建
      apiURL: gitlab.com/some/path/to/your/api # 可选，如果需要从模板创建repo
    # 如果 repoTemplate 不为空，则根据 scafoldingRepo 模板创建 repo
    repoTemplate: # 可选
      scmType: github
      owner: devstream-io
      org: devstream-io # owner, org 二选一
      name: dtm-scaffolding-golang
      url: github.com/devstream-io/dtm-scaffolding-golang # 可选，如果有url，就不需要scm/owner/org/name了    # 如果 repoTemplate 为空，则默认 repo 已经建立好了，不需要再创建
    ci:
      - type: template
        templateName: ci-pipeline-1
        options: # 覆盖模板中的options
          docker:
            registry:
              type: "[[ var3 ]]" # 在覆盖的同时，使用全局变量
        vars: # 可选, 传给模板用的变量（当且仅当type为template时有效）
          dockerUser: dockerUser1
          app: service-A
    cd:
      - type: cd-pipeline-custom
        options: # 自定义的options
          key: value
      - type: template
        templateName: cd-pipeline-1
        options: # 覆盖模板中的options
          destination:
            namespace: devstream-io
        vars: # 可选, 传给模板用的变量（当且仅当type为template时有效）
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
    type: githubactions # 对应一个插件
    options:
      branch: main # 可选 默认值
      docker:
        registry:
          type: dockerhub
          username: "[[ dockerUser ]]"
          repository: "[[ app ]]"
  - name: cd-pipeline-1
    type: argocdapp
    options:
      app:
        namespace: "[[ argocdNamespace ]]" # templates 内也可以引用全局变量
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
				RepoCommon: &RepoCommon{
					ScmType: "github",
					Owner:   "devstream-io",
					Org:     "devstream-io",
					Name:    "service-A",
					URL:     "github.com/devstream-io/service-A",
				},
				ApiURL: "gitlab.com/some/path/to/your/api",
			},
			RepoTemplate: &RepoTemplate{
				RepoCommon: &RepoCommon{
					ScmType: "github",
					Owner:   "devstream-io",
					Org:     "devstream-io",
					Name:    "dtm-scaffolding-golang",
					URL:     "github.com/devstream-io/dtm-scaffolding-golang",
				},
			},
			CIs: []PipelineTemplate{
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
			CDs: []PipelineTemplate{
				{
					Name: "cd-pipeline-custom",
					Type: "cd-pipeline-custom",
					Options: RawOptions{
						"key": "value",
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
