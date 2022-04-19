package plugin

var VaultDefaultConfig = `tools:
- name: vault
  # name of the plugin
  plugin: vault
  options:
    # need to create the namespace or not, default: false
    create_namespace: true
    repo:
      # name of the Helm repo
      name: hashicorp
      # url of the Helm repo
      url: https://helm.releases.hashicorp.com
    # Helm chart information
    chart:
      # name of the chart
      chart_name: hashicorp/vault
      # release name of the chart
      release_name: vault
      # k8s namespace where Vault will be installed
      namespace: devstreamtest1
      # whether to wait for the release to be deployed or not
      wait: true
      # the time to wait for any individual Kubernetes operation (like Jobs for hooks). This defaults to 5m0s
      timeout: 5m
      values_yaml: |
      server:
        affinity: ""
        ha:
          enabled: true
        namespaceSelector:
          matchLabels:
            injection: enabled`
