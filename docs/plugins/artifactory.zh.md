# artifactory 插件

这个插件使用 helm 在已有的 k8s 集群上安装 [artifactory](https://jfrog.com/artifactory/)。

## 使用方法

### 测试环境

如果你想在**本地测试插件**， 可以使用如下 `valuesYaml` 配置。

```yaml
valuesYaml: |
  artifactory:
    service:
      type: NodePort
    nodePort: 30002
  nginx:
    enabled: false
```

在该配置下

- helm 会自动创建依赖的 Postgresql；
- 数据挂载的磁盘默认会使用集群上机器的本地磁盘；
- 通过 `NodePort` 对外暴露服务，可使用 `http://{{k8s 节点ip}}:30002` 域名来访问，默认账号名密码为 admin/password (生产环境请替换默认账号密码)。

### 生产环境

#### 外部存储

- PostgreSQL：设置  `database.url` 来设置数据库地址，具体配置可参考 [Config](https://www.jfrog.com/confluence/display/JFROG/Configuring+the+Database) 中的选项。

#### 磁盘存储

可以设置 `customVolumes` 和 `customVolumeMounts` 来配置挂载磁盘，具体配置可参考  [Config](https://www.jfrog.com/confluence/display/JFROG/Configuring+the+Filestore)。

#### 网络层配置

该插件支持 `Ingress`, `ClusterIP`, `NodePort`, `LoadBalancer` 对外暴露的模式，可以基于需求进行选择。

### 配置

下面的配置文件展示的是"tool file"的内容。

关于更多关于DevStream的主配置、tool file、var file的信息，请阅读[核心概念概览](../core-concepts/core-concepts.zh.md)和[DevStream配置](../core-concepts/config.zh.md).

```yaml
--8<-- "artifactory.yaml"
```

#### 默认配置

| key                | default value           | description                                        |
| ----               | ----                    | ----                                               |
| chart.chartPath    | ""                      | 本地 chart 包路径                                   |
| chart.chartName    | jfrog/artifactory       | helm 包名称                                        |
| chart.timeout      | 10m                     | 等待部署成功的时间                                   |
| chart.upgradeCRDs  | true                    | 默认更新 CRD 配置（如果存在的话）                      |
| chart.releaseName  | artifactory             | helm 发布名称                                      |
| chart.wait         | true                    | 是否等待部署完成                                    |
| chart.namespace    | artifactory             | helm 部署的命名空间名称                              |
| repo.url           | https://charts.jfrog.io | helm 官方仓库地址                                   |
| repo.name          | jfrog                   | helm 仓库名                                        |

目前除了 `valuesYaml` 字段和默认配置，其它所有示例参数均为必填项。
