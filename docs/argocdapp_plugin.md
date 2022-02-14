## 1 ArgoCD Application Plugin

This plugin installs an ArgoCD application.

**Notes:**
- ArgoCD itself must have been already installed before the usage of this plugin. To install ArgoCD, use [this plugin](https://github.com/merico-dev/stream/blob/main/docs/argocd_plugin.md).
- Currently, only the Helm chart is supported when creating the ArgoCD application.
- At the moment, DevStream doesn't support dependency or concurrency yet. So, in the config file, the ArgoCD app plugin must be placed _after_ the ArgoCD plugin if you want to install ArgoCD first then create the ArgoCD application.

## 2 Usage:

```yaml
tools:
- name: argocdapp
  plugin:
    # name of the plugin
    kind: argocdapp
    # version of the plugin
    # checkout the version from the GitHub releases
    version: 0.0.2
  # options for the plugin
  options:
    # information on the ArgoCD application
    app:
      # name of the ArgoCD application
      name: hello
      # where the ArgoCD application CRD will be created
      namespace: argocd
    # destination of the application
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
      # Helm chart repo URL
      repoURL: https://github.com/ironcore864/openstream-gitops-test.git
```

Currently, all the parameters in the example above are mandatory.
