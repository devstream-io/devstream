# harbor 插件

这个插件使用 helm 在已有的 k8s 集群上安装 [harbor](https://goharbor.io/)

## 使用方法

### 生产环境
#### 存储配置
##### 外部存储
- Postgresql：生产环境建议使用外部高可用的 Postgresql 集群，具体配置可参考 [config](https://github.com/goharbor/harbor-helm#configuration) 中的 Database 选项
- Redis: 生产环境建议使用外部高可用的 Redis 集群，具体配置可参考 [config](https://github.com/goharbor/harbor-helm#configuration) 中的 Redis 选项

##### 磁盘存储
请将配置中的 `storageClass` 修改为已存在的 StorageClass，集体配置可参考  [config](https://github.com/goharbor/harbor-helm#configuration) 中的 Persistence 选项

#### 网络层配置
该插件支持 `Ingress`, `ClusterIP`, `NodePort`, `LoadBalancer` 对外暴露的模式，可以基于需求进行选择

### 测试环境

如果你想**在本地测试插件**， 可以直接使用默认的参数配置，在该配置下

- helm 会自动创建依赖的 Postgresql 和 Redis
- 数据挂载的磁盘默认会使用集群上机器的本地磁盘
- helm 会自动创建 `Service` 指向启动的 harbor 服务，可以使用 `kubectl proxy` 来本地访问这些服务

### 配置

```yaml
tools:
# name of the tool
- name: harbor
  # id of the tool instance
  instanceID: default
  # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ ]
  # options for the plugin
  options:
    create_namespace: true
    repo:
      name: harbor
      # url of the Helm repo, use self host helm config beacuse official helm does'nt support namespace config
      url: https://helm.goharbor.io
    # Helm chart information
    chart:
      # name of the chart
      chart_name: harbor/harbor
      # k8s namespace where Tekton will be installed
      namespace: harbor
      # release name of the chart
      release_name: harbor
      # whether to wait for the release to be deployed or not
      wait: true
      # the time to wait for any individual Kubernetes operation (like Jobs for hooks). This defaults to 5m0s
      timeout: 10m
      # whether to perform a CRD upgrade during installation
      upgradeCRDs: true
      values_yaml: |
        notary.enabled: false
```

当前，所有配置项均为必填。