# Config

This document summarizes the config file for DevStream.

DevStream uses a single YAML file to store your DevOps toolchain configuration.

## Default Config File

By default, `dtm` tries to use `./config.yaml` (under your current directory.)

## Specifying a Config File Explicitly

You can override the default value with `-f` or `--config-file`. Examples:

```shell
dtm apply -f path/to/your/config.yaml
dtm apply --config-file path/to/your/config.yaml
```

## Config File Content

The config file only contains:

- Only one section (at the moment), which is `tools`.
- `tools` is a list of dictionaries.
- Each dictionary defines a DevOps "tool" which is managed by a DevStream plugin
- Each dictionary (tool) has the following mandatory fields:
    - `name`: the name of the tool, string, without underscore
    - `plugin`: the plugin to be used
    - you can have duplicated `name` in one config file, you can also have duplicated `plugin` in one config file, but the `name + plugin` combination must be unique in one config file
- Each dictionary (tool) has an optional field which is `options`, which in turn is a dictionary containing parameters for that specific plugin. For plugins' parameters, see [plugins](./plugins.md).
- Each directory (tool) has an optional field which is `dependsOn`. Continue reading for detail about dependencies.

## Example Config File

`config.yaml`:

```yaml
tools:
- name: go-webapp-repo
  plugin: github-repo-scaffolding-golang
  options:
    org: devstream-io
    repo: dtm-e2e-go
    branch: main
    image_repo: dtme2etest/dtm-e2e-go
```

## Dependencies

If you want tool A to be installed before tool B, you can let tool B depend on tool A.

The syntax for dependency is:
    
```yaml
dependsOn: ["NAME1.PLUGIN1"]'
```

Since `dependsOn` is a list, a tool can have multiple dependencies:

```
dependsOn: [ "NAME1.PLUGIN1", "NAME2.PLUGIN2", "..."]
```

In the following example, tool "go-webapp-repo" (using plugin github-repo-scaffolding-golang) will be installed before tool "golang-demo-actions" (using plugin githubactions-golang):

```yaml
tools:
- name: go-webapp-repo
  plugin: github-repo-scaffolding-golang
  options:
    org: devstream-io
    repo: dtm-e2e-go
    branch: main
    image_repo: dtme2etest/dtm-e2e-go
- name: golang-demo-actions
  plugin: githubactions-golang
  dependsOn: ["go-webapp-repo.github-repo-scaffolding-golang"]
  options:
    org: ${{go-webapp-repo.github-repo-scaffolding-golang.outputs.org}}
    repo: ${{go-webapp-repo.github-repo-scaffolding-golang.outputs.repo}}
    language:
      name: go
      version: "1.17"
    branch: main
    build:
      enable: True
    test:
      enable: True
      coverage:
        enable: True
```

## Variables

To not repeat yourself (DRY,) we can define some variables in a var file and use the vars in the config file.

### Default Variable File

The default var file is located at `./variables.yaml`. If this file doesn't exist, and no user-specified var file is provided, the config won't be rendered with any variables, apparently.

### Specifying a Variable File Explicitly

To override the default location of the variables file, use the `--var-file` option:

```shell
dtm apply -f path/to/your/config.yaml --var-file path/to/your/variables.yaml
```

### Variable File Content

The var file is a YAML file containing key-value pairs. Example:

"variables.yaml":

```yaml
gitlabUser: ironcore864
defaultBranch: main
gitlabCIGolangTemplate: https://gitlab.com/ironcore864/go-hello-world/-/raw/main/go.tpl
```

At the moment, nested/composite values (for example, the value is a list/dictionary) are not supported yet.

### Using a Variables File

To use these variables in a config file, use the following syntax:

```yaml
[[ varNameHere ]]
```

### Example Config File with the Use of Variables

`variables.yaml`:

```yaml
gitlabUser: ironcore864
defaultBranch: main
gitlabCIGolangTemplate: https://gitlab.com/ironcore864/go-hello-world/-/raw/main/go.tpl
```

Example config with the variables defined in the above `variables.yaml`:

`config.yaml`:

```yaml
tools:
- name: myapp
  plugin: gitlabci-generic
  options:
    pathWithNamespace: [[ gitlabUser ]]/go-hello-world
    branch: [[ defaultBranch ]]
    templateURL: [[ gitlabCIGolangTemplate ]]
    templateVariables:
      App: hello
```

DevStream will render the config with the provided var file. After rendering, the config above is equivalent to the following content:

```yaml
tools:
- name: myapp
  plugin: gitlabci-generic
  options:
    pathWithNamespace: ironcore864/go-hello-world
    branch: main
    templateURL: https://gitlab.com/ironcore864/go-hello-world/-/raw/main/go.tpl
    templateVariables:
      App: hello
```

```{toctree}
---
maxdepth: 1
---
```
