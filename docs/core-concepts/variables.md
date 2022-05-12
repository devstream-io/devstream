# Variables File

To not repeat yourself (DRY,) we can define some variables in a var file and use the vars in the config file.

You can specify the path/to/your/variablesFile.yaml in the main config. See the [config](./config.md) section for more details.

## Default Variable File

If `varFile` isn't specified explicitly in the main config, `variables.yaml` under the current directory will be used. If this file doesn't exist, and no user-specified var file is provided, the config won't be rendered with any variables.

## Using a Variables File

To use these variables in a config file, use the following syntax:

```yaml
[[ varNameHere ]]
```

## Example Config File with the Use of Variables

`variables.yaml`:

```yaml
gitlabUser: ironcore864
defaultBranch: main
gitlabCIGolangTemplate: https://gitlab.com/ironcore864/go-hello-world/-/raw/main/go.tpl
```

Example tool file with the variables defined in the above `variables.yaml`:

`tools.yaml`:

```yaml
tools:
- name: gitlabci-generic
  instanceID: default
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
- name: gitlabci-generic
  instanceID: default
  options:
    pathWithNamespace: ironcore864/go-hello-world
    branch: main
    templateURL: https://gitlab.com/ironcore864/go-hello-world/-/raw/main/go.tpl
    templateVariables:
      App: hello
```
