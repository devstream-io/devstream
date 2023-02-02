# 使用 DevStream 部署 DevLake

## 前缀匹配

`instanceID` 的前缀需要是 `devlake`，最小化 tools 配置示例：

```yaml
tools:
- name: helm-installer
  instanceID: devlake
```

## 默认配置

| 配置项              | 默认值                    | 描述                                 |
| ----               | ----                     | ----                                |
| chart.chartPath    | ""                       | 本地 chart 包路径                     |
| chart.chartName    | devlake/devlake          | chart 包名称                         |
| chart.version      | ""                       | chart 包版本                         |
| chart.timeout      | 10m                      | helm install 的超时时间               |
| chart.upgradeCRDs  | true                     | 是否更新 CRDs（如果有）                |
| chart.releaseName  | devlake                   | helm 发布名称                         |
| chart.namespace    | devlake                   | 部署的命名空间                         |
| chart.wait         | true                     | 是否等待部署完成                       |
| repo.url           | https://apache.github.io/incubator-devlake-helm-chart | helm 仓库地址   |
| repo.name          | devlake                   | helm 仓库名                           |
