# harbor 插件

这个插件使用 helm 在已有的 k8s 集群上安装 [harbor](https://goharbor.io/)。

## 使用方法

### 测试环境

如果你想在**本地测试插件**， 可以使用如下 `values_yaml` 配置。

```yaml
values_yaml: |
  expose:
    type: nodePort
    tls:
      enabled: false
  chartmuseum:
    enabled: false
  clair:
    enabled: false
  notary:
    enabled: false
  trivy:
    enabled: false
```

在该配置下
- helm 会自动创建依赖的 Postgresql 和 Redis；
- 数据挂载的磁盘默认会使用集群上机器的本地磁盘；
- 只安装 `harbor` 主程序而不会安装其余的插件；
- 通过 `NodePort` 对外暴露服务，可使用 `http://{{k8s 节点ip}}:30002` 域名来访问，默认账号名密码为 admin/Harbor12345 (生产环境请替换默认账号密码)。

### 生产环境

Harbor的大部分组件都是无状态的。因此我们可以通过增加 `Pod` 的副本来确保组件被部署到多个工作节点，并使用 `Kubernetes` 的 `Service` 机制来确保跨 `Pod` 的网络连接。

#### 外部存储

> Harbor 正常运行依赖 Postgresql 和 Redis。

- Postgresql：生产环境建议使用外部高可用的 Postgresql 集群，具体配置可参考 [config](https://github.com/goharbor/harbor-helm#configuration) 中的 Database 选项。
- Redis: 生产环境建议使用外部高可用的 Redis 集群，具体配置可参考 [config](https://github.com/goharbor/harbor-helm#configuration) 中的 Redis 选项。

#### 磁盘存储
请将配置中的 `storageClass` 修改为已存在的 StorageClass，集体配置可参考  [config](https://github.com/goharbor/harbor-helm#configuration) 中的 Persistence 选项。

#### 网络层配置
该插件支持 `Ingress`, `ClusterIP`, `NodePort`, `LoadBalancer` 对外暴露的模式，可以基于需求进行选择。

#### 证书配置

- 使用自签名证书
  1. 将 `tls.enabled` 设置为 `true`，并编辑对应的域名 `externalURL`；
  2. 将 Pod `harbor-core` 中目录 `/etc/core/ca` 存储的自签名证书复制到自己的本机；
  3. 在自己的主机上信任该证书。
- 使用公共证书
  1. 将公共证书添加为密钥 (`Secret`)；
  2. 将 `tls.enabled` 设置为 `true`，并编辑对应的域名 `externalURL`；
  3. 配置 `tls.secretName` 使用该公共证书。

### 配置

```yaml
--8<-- "harbor.yaml"
```

目前除了 `values_yaml` 字段，所有示例参数均为必填项。
