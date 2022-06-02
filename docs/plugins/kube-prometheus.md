# kube-prometheus Plugin

This plugin installs [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus) in an existing Kubernetes cluster using the Helm chart.

## Usage

```yaml
--8<-- "kube-prometheus.yaml"
```

Currently, except for `values_yaml`, all the parameters in the example above are mandatory.
