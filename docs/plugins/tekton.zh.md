# tekton 插件

TODO(dtm): 在这里添加文档.

## 用例

```yaml
--8<-- "tekton.yaml"
```

### Default Configs

| key              | default value                                   | description                       |
| ----             | ----                                            | ----                              |
| chart.chart_name | tekton/tekton-pipeline                          | helm 包名称                       |
| chart.timeout    | 5m                                              | 等待部署成功的时间                |
| upgradeCRDs      | true                                            | 默认更新 CRD 配置（如果存在的话） |
| chart.wait       | true                                            | 是否等待部署完成                  |
| repo.url         | https://steinliber.github.io/tekton-helm-chart/ | helm 官方仓库地址                 |
| repo.name        | tekton                                          | helm 仓库名                       |

目前除了 `values_yaml` 和默认配置，以上其它配置项均为必填项。
