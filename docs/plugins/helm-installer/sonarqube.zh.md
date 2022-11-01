# 使用 DevStream 部署 SonarQube

## 1. 前置要求

- 有一个可用的 Kubernetes 集群，版本 1.19+。
- 运行 Sonarqube 小规模服务至少需要 2GB 的 RAM 内存。
- Sonarqube 官方镜像目前只支持 linux/amd64 架构。
- 更多硬件配置可参考[官方网站](https://docs.sonarqube.org/latest/requirements/hardware-recommendations/)。

## 2、部署架构

Sonarqube 内部会使用 Elastcisearch 来做搜索的索引，所以生产环境中需要通过挂载目录的方式持久化数据。
另外 Sonarqube 也需要一个外部数据库来存储数据，目前支持 `PostgreSQL`，`Oracle`，`Microsoft SQL Sever`，默认使用 `PostgreSQL`。

## 3、开始部署

下文将介绍如何配置 `sonarqube` 插件，完成 Sonarqube 应用的部署。

### 3.1、默认配置

`sonarqube` 插件的配置项多数都有默认值，具体默认值信息如下表：

| 配置项             | 默认值                    | 描述                                 |
|-------------------| ----                     | ----                                |
| chart.chartName   | sonarqube/sonarqube      | helm chart 包名称                    |
| chart.timeout     | 10m                      | helm install 的超时时间               |
| chart.version     | ""                       | chart 版本                           |
| chart.upgradeCRDs | true                     | 是否更新 CRDs（如果有）                 |
| chart.releaseName | sonarqube                | helm 发布名称                         |
| chart.wait        | true                     | 是否等待部署完成                       |
| chart.namespace   | sonarqube                | 部署的命名空间                         |
| repo.url          | https://SonarSource.github.io/helm-chart-sonarqube| helm 仓库地址  |
| repo.name         | sonarqube                | helm 仓库名                           |

### 3.2、测试环境

在测试环境中可以使用如下配置：

```yaml
tools:
- name: helm-installer
  instanceID: sonarqube-001
  dependsOn: [ ]
  options:
    valuesYaml: |
      prometheusExporter:
        enabled: false
```

在该配置下：
- sonarqube 存储使用集群默认的 `StorageClass`。
- 默认使用 `PostgreSQL` 数据库来存储数据，使用集群默认的 `StorageClass`。
- 默认使用 `NodePort` 对外暴露 9000 端口。

### 3.3、生产环境

#### 3.3.1、证书配置

- 使用已有证书
  1. 在集群中创建 `Secret` 保存证书信息。
  2. 在 `valuesYaml` 配置项中增加如下证书配置。

```yaml
valuesYaml: |
  tls:
  # secret 名称
  - secretName: chart-example-tls
    hosts:
    # 证书对应的域名
    - chart-example.local
```

#### 3.3.2、存储配置

- 数据库配置（以 PostgreSQL 为例）
  1. 在外部配置高可用数据库。
  2. 在 `valuesYaml` 配置项中增加配置：

```yaml
valuesYaml: |
  postgresql:
    enabled: false
  jdbcOverwrite:
    enabled: true
    # PostgreSQL 数据库连接配置
    jdbcUrl: "jdbc:postgresql://myPostgress/myDatabase?socketTimeout=1500"
    jdbcUsername: "sonarUser"
    jdbcPassword: "sonarPass"
```

- SonarQube 存储配置
  1. 在集群中创建需要的 `StorageClass`。
  2. 在 `valuesYaml` 配置项中增加配置:

```yaml
valuesYaml: |
  persistence:
    enabled: true
    # 使用集群中创建的 storageclass 名称
    storageclass: prod-storageclass
    # 使用的磁盘大小
    size: 20g
```
