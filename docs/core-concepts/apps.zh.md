# 应用(Apps)

## 1 概念

应用在 DevStream 中表示对一个服务的生命周期配置，包括对服务的项目脚手架，CI 流程以及 CD 流程的配置，在 devstream 中使用应用可以使用少数几行配置就构建出服务的整条 CI/CD 流水线配置。

### 1.1 应用(Apps)

有时候，你需要为一个应用/微服务定义多个 工具。例如，对于一个 web 应用程序，你可能需要指定以下工具：

- 仓库脚手架
- 持续集成
- 持续部署

如果你有多个应用需要管理，你需要在配置中创建多个 工具，这可能会很繁琐，而且配置文件难以阅读。

为了更容易地管理多个应用/微服务，DevStream 提供了另一层抽象，称为 应用。你可以在一个应用中定义所有内容（例如，上面提到的仓库脚手架、持续集成、持续部署等），只需要几行配置，就可以使配置文件更容易阅读和管理。

在底层，DevStream 仍然会将你的 应用 配置转换为 工具 的定义，但你不需要关心它。

## 1.2 流水线模板（pipelineTemplates）

一个 DevStream 应用 可以引用一个或多个 流水线模板，这些模板主要是 CI/CD 定义，这样 应用 的定义可以更短，以在多个微服务之间共享公共的 CI/CD 流水线。

## 2 配置方式

### 2.1 应用

应用在配置中由 `apps` 来定义，它是一个数组，每个元素的字段如下：

- name: 应用的名称，必须是唯一的
- spec: 配置应用特定的信息
- repo: 代码仓库信息
- repoTemplate: 可选。结构和 repo 一致。若非空，DevStream 将使用项目脚手架来创建仓库。
- ci: 可选，一个 CI 流水线的数组，每个元素包含以下字段：
    - type: 类型。值为 `template` 或某个插件的名称。
    - templateName: 可选。当 type 为 `template` 时，指定使用的模板名称。
    - vars: 可选。传入模板的变量，仅当 type 为 `template` 时有效。
    - options: 可选。
        - 当 type 是某个插件名时，就是该插件的 Options
        - 当 type 为 `template` 时，覆盖模板中的 Options。你需要写清想要覆盖的选项的完整路径，例如 `options.docker.registry.type`。
- cd: 与 `ci` 类似，是 CD 流水线列表。dtm 会先执行 ci 流水线，再执行 cd 流水线。

### 2.2 流水线模板

流水线模板在配置中由 `pipelineTemplates` 来定义，它是一个数组，每个元素的字段如下：

- name: 流水线模板的名称，必须是唯一的
- type: 对应插件的名称
- options: 插件的 Options （可能和原始的插件 Options 不同）

### 2.3 局部变量

DevStream 已经有了全局变量，所有的 工具 和 应用 都可以直接引用。

但有时，我们想多次引用某个 DevOps 工具，但他们只有一些细微的差别。例如，除了项目名称不同，其他的都相同。

在这种情况下，我们可以定义一个流水线模板，在其中定义一个局部变量，然后在 应用 引用该流水线模板时，传入不同的值。

```yaml hl_lines="13 15 23 30" title="流水线模板的使用与局部变量"
apps:
- name: my-app
  spec:
    language: java
    framework: springboot
    repo: 
      url: https://github.com/testUser/testApp.git
      branch: main
    ci:
    - type: github-actions # 不定义 pipelineTemplates，直接使用插件
    cd:
    - type: template # 使用流水线模板
      templateName: my-cd-template # 对应 pipelineTemplates 中的 name
      vars:
        appName: my-app # 传入模板的变量

pipelineTemplates:
cd:
- name: my-cd-template
  type: argocdapp
  options:
    app:
      name: [[ appName ]] # 定义一个局部变量，在引用使用模板时传入
      namespace: argocd # argocd 的命名空间
    destination:
      server: https://kubernetes.default.svc # 部署的 kubernetes 服务地址
      namespace: default # 应用要部署的命名空间
    source:
      valuefile: values.yaml # 项目中的 helm 变量文件名
      path: charts/[[ appName ]] # 项目中的 helm 配置路径
```

## 3 完整配置示例

应用 的真实示例配置如下：

