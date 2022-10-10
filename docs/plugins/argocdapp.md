# argocdapp Plugin

This plugin creates an [ArgoCD Application](https://argo-cd.readthedocs.io/en/stable/core_concepts/) custom resource.

**Notes:**

- ArgoCD itself must have been already installed before the usage of this plugin.
  To install ArgoCD, use the [argocd plugin](./argocd.md).
  Or you can use both plugins(argocd+argocdapp) at the same time.
  See [GitOps Toolchain](../best-practices/gitops.md) for more info.
- Currently, only the Helm chart is supported when creating the ArgoCD application.
- Modify the file accordingly. Especially remember to modify `ARGOCD_TOOL_NAME`.

## Usage

The following content is an example of the "tool file".

For more information on the main config, the tool file and the var file of DevStream, see [Core Concepts Overview](../core-concepts/core-concepts.md#1-config) and [DevStream Configuration](../core-concepts/config.md).

```yaml
--8<-- "argocdapp.yaml"
```

## Use Together with the `repo-scaffolding` Plugin

This plugin can be used together with the `repo-scaffolding` plugin (see document [here](./repo-scaffolding.md).)

For example, you can first use `repo-scaffolding` to bootstrap a Golang repo, then use this plugin to set up basic GitHub Actions CI workflows. In this scenario:

- This plugin can specify `repo-scaffolding` as a dependency, so that the dependency is first satisfied before executing this plugin.
- This plugin can refer to `repo-scaffolding`'s output to reduce copy/paste human error.

See the example below:

```yaml
---
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
- name: go-webapp-argocd-deploy
  plugin: argocdapp
  dependsOn: ["repo-scaffolding.golang-github"]
  options:
    app:
      name: hello
      namespace: argocd
    destination:
      server: https://kubernetes.default.svc
      namespace: default
    source:
      valuefile: values.yaml
      path: charts/go-hello-http
      repoURL: ${{repo-scaffolding.golang-github.outputs.repoURL}}
```

In the example above:

- We put `repo-scaffolding.golang-github` as dependency by using the `dependsOn` keyword.
- We used `repo-scaffolding.golang-github`'s output as input for the `githubactions-golang` plugin.

Pay attention to the `${{ xxx }}` part in the example. `${{ TOOL_NAME.PLUGIN.outputs.var}}` is the syntax for using an output.
