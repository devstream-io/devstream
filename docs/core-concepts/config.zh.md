# 配置(Config)

DevStream使用 YAML 文件来声明 DevOps 工具链的配置。

## 1 配置内容概览

正如概述中所提到的，配置包含了以下几个部分：

- `config`
- `tools`
- `apps`
- `pipelineTemplates`
- `vars`

其中，`config` 必选，`tools` 和 `apps` 至少有一个不为空，其余部分可选。

## 2 组织方式

DevStream 支持两种组织配置的方式：

- 单文件：你可以把这些部分写到一个 YAML 文件中
- 目录：也可以把它们分开放到同一个文件夹下的多个 YAML 文件中，只要这些文件的名字以 `.yaml` 或 `.yml` 结尾，且内容拼接后包含了 DevStream 所需的各个部分即可。

然后在执行 `init` 等命令时，加上 `-f` 或 `--config-file` 参数指定配置文件/目录的路径。

如：

- 单文件：`dtm init -f config.yaml`
- 目录：`dtm init -f dirname`

## 3 配置内容详解

### 3.1 主配置

指 DevStream 本身的配置，即 `config` 部分，比如状态存储的方式等。详见 [这里](./state.zh.md)

### 3.2 工具的配置

`tools` 部分声明了工具链中的工具，详见 [这里](./tools.zh.md)

### 3.2 应用与流水线模板的配置

`apps` 部分声明了 应用 的配置，`pipelineTemplates` 声明了 流水线模板，详见 [这里](./apps.zh.md)

### 3.3 变量

DevStream 提供了变量语法。使用 `vars` 用来定义变量的值，而后可以在 `tools`、`apps`、`pipelineTemplates` 中使用，语法是 `[[ varName ]]`。

示例：

```yaml
vars:
  githubUsername: daniel-hutao # 定义变量
  repoName: dtm-test-go
  defaultBranch: main

tools:
- name: repo-scaffolding
  instanceID: default
  options:
    destinationRepo:
      owner: [[ githubUsername ]] # 使用变量
      name: [[ repoName ]]
      branch: [[ defaultBranch ]]
      scmType: github
    # <后面略...>
```

## 4 拓展语法

### 4.1 环境变量

类似于"变量"，你可以使用 `[[env "env_key"]]` 的方式来引用环境变量。

### 4.2 输出(Output)

#### 介绍

在 DevStream 的配置文件中，我们在配置 工具 的 Options 时，可以使用其他 工具 的 输出 来填充。

例如，如果 工具 A 有一个输出，我们可以将这个输出值作为 工具 B 的 Options。

注意：

- 当前，若 B 使用了 A 的输出，并不意味着 B "依赖于" A
- 如果 B 确实需要 "依赖于" A，即，我们想要保证在 B 运行之前行运行 A，我们仍然需要使用 `dependsOn` 关键字（详见上一节 "[核心概念](overview.zh.md)")。

#### 语法

我们可以通过以下格式来使用插件输出：

```
${{ TOOL_NAME.TOOL_INSTANCE_ID.outputs.OUTPUT_KEY }}
```

例如，对于下面给定的配置：

```yaml
tools:
- name: trello
  instanceID: default
  options:
    owner: IronCore864
    repo: golang-demo
    kanbanBoardName: golang-demo-board
```

- `TOOL_NAME` 是 "trello"
- `TOOL_INSTANCE_ID` 是 "default"

如果 "trello" 这个 工具 有一个键为 "boardId" 的输出项，那么我们就能通过以下语法来引用对应的输出的值：

```
${{ trello.default.outputs.boardId }}
```

#### 例子——真实使用场景

配置如下：

```yaml hl_lines="2 3 20 31"
tools:
- name: repo-scaffolding
  instanceID: golang-github
  options:
    destinationRepo:
      owner: IronCore864
      name: golang-demo
      branch: main
      scmType: github
    vars:
      imageRepo: "ironcore864/golang-demo"
    sourceRepo:
      org: devstream-io
      name: dtm-scaffolding-golang
      scmType: github
- name: helm-installer
  instanceID: argocd
- name: argocdapp
  instanceID: default
  dependsOn: [ "helm-installer.argocd", "repo-scaffolding.golang-github" ]
  options:
    app:
      name: golang-demo
      namespace: argocd
    destination:
      server: https://kubernetes.default.svc
      namespace: default
    source:
      valuefile: values.yaml
      path: helm/golang-demo
      repoURL: ${{ repo-scaffolding.golang-github.outputs.repoURL }}
```

在这个例子中，

- `argocdapp` 的 "default" 实例依赖于 `repo-scaffolding` 的 "golang-github" 实例
- `argocdapp` 的 "default" 实例中有一个 options 是 "options.source.repoURL"，它引用了 `repo-scaffolding` 的 "golang-github" 实例的 "repoURL" 输出(`${{ repo-scaffolding.golang-github.outputs.repoURL }}`)。

