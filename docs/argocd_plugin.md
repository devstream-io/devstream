## 1 ArgoCD Plugin

This plugin installs [ArgoCD](https://argoproj.github.io/cd/) in an existing Kubernetes cluster using the Helm chart.

## 2 Usage:

```yaml
tools:
# name of the plugin
- name: argocd
  # version of the plugin
  version: 0.0.1
  # options for the plugin
  # Checkout the version from the GitHub releases.
  options:
    # Helm repo information
    repo:
      # name of the Helm repo
      name: argo
      # url of the Helm repo
      url: https://argoproj.github.io/argo-helm
    # Helm chart information
    chart:
      # name of the chart
      name: argo/argo-cd
      # release name of the chart
      release_name: argocd
      # K8s namespace where argocd will be installed
      namespace: argocd
      # need to create the namespace or not
      create_namespace: False
```

Currently, all the parameters in the example above are mandatory.
