# Install Tekton with DevStream

## Default Configs

| key                | default value                                   | description                                        |
| ----               | ----                                            | ----                                               |
| chart.chartPath    | ""                                              | local chart path                                   |
| chart.chartName    | tekton/tekton-pipeline                          | chart name                                         |
| chart.version      | ""                                              | chart version                                      |
| chart.timeout      | 10m                                             | this config will wait 5 minutes to deploy          |
| chart.upgradeCRDs  | true                                            | default update CRD config                          |
| chart.releaseName  | tekton                                          | helm release name                                  |
| chart.wait         | true                                            | whether to wait until installation is complete     |
| chart.namespace    | tekton                                          | namespace where helm to deploy                     |
| repo.url           | https://steinliber.github.io/tekton-helm-chart/ | helm community repo address                        |
| repo.name          | tekton                                          | helm repo name                                     |