```yaml
apps:
- name: testApp # 应用名称
  spec: # 该配置项用于配置应用特定的信息
    language: java #应用所使用的编程语言
    framework: springboot #应用所使用的编程框架
  repo: # 该配置项用于应用的代码仓库信息
    url: https://github.com/testUser/testApp.git
    branch: main
  repoTemplate: # 可选，用于创建应用的脚手架。不为空时，会使用该脚手架创建应用的代码仓库
    url: https://github.com/devstream-io/dtm-repo-scaffolding-java-springboot.git
    vars:
      imageRepoOwner: repoOwner # 用于渲染脚手架模版的变量
  ci: # 配置应用的 ci 流程，如下即为使用 github-actions 运行应用 ci 流程
  - type: github-actions
- name: testApp2
  spec:
    language: go
    framework: gin
  repo: # 该配置项用于应用的代码仓库信息
    owner: test_user
    type: github
    branch: main
  repoTemplate: # 该配置用于创建应用的脚手架
    org: devstream-io
    name: dtm-repo-scaffolding-java-springboot
    type: github
  ci: # 配置应用的 ci 流程，如下即为使用 github-actions 运行应用 ci 流程
  - type: github-actions
    options:
      imageRepo:
        owner: repoOwner # 覆盖 插件/模板 原有的 Options，YAML 路径要写完全
  cd: # 配置应用的 cd，如果为使用 argocdapp 运行应用的 cd 流程
  - type: argocdapp
```

