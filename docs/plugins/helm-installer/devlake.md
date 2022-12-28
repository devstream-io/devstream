# Install DevLake with DevStream

## Default Configs

| key                | default value                        | description                                        |
| ----------------   | ------------------------------------ | ------------------------------------------------   |
| chart.chartPath    | ""                                   | local chart path                                   |
| chart.chartName    | "devlake/devlake                     | chart name                                         |
| chart.version      | ""                                   | chart version                                      |
| chart.timeout      | 10m                                  | this config will wait 10 minutes to deploy DevLake |
| chart.upgradeCRDs  | true                                 | default update CRD config                          |
| chart.releaseName  | devlake                              | helm release name                                  |
| chart.namespace    | devlake                              | namespace where helm to deploy                     |
| chart.wait         | true                                 | whether to wait until installation is complete     |
| repo.url           | https://merico-dev.github.io/devlake-helm-chart | helm official repo address              |
| repo.name          | devlake                              | helm repo name                                     |
