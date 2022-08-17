# harbor 插件

`harbor` 插件用于部署、管理 [Harbor](https://goharbor.io/) 应用。

Harbor 的主流部署方式有2种：**docker compose** 和 **helm**。
现在 DevStream 有2个插件 `harbor-docker` 和 `harbor` 分别支持这2种部署方式，但是目前以 helm 方式为主。
在不久的将来，这两个插件将会被合并成一个。

## 1、前置要求

- 有一个可用的 Kubernetes 集群，版本 1.10+

## 2、部署架构

Harbor 本身并不关注如何实现存储高可用和 PostgreSQL、Redis 的高可用。
所以 Harbor 支持通过 PVCs 的方式持久化数据，支持对接外部 PostgreSQL 和 Redis。

Harbor 部署架构整体如下图所示(图片来自 Harbor 官网)：

![Harbor Architecture](./harbor/ha.png)

## 3、开始部署

下文将介绍如何配置 `harbor` 插件，完成 Harbor 应用的部署。演示环境为一台装有以 minikube 方式部署的单节点 k8s 集群的 Linux 云主机。

### 3.1、快速开始

如果仅是用于开发、测试等目的，希望快速完成 Harbor 的部署，可以使用如下配置快速开始：

```yaml
--8<-- "harbor.yaml"
```

在成功执行 `dtm apply` 命令后，我们可以在 harbor 命名空间下看到下述主要资源：

- **Deployment** (`kubectl get deployment -n harbor`)

```shell
NAME                READY   UP-TO-DATE   AVAILABLE   AGE
harbor-core         1/1     1            1           2m56s
harbor-jobservice   1/1     1            1           2m56s
harbor-nginx        1/1     1            1           2m56s
harbor-portal       1/1     1            1           2m56s
harbor-registry     1/1     1            1           2m56s
```

可以看到几乎所有 Harbor 相关服务都是以 Deployment 方式在运行。

- **StatefulSet** (`kubectl get statefulset -n harbor`)

```shell
NAME              READY   AGE
harbor-database   1/1     3m40s
harbor-redis      1/1     3m40s
```

这两个 StatefulSet 资源对应的是 Harbor 所依赖的 PostgreSQL 和 Redis。
换言之，当前部署方式会自动完成 PostgreSQL 和 Redis 的部署，但是同时需要注意 PostgreSQL 和 Redis 并不是高可用的。

- **Service** (`kubectl get service -n harbor`)

```shell
NAME                TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)             AGE
harbor              NodePort    10.99.177.6      <none>        80:30002/TCP        4m17s
harbor-core         ClusterIP   10.106.220.239   <none>        80/TCP              4m17s
harbor-database     ClusterIP   10.102.102.95    <none>        5432/TCP            4m17s
harbor-jobservice   ClusterIP   10.98.5.49       <none>        80/TCP              4m17s
harbor-portal       ClusterIP   10.105.115.5     <none>        80/TCP              4m17s
harbor-redis        ClusterIP   10.104.100.167   <none>        6379/TCP            4m17s
harbor-registry     ClusterIP   10.106.124.148   <none>        5000/TCP,8080/TCP   4m17s
```

可以看到默认情况下 Harbor 通过 NodePort 方式暴露服务时，端口是 30002。

- **PersistentVolumeClaim** (`kubectl get pvc -n harbor`)

```shell
NAME                              STATUS   VOLUME                                     CAPACITY   ACCESS MODES   STORAGECLASS   AGE
data-harbor-redis-0               Bound    pvc-5b6b5eb4-c40d-4f46-8f19-ff3a8869e56f   1Gi        RWO            standard       5m12s
database-data-harbor-database-0   Bound    pvc-d7ccaf1f-c450-4a16-937a-f55ad0c7c18d   1Gi        RWO            standard       5m12s
harbor-jobservice                 Bound    pvc-9407ef73-eb65-4a56-8720-a9ddbcb76fef   1Gi        RWO            standard       5m13s
harbor-registry                   Bound    pvc-34a2b88d-9ff2-4af4-9faf-2b33e97b971f   5Gi        RWO            standard       5m13s
```

从这里可以看到 Harbor 所需要的存储卷有4个，其中也包括了 PostgreSQL 和 Redis 的存储。

- **PersistentVolume** (`kubectl get pv`)

```shell
pvc-34a2b88d-9ff2-4af4-9faf-2b33e97b971f   5Gi        RWO            Delete           Bound         harbor/harbor-registry                    standard                5m22s
pvc-5b6b5eb4-c40d-4f46-8f19-ff3a8869e56f   1Gi        RWO            Delete           Bound         harbor/data-harbor-redis-0                standard                5m22s
pvc-9407ef73-eb65-4a56-8720-a9ddbcb76fef   1Gi        RWO            Delete           Bound         harbor/harbor-jobservice                  standard                5m22s
pvc-d7ccaf1f-c450-4a16-937a-f55ad0c7c18d   1Gi        RWO            Delete           Bound         harbor/database-data-harbor-database-0    standard                5m22s
```

我们并没有配置 StorageClass，所以这里用的是集群内的 default StorageClass 完成的 pv 创建。

在我的环境里，default StorageClass 如下(`kubectl get storageclass`)：

```shell
NAME                 PROVISIONER                RECLAIMPOLICY   VOLUMEBINDINGMODE   ALLOWVOLUMEEXPANSION   AGE
standard (default)   k8s.io/minikube-hostpath   Delete          Immediate           false                  20h
```

这个 StorageClass 对应的 Provisioner 会以 hostPath 的方式提供 pv。

到这里，我们就可以通过 http://127.0.0.1:3002 访问到 Harbor 登录页面了，如下：

![Harbor Login](./harbor/login.png)

默认登录账号/密码是 `admin/Harbor12345`。登录后，可以看到默认首页如下：

![Harbor Dashboard](./harbor/dashboard.png)

如果是在云主机上部署的 Harbor，可以通过 `kubectl port-forward` 命令来暴露服务：

```shell
ip=YOUR_HOST_IP
kubectl port-forward -n harbor service/harbor --address=${ip} 80
```

需要注意的是这里得使用主机真实网卡 ip，而我们在浏览器上输入的 ip 是云主机的公网 ip，两者并不一样。

### 3.2、最小化配置

`harbor` 插件的配置项默认值如下表：

| key                | default value            | description                         |
| ----               | ----                     | ----                                |
| chart.chart_name   | harbor/harbor            | helm chart 包名称                    |
| chart.timeout      | 10m                      | helm install 的超时时间               |
| chart.upgradeCRDs  | true                     | 是否更新 CRDs（如果有）                 |
| chart.release_name | harbor                   | helm 发布名称                         |
| chart.wait         | true                     | 是否等待部署完成                       |
| chart.namespace    | harbor                   | 部署的命名空间                         |
| repo.url           | https://helm.goharbor.io | helm 仓库地址                         |
| repo.name          | harbor                   | helm 仓库名                           |
| create_namespace   | false                    | 是否需要 dtm 来新建命名空间              |

目前除了 `values_yaml` 字段和默认配置，其它所有示例参数均为必填项。 因此，我们也可以使用如下最小化配置达到和上面配置文件一样的效果：

```yaml
tools:
- name: harbor
  instanceID: default
  dependsOn: [ ]
  options:
    create_namespace: true
    chart:
      values_yaml: |
        externalURL: http://127.0.0.1
        expose:
          type: nodePort
          tls:
            enabled: false
        chartmuseum:
          enabled: false
        notary:
          enabled: false
        trivy:
          enabled: false
```

### 3.3、持久化存储数据

前面我们已经看到了如果不配置 StorageClass，Harbor 会使用集群内的 default StorageClass。比如在当前演示环境下，harbor-registry 所使用的 pv 配置大致如下：

```yaml
apiVersion: v1
kind: PersistentVolume
metadata:
  name: pvc-34a2b88d-9ff2-4af4-9faf-2b33e97b971f
spec:
  accessModes:
  - ReadWriteOnce
  capacity:
    storage: 5Gi
  claimRef:
    apiVersion: v1
    kind: PersistentVolumeClaim
    name: harbor-registry
    namespace: harbor
  hostPath:
    path: /tmp/hostpath-provisioner/harbor/harbor-registry
  persistentVolumeReclaimPolicy: Delete
  storageClassName: standard
  volumeMode: Filesystem
status:
  phase: Bound
```

可见数据其实挂在到了主机的 `/tmp/hostpath-provisioner/harbor/harbor-registry` 目录下了。

Harbor 支持3种持久化存储数据配置方式：

1. 配置 StorageClass（默认使用 default StorageClass）；
2. 使用已经存在的 pvc（手动创建）；
3. 对接 azure、gcs、s3、swift、oss 等实现镜像和 Charts 云端存储。

我们暂时只介绍第一种方式，也就是指定 StorageClass 实现存储数据持久化。
registry、jobservice、chartmuseum、database、redis、trivy 等组件都可以单独指定 StorageClass。
假设我们现在有一个新的 StorageClass 叫做 nfs，这时候要实现前面部署的 Harbor 所有数据持久化到外部 pv，我们可以这样配置：

```yaml
tools:
- name: harbor
  instanceID: default
  dependsOn: [ ]
  options:
    create_namespace: true
    chart:
      values_yaml: |
        persistence:
          persistentVolumeClaim:
            registry:
              storageClass: "nfs"
              accessMode: ReadWriteOnce
              size: 5Gi
            jobservice:
              storageClass: "nfs"
              accessMode: ReadWriteOnce
              size: 1Gi
            database:
              storageClass: "nfs"
              accessMode: ReadWriteOnce
              size: 1Gi
            redis:
              storageClass: "nfs"
              accessMode: ReadWriteOnce
              size: 1Gi
```

### 3.4、服务暴露

Harbor 可以以 ClusterIP、LoadBalancer、NodePort 和 Ingress 等方式对外暴露服务。我们前面使用的就是 NodePort 方式：

```yaml
tools:
- name: harbor
  instanceID: default
  dependsOn: [ ]
  options:
    create_namespace: true
    chart:
      values_yaml: |
        externalURL: http://127.0.0.1
        expose:
          type: nodePort
```

本文我们再介绍一下如何使用 Ingress 方式暴露服务：

```yaml
tools:
- name: harbor
  instanceID: default
  dependsOn: [ ]
  options:
    create_namespace: true
    chart:
      values_yaml: |
        externalURL: http://127.0.0.1
        expose:
          type: ingress
          tls:
            enabled: false
          hosts:
            core: "core.harbor.domain"
```

注意：如果没有开启 TLS，这种方式暴露 Harbor 服务后 docker push/pull 命令必须带上端口。

### 3.5、PostgreSQL 和 Redis 高可用

Harbor 依赖 PostgreSQL 和 Redis 服务，默认情况下自动部署的 PostgreSQL 和 Redis 服务都是非高可用模式。
换言之如果我们单独部署高可用的 PostgreSQL 和 Redis 服务，Harbor 是支持去对接外部 PostgreSQL 和 Redis 服务的。

TODO(daniel-hutao): 本节待细化

### 3.6、https 配置

TODO(daniel-hutao): 本节待细化

- 使用自签名证书
  1. 将 `tls.enabled` 设置为 `true`，并编辑对应的域名 `externalURL`；
  2. 将 Pod `harbor-core` 中目录 `/etc/core/ca` 存储的自签名证书复制到自己的本机；
  3. 在自己的主机上信任该证书。
- 使用公共证书
  1. 将公共证书添加为密钥 (`Secret`)；
  2. 将 `tls.enabled` 设置为 `true`，并编辑对应的域名 `externalURL`；
  3. 配置 `tls.secretName` 使用该公共证书。

## 4、典型场景

// TODO(daniel-hutao): 本节待补充

### 4.1、HTTP

#### 4.1.1 HTTP + Registry + Internal Database + Internal Redis

#### 4.1.2 HTTP + Registry + Chartmuseum + Internal Database + Internal Redis

#### 4.1.3 HTTP + Registry + Chartmuseum + External Database + External Redis

### 4.2、HTTPS

#### 4.1.1 HTTPS + Registry + Internal Database + Internal Redis

#### 4.1.2 HTTPS + Registry + Chartmuseum + Internal Database + Internal Redis

#### 4.1.3 HTTPS + Registry + Chartmuseum + External Database + External Redis
