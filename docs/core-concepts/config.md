# Config

Now let's have a look at some config examples.

TL;DR: [see the complete config file examle at the end of this doc](#7-putting-it-all-together).

---

## 1 The Config File

As mentioned in the overview section, the main config contains many sections, like:

- `config`
- `vars`
- `tools`
- `apps`
- `pipelineTemplates`

By default, `dtm` tries to use `./config.yaml` (under your working directory) as the main config.

You can override the default value with `-f` or `--config-file`. Examples:

```shell
dtm apply -f path/to/your/config.yaml
dtm apply --config-file path/to/your/config.yaml
```

---

## 2 State Config

The `config` section specifies where to store the DevStream state.

### 2.1 Local

```yaml
config:
  state:
    backend: local
    options:
      stateFile: devstream.state
```

_Note: The `stateFile` under the `options` section is mandatory for the local backend._

### 2.2 Kubernetes

```yaml
config:
  state:
    backend: k8s
    options:
      namespace: devstream # optional, default is "devstream", will be created if not exists
      configmap: state     # optional, default is "state", will be created if not exists
```

### 2.3 S3

```yaml
config:
  state:
    backend: s3
    options:
      bucket: devstream-remote-state
      region: ap-southeast-1
      key: devstream.state
```

_Note: the `bucket`, `region`, and `key` under the `options` section are all mandatory fields for the s3 backend._

---

## 3 Variables

To not repeat yourself (Don't repeat yourself, DRY, see [here](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself),) we can define some variables in a var file and use the vars in the config file.

The vars section is a YAML dict containing key-value pairs. Example:

Example:

```yaml
vars:
  repoOwner: RepoOwner
  repoTemplateBaseURL: github.com/devstream-io
  imageRepoOwner: ImageRepoOwner
```

To use these variables in a config file, use the following syntax:

```yaml
[[ varNameHere ]]
```

---

## 4 Tools

The tools section contains a list of tools.

```yaml
tools:
- name: repo-scaffolding
  instanceID: myapp
  options:
    destinationRepo:
      owner: [[ githubUser ]]
      name: [[ app ]]
      branch: main
      scmType: github
    sourceRepo:
      org: devstream-io
      name: dtm-scaffolding-python
      scmType: github
    vars:
      imageRepo: [[ dockerUser ]]/[[ app ]]
- name: github-actions
  instanceID: default
  dependsOn: [ repo-scaffolding.myapp ]
  options:
    scm:
      owner: [[ githubUser ]]
      name: [[ app ]]
      branch: main
      scmType: github
    pipeline:
      language:
        name: python
        framework: flask
      imageRepo:
        user: [[ dockerUser ]]
- name: helm-installer
  instanceID: argocd
- name: argocdapp
  instanceID: default
  dependsOn: [ "helm-installer.argocd", "github-actions.default" ]
  options:
    app:
      name: [[ app ]]
      namespace: argocd
    destination:
      server: https://kubernetes.default.svc
      namespace: default
    source:
      valuefile: values.yaml
      path: helm/[[ app ]]
      repoURL: ${{repo-scaffolding.myapp.outputs.repoURL}}
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

---

## 5 Apps

The Apps section looks like the following:

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
    name: service-A   # defaults to "app.name"
    url: github.com/devstream-io/repo-name   # optional. if exists, no need for the scm/owner/org/name sections
    apiURL: gitlab.com/some/path/to/your/api # optional, in case of self-managed git
  repoTemplate:   # optional. if repoTemplate isn't empty, create repo according to the template
    scmType: github
    owner: devstream-io
    org: devstream-io # either owner or org must exist
    name: dtm-repo-scaffolding-golang
    url: github.com/devstream-io/repo-name   # optional. if exists, no need for the scm/owner/org/name sections
    vars:  # optional
      foo: bar  # variables used for repoTemplate specifically
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
- name: myapp1
  spec:
    language: python
    framework: django
  repo:
    url: github.com/ironcore864/myapp1
  repoTemplate:
    url: github.com/devstream-io/dtm-repo-scaffolding-python
  ci:
  - type: github-actions
  cd:
  - type: argocdapp
```

---

## 6 Pipeline Templates

You can define some pipeline templates in the pipelineTemplates section:

```yaml
pipelineTemplates:
- name: ci-pipeline-1
  type: github-actions # corresponds to a DevStream plugin
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

Here's a complete example:

```yaml
config:
  state:
  backend: local
  options:
    stateFile: devstream.state

tools:
- name: helm-installer
  instanceID: argocd

apps:
- name: myapp1
  spec:
    language: python
    framework: django
  repo:
    url: github.com/ironcore864/myapp1
  repoTemplate:
    url: github.com/devstream-io/dtm-repo-scaffolding-python
  ci:
  - type: github-actions
  cd:
  - type: argocdapp
```
