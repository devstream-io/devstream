# argocdapp Plugin

This plugin creates an [ArgoCD Application](https://argo-cd.readthedocs.io/en/stable/core_concepts/) custom resource.

**Notes:**
- ArgoCD itself must have been already installed before the usage of this plugin.
  To install ArgoCD, use the [argocd plugin](argocd).
  Or you can use both plugins(argocd+argocdapp) at the same time.
  See [GitOps Toolchain](/docs/best-practices/gitops) for more info.
- Currently, only the Helm chart is supported when creating the ArgoCD application.
- Modify the file accordingly. Especially remember to modify `ARGOCD_TOOL_NAME`.

## Usage

```yaml
tools:
# name of the tool
- name: argocdapp
  # id of the tool instance
  instanceID: default
  # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ "argocd.ARGOCD_INSTANCE_NAME" ]
  # options for the plugin
  options:
    # information on the ArgoCD Application
    app:
      # name of the ArgoCD Application
      name: hello
      # where the ArgoCD Application custom resource will be created
      namespace: argocd
    # destination of the ArgoCD Application
    destination:
      # on which server to deploy
      server: https://kubernetes.default.svc
      # in which namespace to deploy
      namespace: default
    # source of the application
    source:
      # which values file to use in the Helm chart
      valuefile: values.yaml
      # path of the Helm chart
      path: charts/go-hello-http
      # Helm chart repo URL, this is only an example, do not use this
      repoURL: YOUR_CHART_REPO_URL
```

## Use Together with the `github-repo-scaffolding-golang` Plugin

This plugin can be used together with the `github-repo-scaffolding-golang` plugin (see document [here](./github-repo-scaffolding-golang.md).)

For example, you can first use `github-repo-scaffolding-golang` to bootstrap a Golang repo, then use this plugin to set up basic GitHub Actions CI workflows. In this scenario:

- This plugin can specify `github-repo-scaffolding-golang` as a dependency, so that the dependency is first satisfied before executing this plugin.
- This plugin can refer to `github-repo-scaffolding-golang`'s output to reduce copy/paste human error.

See the example below:

```yaml
---
tools:
- name: go-webapp-repo
  plugin: github-repo-scaffolding-golang
  options:
    owner: IronCore864
    repo: go-webapp-devstream-demo
    branch: main
    image_repo: ironcore864/go-webapp-devstream-demo
- name: go-webapp-argocd-deploy
  plugin: argocdapp
  dependsOn: ["go-webapp-repo.github-repo-scaffolding-golang"]
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
      repoURL: ${{go-webapp-repo.github-repo-scaffolding-golang.outputs.repoURL}}
```

In the example above:

- We put `go-webapp-repo.github-repo-scaffolding-golang` as dependency by using the `dependsOn` keyword.
- We used `go-webapp-repo.github-repo-scaffolding-golang`'s output as input for the `githubactions-golang` plugin.

Pay attention to the `${{ xxx }}` part in the example. `${{ TOOL_NAME.PLUGIN.outputs.var}}` is the syntax for using an output.
