## 1 kube-prometheus Plugin

This plugin installs [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus) in an existing Kubernetes cluster using the Helm chart.

## 2 Usage:

```yaml
tools:
# name of the instance with kube-prometheus
- name: kube-prometheus-dev
  plugin:
    # kind of the plugin
    kind: kube-prometheus
    # version of the plugin
    version: 0.0.1
  # options for the plugin
  options:
    # need to create the namespace or not, default: false
    create_namespace: false
    # Helm repo information
    repo:
      # name of the Helm repo
      name: prometheus-community
      # url of the Helm repo
      url: https://prometheus-community.github.io/helm-charts
    # Helm chart information
    chart:
      # name of the chart
      chart_name: prometheus-community/kube-prometheus-stack
      # release name of the chart
      release_name: dev
      # k8s namespace where kube-prometheus will be installed
      namespace: monitoring
      # whether to wait for the release to be deployed or not
      wait: true
      # the time to wait for any individual Kubernetes operation (like Jobs for hooks)
      timeout: 5m
      # whether to perform a CRD upgrade during installation
      upgradeCRDs: true
```

Currently, all the parameters in the example above are mandatory.
