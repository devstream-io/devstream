# kube-prometheus Plugin

This plugin installs [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus) in an existing Kubernetes cluster using the Helm chart.

## Usage

```yaml
--8<-- "kube-prometheus.yaml"
```

### Default Configs

| key              | default value                                      | description                                    |
| ----             | ----                                               | ----                                           |
| chart.chart_name | prometheus-community/kube-prometheus-stack         | chart name                                     |
| chart.timeout    | 5m                                                 | this config will wait 5 minutes to deploy      |
| upgradeCRDs      | true                                               | default update CRD config                      |
| chart.wait       | true                                               | whether to wait until installation is complete |
| repo.url         | https://prometheus-community.github.io/helm-charts | helm official repo address                     |
| repo.name        | prometheus-community                               | helm repo name                                 |

Currently, except for `values_yaml` and default configs, all the parameters in the example above are mandatory.
