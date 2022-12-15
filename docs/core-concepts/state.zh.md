# 状态(State)

## 1 概念

状态记录了 DevStream 定义和创建的 DevOps 工具链和平台的当前状态。

状态包含了所有组件的配置和它们对应的状态，这样 DevStream 核心模块就可以依靠它计算出，达到配置中定义的状态所需要的操作。

我们可以在配置中指定使用哪种 `backend` 来存储状态，目前支持的 `backend` 有：

- local
- k8s
- s3

## 2 配置方式

配置中的 `config.state` 部分指定了如何存储 DevStream 状态。

### 2.1 本地文件

```yaml
config:
  state:
    backend: local
    options:
      stateFile: devstream.state # 可选，默认为 devstream.state
```

### 2.2 Kubernetes

```yaml
config:
  state:
    backend: k8s
    options:
      namespace: devstream # 可选, 默认是 "devstream", 不存在则自动创建
      configmap: state     # 可选, 默认是 "state", 不存在则自动创建
```

### 2.3 S3

```yaml
config:
  state:
    backend: s3
    options:
      bucket: devstream-remote-state
      region: ap-southeast-1
      key: devstream.state
```

_注意：`options` 下的 `bucket`、`region` 和 `key` 是 s3 后端的必填字段。_
