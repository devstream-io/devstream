# argocd Plugin

This plugin installs [ArgoCD](https://argoproj.github.io/cd/) in an existing Kubernetes cluster using the Helm chart.

## Usage

The following content is an example of the "tool file".

For more information on the main config, the tool file and the var file of DevStream, see [Core Concepts Overview](../core-concepts/core-concepts.md#1-config) and [DevStream Configuration](../core-concepts/config.md).

```yaml
--8<-- "argocd.yaml"
```

### Default Configs

| key                | default value                        | description                                        |
| ----------------   | ------------------------------------ | ------------------------------------------------   |
| chart.chartPath    | ""                                   | local chart path                                   |
| chart.chartName    | argo/argo-cd                         | chart name                                         |
| chart.timeout      | 5m                                   | this config will wait 5 minutes to deploy argocd   |
| chart.upgradeCRDs  | true                                 | default update CRD config                          |
| chart.releaseName  | argocd                               | helm release name                                  |
| chart.namespace    | argocd                               | namespace where helm to deploy                     |
| chart.wait         | true                                 | whether to wait until installation is complete     |
| repo.url           | https://argoproj.github.io/argo-helm | helm official repo address                         |
| repo.name          | argo                                 | helm repo name                                     |

Currently, except for `valuesYaml` and default configs, all the parameters in the example above are mandatory.
