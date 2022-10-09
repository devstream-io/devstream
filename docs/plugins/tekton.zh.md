# tekton 插件

## 用例

下面的配置文件展示的是"tool file"的内容。

关于更多关于DevStream的主配置、tool file、var file的信息，请阅读[核心概念概览](../core-concepts/core-concepts.zh.md)和[DevStream配置](../core-concepts/config.zh.md).

```yaml
--8<-- "tekton.yaml"
```

### Default Configs

| key                | default value                                   | description                                      |
| ----               | ----                                            | ----                                             |
| chart.chartPath    | ""                                              | 本地 chart 包路径                                  |
| chart.chartName    | tekton/tekton-pipeline                          | helm 包名称                                       |
| chart.timeout      | 5m                                              | 等待部署成功的时间                                  |
| chart.upgradeCRDs  | true                                            | 默认更新 CRD 配置（如果存在的话）                     |
| chart.releaseName  | tekton                                          | helm 发布名称                                     |
| chart.wait         | true                                            | 是否等待部署完成                                   |
| chart.namespace    | tekton                                          | helm 部署的命名空间名称                             |
| repo.url           | https://steinliber.github.io/tekton-helm-chart/ | helm 官方仓库地址                                  |
| repo.name          | tekton                                          | helm 仓库名                                       |

目前除了 `valuesYaml` 和默认配置，以上其它配置项均为必填项。
