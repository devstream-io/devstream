# tekton Plugin
This plugin installs [tekton](https://tekton.dev/) in an existing Kubernetes cluster using the Helm chart.

## Usage

```yaml
--8<-- "tekton.yaml"
```

### Default Configs

| key                | default value                                   | description                                        |
| ----               | ----                                            | ----                                               |
| chart.chart_name   | tekton/tekton-pipeline                          | chart name                                         |
| chart.timeout      | 5m                                              | this config will wait 5 minutes to deploy          |
| chart.upgradeCRDs  | true                                            | default update CRD config                          |
| chart.release_name | tekton                                          | helm release name                                  |
| chart.wait         | true                                            | whether to wait until installation is complete     |
| chart.namespace    | tekton                                          | namespace where helm to deploy                     |
| repo.url           | https://steinliber.github.io/tekton-helm-chart/ | helm community repo address                        |
| repo.name          | tekton                                          | helm repo name                                     |
| create_namespace   | false                                           | whether to create namespace if namespace not eixst |


Currently, except for `values_yaml` and default configs, all the parameters in the example above are mandatory.
