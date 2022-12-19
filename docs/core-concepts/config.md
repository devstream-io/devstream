# Config

DevStream uses YAML files to define your DevOps platform.

## 1 Sections Overview

As aforementioned in the overview, there are several sections in the config:

- `config`
- `tools`
- `apps`
- `pipelineTemplates`
- `vars`

Among which, `config` is mandatory, and you should have at least either `tools` or `apps`. Others are optional.

## 2 Config File V.S. Config Folder

DevStream supports both:

- a single config file: put all sections of the config into one YAML file
- a directory: put multiple files in one directory, as long as the file names end with `.yaml` or `.yml`

When you run the `init` (or other commands), add `-f` or `--config-file`.

Examples:

- single file: `dtm init -f config.yaml`
- directory: `dtm init -f dirname`

## 3 Sections Detail

### 3.1 The Main Config

DevStream's own config, the `config` section, contains where to store the state. See [here](./state.md) for more details.


### 3.2 Tools Config

The `tools` section defines DevOps tools. Read [this](./tools.md) for all the detail.

### 3.3 Apps/pipelineTemplates Config

The `apps` section defines Apps and the `pipelineTemplates` defines templates for pipelines. See [here](./apps.md) for more detail.

### 3.4 Variables

DevStream supports variables. Define key/values in the `vars` section and refer to it in the `tools`, `apps` and `pipelineTemplates` sections.

Use double square brackets when referring to a variable: `[[ varName ]]`.

example:

```yaml
vars:
  githubUsername: daniel-hutao # define a variable
  repoName: dtm-test-go
  defaultBranch: main

tools:
- name: repo-scaffolding
  instanceID: default
  options:
    destinationRepo:
      owner: [[ githubUsername ]] # refer to the pre-defined variable
      name: [[ repoName ]]
      branch: [[ defaultBranch ]]
      scmType: github
    # ...
```

## 4 Expanded Features

### 4.1 Environment Variables

Similar to "variables", you can use `[[env "env_key"]]` to refer to environment variables.

### 4.2 Output

#### Introduction

In DevStream's configuration file, we can use Output from one Tool as the options values for another Tool.

For example, if Tool A has an output, we can use its value as Tool B's options.

Notes:

- At the moment, B using A's output doesn't mean B "depends on" A.
- If B needs to "depend on" A, i.e., we want to make sure A runs first before B runs, we still need to use the `dependsOn` keyword (see the previous section "[Core Concepts](overview.md)" for more details.)

#### Syntax

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

- `TOOL_NAME` is "trello"
- `TOOL_INSTANCE_ID` is "default"

If the "trello" tool/plugin has an output key name "boardId", then we can use its value by the following syntax:

```
${{ trello.default.outputs.boardId }}
```

#### Real-World Usage Example

Config:

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
      repoURL: ${{ repo-scaffolding.golang-github.outputs.repoURL }} # pay attention here
```

In this example:

- The "default" instance of tool `argocdapp` depends on the "golang-github" instance of tool `repo-scaffolding`
- The "default" instance of tool `argocdapp` has a user option "options.source.repoURL", which uses the "golang-github" instance of tool `repo-scaffolding`'s output "repoURL" (`${{ repo-scaffolding.golang-github.outputs.repoURL }}`)

