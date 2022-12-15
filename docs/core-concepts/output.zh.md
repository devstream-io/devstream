# 输出(Output)

## 介绍

在 DevStream 的配置文件中，我们在配置 工具 的 Options 时，可以使用其他 工具 的 输出 来填充。

例如，如果 工具 A 有一个输出，我们可以将这个输出值作为 工具 B 的 Options。

注意：

- 当前，若 B 使用了 A 的输出，并不意味着 B "依赖于" A
- 如果 B 确实需要 "依赖于" A，即，我们想要保证在 B 运行之前行运行 A，我们仍然需要使用 `dependsOn` 关键字（详见上一节 "[核心概念](overview.zh.md)")。

## 语法

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

## 例子——真实使用场景

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
