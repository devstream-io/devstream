# Output

## Introduction

In DevStream's configuration file, we can use _Output_ from one _Tool_ as the options values for another _Tool_.

For example, if _Tool_ A has an output, we can use its value as _Tool_ B's options.

Notes:

- At the moment, B using A's output doesn't mean B "depends on" A.
- If B really needs to "depend on" A, i.e., we want to make sure A runs first before B runs, we still need to use the `dependsOn` keyword (see the previous section "[Core Concepts](./core_concepts.md)" for more details.)

## Syntax

To use the output, follow this format:

```
${{ TOOL_NAME.PLUGIN.outputs.OUTPUT_KEY }}
```

For example, given config:

```yaml
tools:
- name: kanban
  plugin: trello
  options:
    owner: IronCore864
    repo: golang-demo
    kanbanBoardName: golang-demo-board
```

- TOOL_NAME is "kanban"
- PLUGIN is "trello"

If the "trello" plugin has an output key name "boardId", then we can use its value by the following syntax:

```
${{ kanban.trello.outputs.boardId }}
```

## Real-World Usage Example

Config:

```yaml
---
tools:
- name: repo
  plugin: github-repo-scaffolding-golang
  options:
    owner: IronCore864
    repo: golang-demo
    branch: main
    image_repo: ironcore864/golang-demo
- name: cd
  plugin: argocd
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
- name: demo
  plugin: argocdapp
  dependsOn: [ "cd.argocd", "demo.github-repo-scaffolding-golang" ]
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
      repoURL: ${{ demo.github-repo-scaffolding-golang.outputs.repoURL }} # pay attention here
```

In this example:
- Tool "demo" (plugin: argocdapp) depends on tool "repo" (plugin: github-repo-scaffolding-golang);
- tool "demo" has an user option "options.source.repoURL", which uses tool "repo" output "repoURL" (`${{ demo.github-repo-scaffolding-golang.outputs.repoURL }}`)


```{toctree}
---
maxdepth: 1
---
```