# Install kube-prometheus with DevStream

## Default Configs

| key                | default value                        | description                                        |
| ----------------   | ------------------------------------ | ------------------------------------------------   |
| chart.chartPath    | ""                                   | local chart path                                   |
| chart.chartName    | prometheus-community/kube-prometheus-stack       | chart name                             |
| chart.version      | ""                                   | chart version                                      |
| chart.timeout      | 10m                                  | this config will wait 10 minutes to deploy Argo CD |
| chart.upgradeCRDs  | true                                 | default update CRD config                          |
| chart.releaseName  | prometheus                           | helm release name                                  |
| chart.namespace    | prometheus                           | namespace where helm to deploy                     |
| chart.wait         | true                                 | whether to wait until installation is complete     |
| repo.url           | https://prometheus-community.github.io/helm-charts | helm official repo address           |
| repo.name          | prometheus-community                 | helm repo name                                     |
