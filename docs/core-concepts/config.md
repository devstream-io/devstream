# Config

Now let's have a look at some config examples.

---

## 1 The Main Config File

As mentioned in the overview section, the main config contains many settings, like:

- `pluginDir`
- `state`
- `varFile`
- `toolFile`[here](./tools-apps.md).
- `appFile`
- `templateFile`

Here's an example of the main config file:

```yaml
varFile: "./variables.yaml"
toolFile: "./tools.yaml"
pluginDir: ""    # empty, use the default value: ~/.devstream/plugins
appFile: "./apps.yaml"
templateFile: "./templates.yaml"
state:           # state config, backend can be "local", "s3", or "k8s"
  backend: local
  options:
    stateFile: devstream.state
```

By default, `dtm` tries to use `./config.yaml` (under your working directory) as the main config.

You can override the default value with `-f` or `--config-file`. Examples:

```shell
dtm apply -f path/to/your/config.yaml
dtm apply --config-file path/to/your/config.yaml
```

No default values are provided for `varFile` and `toolFile`. If `varFile` isn't specified in the main config, `dtm` will not use any var files, even if there is already a file named `variables.yaml` under the current directory. Similarly, if `toolFile` isn't specified in the main config, `dtm` will throw an error, even if there is a `tools.yaml` file under the current directory.

---

## 2 Variables File

To not repeat yourself (Don't repeat yourself, DRY, see [here](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself),) we can define some variables in a var file and use the vars in the config file.

The var file is a YAML file containing key-value pairs.

_Note: at the moment, nested/composite values (for example, the value is a list/dictionary) are not supported yet._

Example:

```yaml
githubUsername: daniel-hutao
repoName: dtm-test-go
defaultBranch: main
dockerhubUsername: exploitht
```

To use these variables in a config file, use the following syntax:

```yaml
[[ varNameHere ]]
```

---

## 3 State Config

The `state` section specifies where to store the DevStream state.

### 3.1 Local

```yaml
# part of the main config file
state:
  backend: local
  options:
    stateFile: devstream.state
```

_Note: The `stateFile` under the `options` section is mandatory for the local backend._

### 3.2 Kubernetes

```yaml
# part of the main config file
state:
  backend: k8s
  options:
    namespace: devstream # optional, default is "devstream", will be created if not exists
    configmap: state     # optional, default is "state", will be created if not exists
```

### 3.3 S3

```yaml
# part of the main config file
state:
  backend: s3
  options:
    bucket: devstream-remote-state
    region: ap-southeast-1
    key: devstream.state
```

_Note: the `bucket`, `region`, and `key` under the `options` section are all mandatory fields for the s3 backend._

---

## 4 Tool File

The tool file contains a list of tools.

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

Example:

```yaml
tools:
- name: repo-scaffolding
  instanceID: golang-github
  options:
    destinationRepo:
      owner: [[ githubUsername ]]
      org: ""
      repo: [[ repoName ]]
      branch: [[ defaultBranch ]]
      repoType: github
    vars:
      ImageRepo: "[[ dockerhubUsername ]]/[[ repoName ]]"
    sourceRepo:
      org: devstream-io
      repo: dtm-scaffolding-golang
      repoType: github
- name: jira-github-integ
  instanceID: default
  dependsOn: [ "repo-scaffolding.golang-github" ]
  options:
    owner: [[ githubUsername ]]
    repo: [[ repoName ]]
    jiraBaseUrl: https://xxx.atlassian.net
    jiraUserEmail: foo@bar.com
    jiraProjectKey: zzz
    branch: main
```


If you want tool A to be installed before tool B, you can let tool B depend on tool A.

The syntax for dependency is:
    
```yaml
dependsOn: [ "ToolName.ToolInstanceID" ]
```

Since `dependsOn` is a list, a tool can have multiple dependencies:

```
dependsOn: [ "ToolName1.ToolInstanceID1", "ToolName2.ToolInstanceID2", "..." ]
```

The example config above shows that the second tool depends on the first.

---

## 5 App File

A full app file looks like the following:

```yaml
apps:
- name: service-A
  spec:
    language: python
    framework: django
  repo:
    scmType: github
    owner: devstream-io
    org: devstream-io # either owner or org must exist
    name: service-A   # defaults to "name"
    url: github.com/devstream-io/repo-name   # optional. if exists, no need for the scm/owner/org/name sections
    apiURL: gitlab.com/some/path/to/your/api # optional, in case of self-managed git
  repoTemplate:   # optional. if repoTemplate isn't empty, create repo according to the template
    scmType: github
    owner: devstream-io
    org: devstream-io # either owner or org must exist
    name: dtm-scaffolding-golang
    url: github.com/devstream-io/repo-name   # optional. if exists, no need for the scm/owner/org/name sections
  ci:
  - type: template              # type template means it's a reference to the pipeline template. read the next section.
    templateName: ci-pipeline-1
    options: # optional, used to override defaults in the template
    vars:    # optional, key/values to be passed to the template
      key: value
  cd:
  - type: template              # type template means it's a reference to the pipeline template. read the next section.
    templateName: cd-pipeline-1
    options: # optional, used to override defaults in the template
    vars:    # optional, key/values to be passed to the template
      key: value
```

Since many key/values have default values, it's possible to use the following minimum config for the apps section:

```yaml
apps:
- name: service-B # minimum config demo
  spec:
    language: python
    framework: django
  repo:
    url: github.com/devstream-io/repo-name
  repoTemplate:
    url: github.com/devstream-io/repo-name
  ci:
    - type: githubactions
  cd:
    - type: argocdapp
```

---

## 6 Template File

You can define some pipeline templates in the template file, like the following:

```yaml
pipelineTemplates:
- name: ci-pipeline-1
  type: githubactions # corresponds to a DevStream plugin
  options:
    branch: main      # optional
    docker:
      registry:
        type: dockerhub
        username: [[ dockerUser ]]
        repository: [[ app ]]
- name: cd-pipeline-1
  type: argocdapp
  options:
    app:
      namespace: argocd
    destination:
      server: https://kubernetes.default.svc
      namespace: default
    source:
      valuefile: values.yaml
      path: helm/[[ app ]]
      repoURL: ${{repo-scaffolding.myapp.outputs.repoURL}}
```

Then, you can refer to these pipeline templates in the apps file.

---

## 7 Putting It All Together

Here's a complete example.

Main config file:

```yaml
toolFile: "./tools.yaml"
appFile: "./apps.yaml"
templateFile: "./templates.yaml"
state:
  backend: local
  options:
    stateFile: devstream.state
```

Tool file:

```yaml
tools:
- name: argocd
  instanceID: default
```

App file:

```yaml
apps:
- name: service-B
  spec:
    language: python
    framework: django
  repo:
    url: github.com/devstream-io/repo-name
  repoTemplate:
    url: github.com/devstream-io/repo-name
  ci:
    - type: githubactions
  cd:
    - type: argocdapp
```

In this example, we didn't use a var file, or template file, since they are optional and we don't need them in this example.
