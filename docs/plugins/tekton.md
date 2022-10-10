# tekton Plugin
This plugin installs [tekton](https://tekton.dev/) in an existing Kubernetes cluster using the Helm chart.

## Usage

The following content is an example of the "tool file".

For more information on the main config, the tool file and the var file of DevStream, see [Core Concepts Overview](../core-concepts/core-concepts.md#1-config) and [DevStream Configuration](../core-concepts/config.md).

```yaml
--8<-- "tekton.yaml"
```

### Default Configs

| key                | default value                                   | description                                        |
| ----               | ----                                            | ----                                               |
| chart.chartPath    | ""                                              | local chart path                                   |
| chart.chartName    | tekton/tekton-pipeline                          | chart name                                         |
| chart.timeout      | 5m                                              | this config will wait 5 minutes to deploy          |
| chart.upgradeCRDs  | true                                            | default update CRD config                          |
| chart.releaseName  | tekton                                          | helm release name                                  |
| chart.wait         | true                                            | whether to wait until installation is complete     |
| chart.namespace    | tekton                                          | namespace where helm to deploy                     |
| repo.url           | https://steinliber.github.io/tekton-helm-chart/ | helm community repo address                        |
| repo.name          | tekton                                          | helm repo name                                     |


Currently, except for `valuesYaml` and default configs, all the parameters in the example above are mandatory.
