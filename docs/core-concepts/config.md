# Config

DevStream uses YAML files to define your DevOps platform.

## Sections

As aforementioned in the overview, there are several sections in the config:

- `config`
- `vars`
- `tools`
- `apps`
- `pipelineTemplates`

Among which, `config` is mandatory, and you should have at least either `tools` or `apps`. Others are optional.

## Config File V.S. Config Folder

DevStream supports both:

- a single config file: put all sections of the config into one YAML file
- a directory: put multiple files in one directory, as long as the file names end with `.yaml` or `.yml`

When you run the `init` (or other commands), add `-f` or `--config-file`.

Examples:

- single file: `dtm init -f config.yaml`
- directory: `dtm init -f dirname`

## The Main Config

DevStream's own config, the `config` section, contains where to store the state. See [here](./state.md) for more details.

## Variables

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

## Tools Config

The `tools` section defines DevOps tools. Read [this](./tools.md) for all the detail.

## Apps/pipelineTemplates Config

The `apps` section defines Apps and the `pipelineTemplates` defines templates for pipelines. See [here](./apps.md) for more detail.
