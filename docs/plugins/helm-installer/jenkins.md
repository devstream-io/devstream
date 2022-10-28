# Install Jenkins with DevStream

//TODO(daniel-hutao): to be updated

This plugin installs [Jenkins](https://jenkins.io) in an existing Kubernetes cluster using the Helm chart.

It also installs [GitHub Pull Request Builder(ghprb)](https://plugins.jenkins.io/ghprb/) and [OWASP Markup Formatter](https://plugins.jenkins.io/antisamy-markup-formatter/) plugins. Then enable HTML parsing using OWASP Markup Formatter Plugin , useful with ghprb plugin.

## Default Configs

| key                | default value             | description                                        |
| ----               | ----                      | ----                                               |
| chart.chartPath    | ""                        | local chart path                                   |
| chart.chartName    | jenkins/jenkins           | chart name                                         |
| chart.version      | ""                        | chart version                                      |
| chart.timeout      | 5m                        | this config will wait 5 minutes to deploy          |
| chart.upgradeCRDs  | true                      | default update CRD config                          |
| chart.releaseName  | jenkins                   | helm release name                                  |
| chart.namespace    | jenkins                   | namespace where helm to deploy                     |
| chart.wait         | true                      | whether to wait until installation is complete     |
| repo.url           | https://charts.jenkins.io | helm official repo address                         |
| repo.name          | jenkins                   | helm repo name                                     |

Please be sure to change the `storageClass` in the options of the config to an existing StorageClass.

Currently, expect default configs all the parameters in the example above are mandatory.

## Outputs

This plugin has two outputs:

- `jenkinsURL` (format: `hostname:port`, example: "localhost:8080")
- `jenkinsPasswordOfAdmin` 
