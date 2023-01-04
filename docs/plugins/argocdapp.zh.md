# argocdapp 插件

该插件会创建一个在 Kubernetes 上 [Argo CD Application](https://argo-cd.readthedocs.io/en/stable/core_concepts/) 的自定义资源。

**注意:**

- 在使用该插件前需要先安装 Argocd CD。你可以使用 [helm-installer 插件](./helm-installer/argocd.md) 来安装它。
- 目前该插件只支持 Helm chart 的配置方式。

## 用例

以下内容是该插件的示例配置文件。

```yaml
--8<-- "argocdapp.yaml"
```

### 自动创建 Helm 配置

如果你不想要自己创建 Helm 配置，该插件支持把 Devstream 提供的默认 Helm 配置上传到 `source.path` 的配置路径上，这样你就可以直接使用该插件。配置示例如下：

```yaml
---
tools:
- name: go-webapp-argocd-deploy
  plugin: argocdapp
  dependsOn: ["repo-scaffolding.golang-github"]
  options:
    app:
      name: hello
      namespace: argocd
    destination:
      server: https://kubernetes.default.svc
      namespace: default
    source:
      valuefile: values.yaml
      path: charts/go-hello-http
      repoURL: https://github.com/devstream-io/testrepo.git
    imageRepo:
      url: http://test.barbor.com/library
      user: test_owner
      tag: "1.0.0"
```

这个示例配置将会把 [Helm 配置](https://github.com/devstream-io/dtm-pipeline-templates/tree/main/argocdapp/helm) 上传到 [testrepo](https://github.com/devstream-io/testrepo.git) 仓库中，生成的 Helm 配置会使用 `http://test.barbor.com/library/test_owner/hello:1.0.0` 作为 Helm 应用的启动镜像。