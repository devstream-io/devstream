# helm-installer Plugin

## Usage

``` yaml
--8<-- "helm-installer.yaml"
```

### Argo CD

#### Default Configs

| key                | default value                        | description                                        |
| ----------------   | ------------------------------------ | ------------------------------------------------   |
| chart.chartPath    | ""                                   | local chart path                                   |
| chart.chartName    | argo/argo-cd                         | chart name                                         |
| chart.version      | ""                                   | chart version
| chart.timeout      | 5m                                   | this config will wait 5 minutes to deploy Argo CD   |
| chart.upgradeCRDs  | true                                 | default update CRD config                          |
| chart.releaseName  | argocd                               | helm release name                                  |
| chart.namespace    | argocd                               | namespace where helm to deploy                     |
| chart.wait         | true                                 | whether to wait until installation is complete     |
| repo.url           | https://argoproj.github.io/argo-helm | helm official repo address                         |
| repo.name          | argo                                 | helm repo name                                     |
