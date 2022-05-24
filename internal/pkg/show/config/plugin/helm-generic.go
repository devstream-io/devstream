package plugin

var HelmGenericDefaultConfig = `tools:
# name of the tool
- name: helm-generic
  # id of the tool instance
  instanceID: default
  # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: []
  # options for the plugin
  options:
    # need to create the namespace or not, default: false
    create_namespace: true
    repo:
      # name of the Helm repo, e.g. argo
      name: HELM_REPO_NAME
      # url of the Helm repo, e.g. https://argoproj.github.io/argo-helm
      url: YOUR_CHART_REPO
    # Helm chart information
    chart:
      # name of the chart, e.g. argo/argo-cd
      chart_name: CHART_NAME
      # release name of the chart, e.g. argocd
      release_name: RELEASE_NAME
      # k8s namespace, e.g. argocd
      namespace: YOUR_CHART_NAMESPACE
      # whether to wait for the release to be deployed or not
      wait: true
      # the time to wait for any individual Kubernetes operation (like Jobs for hooks). This defaults to 5m0s
      timeout: 5m
      # whether to perform a CRD upgrade during installation
      upgradeCRDs: true
      # custom configuration (Optional). e.g. You can refer to [ArgoCD values.yaml](https://github.com/argoproj/argo-helm/blob/master/charts/argo-cd/values.yaml)
      values_yaml: |
        FOO: BAR`
