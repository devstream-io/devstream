# 使用 DevStream 部署 Vault

## 默认配置

| 配置项              | 默认值                    | 描述                                 |
| ----               | ----                     | ----                                |
| chart.chartPath    | ""                       | 本地 chart 包路径                     |
| chart.chartName    | hashicorp/vault          | chart 包名称                         |
| chart.version      | ""                       | chart 包版本                         |
| chart.timeout      | 10m                      | helm install 的超时时间               |
| chart.upgradeCRDs  | true                     | 是否更新 CRDs（如果有）                |
| chart.releaseName  | vault                    | helm 发布名称                         |
| chart.namespace    | vault                    | 部署的命名空间                         |
| chart.wait         | true                     | 是否等待部署完成                       |
| repo.url           | https://helm.releases.hashicorp.com | helm 仓库地址              |
| repo.name          | hashicorp                | helm 仓库名                           |
