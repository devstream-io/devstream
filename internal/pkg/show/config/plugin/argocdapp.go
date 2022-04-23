package plugin

var ArgocdappDefaultConfig = `tools:
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
      repoURL: YOUR_CHART_REPO_URL`
