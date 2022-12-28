# 使用 DevStream 部署 Tekton

## 默认配置

| 配置项              | 默认值                                           | 描述                                 |
| ----               | ----                                            | ----                                |
| chart.chartPath    | ""                                              | 本地 chart 包路径                     |
| chart.chartName    | tekton/tekton-pipeline                          | helm 包名称                          |
| chart.version      | ""                                              | chart 版本                           |
| chart.timeout      | 5m                                              | 等待部署成功的时间                     |
| chart.upgradeCRDs  | true                                            | 默认更新 CRD 配置（如果存在的话）        |
| chart.releaseName  | tekton                                          | helm 发布名称                        |
| chart.wait         | true                                            | 是否等待部署完成                      |
| chart.namespace    | tekton                                          | helm 部署的命名空间名称                |
| repo.url           | https://steinliber.github.io/tekton-helm-chart/ | helm 官方仓库地址                     |
| repo.name          | tekton                                          | helm 仓库名                          |
