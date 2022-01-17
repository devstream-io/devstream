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
    # Helm repo information
    repo:
      # name of the Helm repo
      name: prometheus-community
      # url of the Helm repo
      url: https://prometheus-community.github.io/helm-charts
    # Helm chart information
    chart:
      # name of the chart
      name: prometheus-community/kube-prometheus-stack
      # release name of the chart
      release_name: dev
      # k8s namespace where kube-prometheus will be installed
      namespace: monitoring
      # need to create the namespace or not
      create_namespace: True
```

Currently, all the parameters in the example above are mandatory.
