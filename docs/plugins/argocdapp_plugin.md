## 1 `argocdapp` Plugin

This plugin creates an [ArgoCD Application](https://argo-cd.readthedocs.io/en/stable/core_concepts/) custom resource.

**Notes:**
- ArgoCD itself must have been already installed before the usage of this plugin. To install ArgoCD, use the [argocd plugin](https://github.com/merico-dev/stream/blob/main/docs/argocd_plugin.md).
- Currently, only the Helm chart is supported when creating the ArgoCD application.
- At the moment, DevStream doesn't support dependency or concurrency yet. So, in the config file, the ArgoCD app plugin must be placed _after_ the ArgoCD plugin if you want to install ArgoCD first then create the ArgoCD application.

## 2 Usage:

```yaml
tools:
- name: helloworld
  plugin:
    # name of the plugin
    kind: argocdapp
    # version of the plugin
    # checkout the version from the GitHub releases
    version: 0.2.0
  # optional; if specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ "TOOL1_NAME.TOOL1_KIND", "TOOL2_NAME.TOOL2_KIND" ]
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
      repoURL: https://github.com/ironcore864/openstream-gitops-test.git
```

## 3. Use Together with the `github-repo-scaffolding-golang` Plugin

This plugin can be used together with the `github-repo-scaffolding-golang` plugin (see document [here](./github-repo-scaffolding-golang_plugin.md).)

For example, you can first use `github-repo-scaffolding-golang` to bootstrap a Golang repo, then use this plugin to set up basic GitHub Actions CI workflows. In this scenario:

- This plugin can specify `github-repo-scaffolding-golang` as a dependency, so that the dependency is first satisfied before executing this plugin.
- This plugin can refer to `github-repo-scaffolding-golang`'s output to reduce copy/paste human error.

See the example below:

```yaml
---
tools:
- name: go-webapp-repo
  plugin:
    kind: github-repo-scaffolding-golang
    version: 0.2.0
  options:
    owner: IronCore864
    repo: go-webapp-devstream-demo
    branch: main
    image_repo: ironcore864/go-webapp-devstream-demo
- name: go-webapp-argocd-deploy
  plugin:
    kind: argocdapp
    version: 0.2.0
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

Pay attention to the `${{ xxx }}` part in the example. `${{ TOOL_NAME.TOOL_KIND.outputs.var}}` is the syntax for using an output.
