# Output

## Introduction

In DevStream's configuration file, we can use _Output_ from one _Tool_ as the options values for another _Tool_.

For example, if _Tool_ A has an output, we can use its value as _Tool_ B's options.

Notes:

- At the moment, B using A's output doesn't mean B "depends on" A.
- If B really needs to "depend on" A, i.e., we want to make sure A runs first before B runs, we still need to use the `dependsOn` keyword (see the previous section "[Core Concepts](core-concepts)" for more details.)

## Syntax

To use the output, follow this format:

```
${{ TOOL_NAME.TOOL_INSTANCE_ID.outputs.OUTPUT_KEY }}
```

For example, given config:

```yaml
tools:
- name: trello
  instanceID: default
  options:
    owner: IronCore864
    repo: golang-demo
    kanbanBoardName: golang-demo-board
```

- TOOL_NAME is "trello"
- TOOL_INSTANCE_ID is "default"

If the "trello" tool/plugin has an output key name "boardId", then we can use its value by the following syntax:

```
${{ trello.default.outputs.boardId }}
```

## Real-World Usage Example

Config:

```yaml
---
tools:
- name: github-repo-scaffolding-golang
  instanceID: default
  options:
    owner: IronCore864
    repo: golang-demo
    branch: main
    image_repo: ironcore864/golang-demo
- name: argocd
  instanceID: default
  options:
    create_namespace: true
    repo:
      name: argo
      url: https://argoproj.github.io/argo-helm
    chart:
      chart_name: argo/argo-cd
      release_name: argocd
      namespace: argocd
      wait: true
      timeout: 10m
      upgradeCRDs: true
- name: argocdapp
  instanceID: default
  dependsOn: [ "argocd.default", "github-repo-scaffolding-golang.default" ]
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
      repoURL: ${{ github-repo-scaffolding-golang.default.outputs.repoURL }} # pay attention here
```

In this example:
- The "default" instance of tool `argocdapp` depends on the "default" instance of tool `github-repo-scaffolding-golang` 
- The "default" instance of tool `argocdapp` has an user option "options.source.repoURL", which uses the "default" instance of tool `github-repo-scaffolding-golang`'s output "repoURL" (`${{ github-repo-scaffolding-golang.default.outputs.repoURL }}`)
