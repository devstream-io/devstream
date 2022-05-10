# argocd Plugin

This plugin installs [ArgoCD](https://argoproj.github.io/cd/) in an existing Kubernetes cluster using the Helm chart.

## Usage

```yaml
tools:
# name of the tool
- name: argocd
  # id of the tool instance
  instanceID: default
  # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: []
  # options for the plugin
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
      # k8s namespace where ArgoCD will be installed
      namespace: argocd
      # whether to wait for the release to be deployed or not
      wait: true
      # the time to wait for any individual Kubernetes operation (like Jobs for hooks). This defaults to 5m0s
      timeout: 5m
      # whether to perform a CRD upgrade during installation
      upgradeCRDs: true
      # custom configuration (Optional). You can refer to [ArgoCD values.yaml](https://github.com/argoproj/argo-helm/blob/master/charts/argo-cd/values.yaml)
      values_yaml: |
        controller:
          service: 
            port: 8080
        redis:
          image:
            tag: 6.2.6-alpine3.15

```

Currently, except for `values_yaml`, all the parameters in the example above are mandatory.
