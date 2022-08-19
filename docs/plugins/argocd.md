# argocd Plugin

This plugin installs [ArgoCD](https://argoproj.github.io/cd/) in an existing Kubernetes cluster using the Helm chart.

## Usage

```yaml
--8<-- "argocd.yaml"
```

### Default Configs

| key                | default value                        | description                                        |
| ----------------   | ------------------------------------ | ------------------------------------------------   |
| chart.chart_name   | argo/argo-cd                         | chart name                                         |
| chart.timeout      | 5m                                   | this config will wait 5 minutes to deploy argocd   |
| chart.upgradeCRDs  | true                                 | default update CRD config                          |
| chart.release_name | argocd                               | helm release name                                  |
| chart.namespace    | argocd                               | namespace where helm to deploy                     |
| chart.wait         | true                                 | whether to wait until installation is complete     |
| repo.url           | https://argoproj.github.io/argo-helm | helm official repo address                         |
| repo.name          | argo                                 | helm repo name                                     |
| create_namespace   | true                                 | make sure namespace exist                          |

Currently, except for `values_yaml` and default configs, all the parameters in the example above are mandatory.
