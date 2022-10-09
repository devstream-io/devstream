# kube-prometheus Plugin

This plugin installs [kube-prometheus](https://github.com/prometheus-operator/kube-prometheus) in an existing Kubernetes cluster using the Helm chart.

## Usage

The following content is an example of the "tool file".

For more information on the main config, the tool file and the var file of DevStream, see [Core Concepts Overview](../core-concepts/core-concepts.md#1-config) and [DevStream Configuration](../core-concepts/config.md).

```yaml
--8<-- "kube-prometheus.yaml"
```

### Default Configs

| key                | default value                                      | description                                        |
| ----               | ----                                               | ----                                               |
| chart.chartPath    | ""                                                 | local chart path                                   |
| chart.chartName    | prometheus-community/kube-prometheus-stack         | chart name                                         |
| chart.timeout      | 5m                                                 | this config will wait 5 minutes to deploy          |
| chart.releaseName  | prometheus                                         | helm release name                                  |
| chart.upgradeCRDs  | true                                               | default update CRD config                          |
| chart.wait         | true                                               | whether to wait until installation is complete     |
| chart.namespace    | prometheus                                         | namespace where helm to deploy                     |
| repo.url           | https://prometheus-community.github.io/helm-charts | helm official repo address                         |
| repo.name          | prometheus-community                               | helm repo name                                     |

Currently, except for `valuesYaml` and default configs, all the parameters in the example above are mandatory.
