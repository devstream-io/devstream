# helm-installer 插件

`helm-installer` 插件实现了比 `helm` 更加简单和容易上手的方式来快速部署提供了 Helm Chart 的应用。

## 1、快速开始

只需要一个最小化配置，你就可以快速使用默认配置部署一个 Helm Chart。你可以将如下配置内容保存到本地 config.yaml 文件中：

```yaml
config:
  state:
    backend: local
    options:
      stateFile: devstream.state
tools:
- name: helm-installer
  instanceID: argocd-001
```

在这个配置文件里，和插件相关的配置 name 和 instanceID，前者表示你将使用 `helm-installer` 插件，后者表示插件实例名。
请注意这个 instanceID 使用了 "argocd-" 前缀，DevStream 会识别这个前缀，尝试寻找 Argo CD 应用对应的 Chart，并设置一系列默认值，然后开始部署。

你可以在 [Install Argo CD with DevStream](./argocd.zh.md) 中查看 DevStream 为你设置了哪些默认值。

接着执行如下命令开始部署：

```sh
./dtm init -f config.yaml
./dtm apply -f config.yaml -y
```

## 2、插件介绍

`helm-installer` 插件的完整配置格式如下：

```yaml
tools:
- name: helm-installer
  instanceID: argocd-001
  dependsOn: [ ]
  options:
    repo:
      name: ""
      url: ""
    chart:
      chartPath: ""
      chartName: ""
      version: ""
      namespace: ""
      wait: true
      timeout: 10m
      upgradeCRDs: true
    valuesYaml: ""
```

这里有一些细节需要注意，下述几个小节将详细为你介绍。

### 2.1、instanceID 使用技巧

instanceID 的前缀如果能够匹配到某个已经被支持的工具（详见文末列表），那么 DevStream 会为你设置一系列的默认值。
比如 "argocd-001" 的前缀 "argocd-" 能够匹配到 "argocd" + "-"，因此 Argo CD 的默认 Chart 配置会被应用，于是如下最小化配置：

```yaml
tools:
- name: helm-installer
  instanceID: argocd-001
```

将会被 DevStream 直接补全成：

```yaml
- name: helm-installer
  instanceID: argocd-001
  dependsOn: [ ]
  options:
    repo:
      name: ""
      url: ""
    chart:
      chartPath: ""
      chartName: ""
      version: ""
      namespace: ""
      wait: true
      timeout: 10m
      upgradeCRDs: true
    valuesYaml: ""
```

### 2.2、自定义 Chart 配置

如果你想使用自定义 Chart 的 values.yaml 配置，只需要将 values.yaml 的文件路径或者内容直接加到 helm-installer 插件配置 options 部分的 chart.valuesYaml 里。
两种配置方式分别如下：

- 使用本地 values.yaml 文件

```yaml
- name: helm-installer
  instanceID: argocd-001
  dependsOn: [ ]
  options:
    valuesYaml: "./values.yaml"
```

- 直接使用 values.yaml 文件内容

```yaml
- name: helm-installer
  instanceID: argocd-001
  dependsOn: [ ]
  options:
    valuesYaml: |
      foo: bar
```

## 3、当前支持的工具列表

当前 DevStream 支持使用"极简配置"部署如下应用（也就是能够根据 instanceID 配置识别到 Chart 地址等信息，并设置一系列默认值，直接开始部署流程）：

- [Install Argo CD with DevStream](./argocd.zh.md)
- [Install Artifactory with DevStream](./artifactory.zh.md)
- [Install DevLake with DevStream](./devlake.zh.md)
- [Install Harbor with DevStream](./harbor.zh.md)
- [Install Jenkins with DevStream](./jenkins.zh.md)
- [Install Kube Prometheus with DevStream](./kube-prometheus.zh.md)
- [Install OpenLDAP with DevStream](./openldap.zh.md)
- [Install SonarQube with DevStream](./sonarqube.zh.md)
- [Install Tekton with DevStream](./tekton.zh.md)
- [Install Vault with DevStream](./vault.zh.md)
