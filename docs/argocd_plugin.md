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
  # checkout the version from the GitHub releases
  options:
    # need to create the namespace or not, default: false
    create_namespace: false
    # Helm repo information
    repo:
      # name of the Helm repo
      name: argo
      # url of the Helm repo
      url: https://argoproj.github.io/argo-helm
    # Helm chart information
    chart:
      # name of the chart
      chart_name: argo/argo-cd
      # release name of the chart
      release_name: argocd
      # k8s namespace where argocd will be installed
      namespace: argocd
      # whether to wait for the release to be deployed or not
      wait: true
      # the time to wait for any individual Kubernetes operation (like Jobs for hooks)
      timeout: 5m
      # whether to perform a CRD upgrade during installation
      upgradeCRDs: true
```

Currently, all the parameters in the example above are mandatory.
