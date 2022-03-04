## 1 `argocd` Plugin

This plugin installs [ArgoCD](https://argoproj.github.io/cd/) in an existing Kubernetes cluster using the Helm chart.

## 2 Usage:

```yaml
tools:
- name: argocd
  plugin:
    # name of the plugin
    kind: argocd
    # version of the plugin
    version: 0.2.0
  options:
    # need to create the namespace or not, default: false
    create_namespace: true
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
      # the time to wait for any individual Kubernetes operation (like Jobs for hooks). This defaults to 5m0s
      timeout: 5m
      # whether to perform a CRD upgrade during installation
      upgradeCRDs: true
```

Currently, all the parameters in the example above are mandatory.
