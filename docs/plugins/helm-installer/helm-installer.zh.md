# helm-installer 插件

`helm-installer` 插件实现了比 `helm` 更加简单和容易上手的方式来快速部署提供了 Helm Chart 的应用。

## 快速开始

只需要一个最小化配置，你就可以快速使用默认配置部署一个 Helm Chart。你可以将如下配置内容保存到本地 config.yaml 文件中：

```yaml
---
varFile: ""
toolFile: ""
pluginDir: ""
state:
  backend: local
  options:
    stateFile: devstream.state

---
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

## DevStream vs Helm

// TODO(daniel-hutao): add document here later.

## 当前支持的工具

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
