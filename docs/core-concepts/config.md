# Config

DevStream uses YAML files to describe your DevOps toolchain configuration.

## Main Config File

By default, `dtm` tries to use `./config.yaml` (under your current directory) as the main config.

The main config contains three sections:

- `varFile`: the path/to/your/variables file
- `toolFile`: the path/to/your/tools configuration file
- `state`: configuration of DevStream state. For example, 

### Example Main Config File

See the `config.yaml` example below:

```yaml
varFile: variables.yaml

toolFile: tools.yaml

state:
  backend: local
  options:
    stateFile: devstream.state
```

### Variables File

The var file is a YAML file containing key-value pairs.

_At the moment, nested/composite values (for example, the value is a list/dictionary) are not supported yet._

See the `variables.yaml` example below:

```yaml
githubUsername: daniel-hutao
repoName: dtm-test-go
defaultBranch: main
dockerhubUsername: exploitht
```

### Tool File

Tool file contains a list of tools.

The tool file contains:

- Only one section (at the moment), which is `tools`.
- `tools` is a list of dictionaries.
- Each dictionary defines a DevOps "tool" which is managed by a DevStream plugin
- Each dictionary (tool) has the following mandatory fields:
    - `name`: the name of the tool/plugin, string, without underscore
    - `instanceID`: the id of this tool instance
    - you can have duplicated `name` in one config file, and you can also have duplicated `instanceID` in one config file, but the `name + instanceID` combination must be unique in one config file
- Each dictionary (tool) has an optional field which is `options`, which in turn is a dictionary containing parameters for that specific plugin. For plugins' parameters, see the "plugins" section of this document.
- Each directory (tool) has an optional field which is `dependsOn`. Continue reading for detail about dependencies.

See the `tools.yaml` example down below:

```yaml
tools:
- name: github-repo-scaffolding-golang
  instanceID: default
  options:
    owner: [[ githubUsername ]]
    org: ""
    repo: [[ repoName ]]
    branch: [[ defaultBranch ]]
    image_repo: [[ dockerhubUsername ]]/[[ repoName ]]
- name: jira-github-integ
  instanceID: default
  dependsOn: [ "github-repo-scaffolding-golang.default" ]
  options:
    owner: [[ githubUsername ]]
    repo: [[ repoName ]]
    jiraBaseUrl: https://xxx.atlassian.net
    jiraUserEmail: foo@bar.com
    jiraProjectKey: zzz
    branch: main
```

### State

The `state` section specifies where to store DevStream state. As of now (v0.5.0), we only support local backend.

From v0.6.0 on, we will support both "local" and "s3" backend store the DevStream state.

Read the section [The State Section in the Main Config](./stateconfig.md) for more details.

## Default Values

By default, `dtm` uses `config.yaml` as the main config file.

### Specifying a Main Config File Explicitly 

You can override the default value with `-f` or `--config-file`. Examples:

```shell
dtm apply -f path/to/your/config.yaml
dtm apply --config-file path/to/your/config.yaml
```

### No Defaults for varFile and toolFile

For `varFile` and `toolFile`, no default values are provided.

If `varFile` isn't specified in the main config, `dtm` will not use any var files, even if there is already a file named `variables.yaml` under the current directory.

Similarly, if `toolFile` isn't specified in the main config, `dtm` will throw an error, even if there is a `tools.yaml` file under the current directory.
