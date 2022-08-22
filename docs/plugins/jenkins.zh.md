# jenkins 插件

本插件使用 helm 在已有的 k8s 集群上安装 [Jenkins](https://jenkins.io)。

并且安装 [GitHub Pull Request Builder](https://plugins.jenkins.io/ghprb/) 插件和 [OWASP Markup Formatter](https://plugins.jenkins.io/antisamy-markup-formatter/) 插件；同时利用 OWASP Markup Formatter 插件激活 HTML 渲染模式。

## 配置

请将配置中的 `storageClass` 修改为已存在的 StorageClass.

```yaml
--8<-- "jenkins.yaml"
```

## 默认配置

| key                | default value             | description                                        |
| ----               | ----                      | ----                                               |
| chart.chart_name   | jenkins/jenkins           | helm 包名称                                        |
| chart.timeout      | 5m                        | 等待部署成功的时间                                 |
| chart.release_name | jenkins                   | helm 发布名称                                      |
| chart.upgradeCRDs  | true                      | 默认更新 CRD 配置（如果存在的话）                  |
| chart.wait         | true                      | 是否等待部署完成                                   |
| chart.namespace    | jenkins                   | helm 部署的命名空间名称                            |
| repo.url           | https://charts.jenkins.io | helm 官方仓库地址                                  |
| repo.name          | jenkins                   | helm 仓库名                                        |

当前，除了默认配置以外，所有配置项均为必填。

## 输出

这个插件有两个输出：

- `jenkinsURL` (格式: `hostname:port`, 例如: "localhost:8080")
- `jenkinsPasswordOfAdmin` 