使用该配置就会在 `gitlab` 仓库中创建两个应用，项目的脚手架均为 DevStream 官方提供的 [SpringBoot](https://github.com/devstream-io/dtm-repo-scaffolding-java-springboot.git) 项目。应用 `testApp` 会在每次代码提交后使用 `github-actions` 运行测试。应用 `testApp2` 会在每次代码提交后使用 `github-actions` 运行测试并构建镜像推送到 `repoOwner` 的镜像仓库中，最后使用 Argo CD 将应用部署到集群中。

### repo/repoTemplate 配置

应用配置中的 `repo` 和 `repoTemplate` 均表示一个代码仓库，支持通过 url 或 详细字段来配置：

!!! note "两种配置代码仓库的方式"

    === "使用 url 来配置"
    
        ```yaml title=""
          repo:
            url: git@gitlab.example.com:root/myapps.git # url 表示仓库的地址，支持 git 地址和 http 地址
            apiURL: https://gitlab.example.com # 非必填，如果使用 gitlab 而且 url 使用的是 git 地址，则需要配置该字段用于表示 gitlab 的 api 请求地址
            branch: "" # 非必填，github 默认为 main 分支，gitlab 默认为 master 分支
        ```
    
        该配置表示使用的仓库为 `gitlab`, 要使用 `git@gitlab.example.com:root/myapps.git` 来克隆代码，devstream 会使用使用 `https://gitlab.example.com` 和 gitlab 进行交互，仓库的主分支为 `master` 分支。
    
    === "使用仓库的详细字段来配置"
    
        ```yaml title=""
          repo:
            org: "" # 非必填，仓库的拥有组织名称，如果使用 github 的组织则需要填写此字段
            owner："test_user" # 如果仓库是非组织的，则需要填写该字段表示仓库拥有者
            name: "" # 非必填，默认为应用名
            baseURL:  https://gitlab.example.com # 非必填，如果使用 gitlab 则需要填写该字段表示 gitlab 的域名
            branch: master #非必填，github 默认为 main 分支，gitlab 默认为 master 分支
            type:  gitlab #必填，表示该仓库的类型，目前支持 gitlab/github
        ```
    
        该配置表示代码仓库使用的是 `gitlab` 且其地址为 `https://gitlab.example.com`，仓库名称为应用名，仓库的所有者为 `test_user`，主分支为 `master` 分支。

### ci 配置

应用配置中的 `ci` 目前支持 `github-actions`/`gitlab-ci`/`jenkins-pipeline`/`template` 4 种类型。

其中 `template` 类型表示使用 流水线模板 来运行流水线，前 3 种类型分别对应了 `github` 中的 actions 流水线，`gitlab` 中 ci 流水线和 `jenkins` 中的 pipeline，它们的具体配置如下：

```yaml
  ci:
  - type: jenkins-pipieline # 表明当前 ci 的类型
    options: # ci 的具体配置项，如果该配置项为空，则 ci 只会运行单元测试然后结束
      jenkins: # 该配置项之用于 jenkins，表示 jenkins 的一些配置信息
        url: jenkins.exmaple.com # jenkins 的地址
        user: admin # jenkins 用户
      imageRepo: # 需要推送的镜像仓库信息，如果设置了该字段，则 ci 流程会在测试成功后构建镜像推送到该镜像仓库
        url: http://harbor.example.com # 镜像仓库地址，若为空则默认为 dockerhub
        owner: admin # 镜像仓库拥有者名称
      dingTalk: # 钉钉通知的配置信息，如果设置了该字段，则 ci 流程中会将最后的构建结构通过钉钉发送通知
        name: dingTalk
        webhook: https://oapi.dingtalk.com/robot/send?access_token=changemeByConfig # 钉钉的回调地址
        securityType: SECRET # 使用 secret 模式来加密钉钉的信息
        securityValue: SECRETDATA # 钉钉的 secret 加密字符串
      sonarqube: # sonarqube 的配置信息，如果设置了该字段，则 ci 流程会和测试并行执行 sonarqube 的代码扫描
        url: http://sonar.example.com # sonarqube 的地址
        token: YOUR_SONAR_TOKEN # soanrqube 的认证 token
        name: sonar_test
```

上述的配置即会在应用更新推送到代码仓库后先并行执行单元测试和代码扫描，然后构建镜像推送到代码仓库，流程都成功后发送通知消息到指定的钉钉群中。如果所有应用都要配置一遍该类型 `ci`, 那配置就会变得比较繁琐，所以 devstream 还提供了 `template` 用于在多个应用间共享 `ci` 的流程配置，具体如下所示：

```yaml
apps:
- name: javaProject1
  spec:
    language: java
    framework: springboot
  repo:
    owner: testUser
    type: github
  repoTemplate:
    url: https://github.com/devstream-io/dtm-repo-scaffolding-java-springboot.git
  ci:
  - type: template # 表示该 ci 流程使用模版
    templateName: ci-pipeline # ci 使用的模版名称
    vars:
      dingdingAccessToken: tokenForProject1 #用于渲染 ci 模版的变量
      dingdingSecretValue: secretValProject1
- name: javaProject2
  spec:
    language: java
    framework: springboot
  repo:
    owner: testUser
    type: github
  repoTemplate:
    url: https://github.com/devstream-io/dtm-repo-scaffolding-java-springboot.git
  ci:
  - type: template # 表示该 ci 流程使用模版
    templateName: ci-pipeline # ci 使用的模版名称
    vars:
      dingdingAccessToken: tokenForProject2 # 设置传入模版 "ci-pipeline" 的变量
      dingdingSecretValue: secretValProject2

pipelineTemplates: # 即配置的 ci/cd 模版
- name: ci-pipeline # 模版名称
  type: jenkins-pipeline #模版类型，支持 jenkins-pipeline，github-actions 和 gitlab-ci
  options: # options 和 ci 的 options 完全一致
    jenkins:
      url: jenkins.exmaple.com
      user: admin
    imageRepo:
      url: http://harbor.example.com
      owner: admin
    dingTalk:
      name: dingTalk
      webhook: https://oapi.dingtalk.com/robot/send?access_token=[[ dingdingAccessToken ]] # 用于被 app 渲染的模版，这样就可以实现不同应用使用同一个模版发送通知到不同的钉钉群
      securityType: SECRET
      securityValue: [[ dingdingSecretValue ]] # 局部变量，应用 在引用此模板时通过 vars 来传入
    sonarqube:
        url: http://sonar.example.com
        token: sonar_token
        name: sonar_test

```

使用以上的配置，就会创建两个和上述流程一致的 `jenkins` 流水线，不同之处只在于两个应用通知的钉钉群不一样。


### cd 配置

应用的 cd 配置目前只支持 `argocdapp`，可以使用 `argocd` 来将应用部署在集群中，具体配置如下：

```yaml
cd:
- type: argocdapp
  options:
    app:
      name: hello # argocd 应用名称
      namespace: argocd # argocd 的命名空间
    destination:
      server: https://kubernetes.default.svc # 部署的 Kubernetes 服务地址
      namespace: default # 应用要部署的命名空间
    source:
      valuefile: values.yaml # 项目中的 helm 变量文件名
      path: charts/go-hello-http # 项目中的 helm 配置路径
```
