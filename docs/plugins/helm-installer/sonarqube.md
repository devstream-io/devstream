# Install SonarQube with DevStream

## InstanceID Prefix

The `instanceID` prefix must be `sonarqube`, the minimum tools configuration example:

```yaml
tools:
- name: helm-installer
  instanceID: sonarqube
```

### Default Configs

| key                | default value                             | description                                        |
| ----               | ----                                      | ----                                               |
| chart.chartPath    | ""                                        | local chart path                                   |
| chart.chartName    | sonarqube/sonarqube                       | community chart name                               |
| chart.version      | ""                                        | chart version                                      |
| chart.timeout      | 10m                                       | this config will wait 5 minutes to deploy          |
| chart.releaseName  | sonarqube                                 | helm release name                                  |
| chart.upgradeCRDs  | true                                      | default update CRD config                          |
| chart.namespace    | sonarqube                                 | namespace where helm to deploy                     |
| chart.wait         | true                                      | whether to wait until installation is complete     |
| repo.url           | https://SonarSource.github.io/helm-chart-sonarqube | helm repo address                         |
| repo.name          | sonarqube                                 | helm repo name                                     |
