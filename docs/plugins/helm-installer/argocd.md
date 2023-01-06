# Install Argo CD with DevStream

## InstanceID Prefix

The `instanceID` prefix must be `argocd`, the minimum tools configuration example:

```yaml
tools:
- name: helm-installer
  instanceID: argocd
```

## Default Configs

| key                | default value                        | description                                        |
| ----------------   | ------------------------------------ | ------------------------------------------------   |
| chart.chartPath    | ""                                   | local chart path                                   |
| chart.chartName    | argo/argo-cd                         | chart name                                         |
| chart.version      | ""                                   | chart version                                      |
| chart.timeout      | 10m                                  | this config will wait 10 minutes to deploy Argo CD |
| chart.upgradeCRDs  | true                                 | default update CRD config                          |
| chart.releaseName  | argocd                               | helm release name                                  |
| chart.namespace    | argocd                               | namespace where helm to deploy                     |
| chart.wait         | true                                 | whether to wait until installation is complete     |
| repo.url           | https://argoproj.github.io/argo-helm | helm official repo address                         |
| repo.name          | argo                                 | helm repo name                                     |
