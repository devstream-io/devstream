# repo-scaffolding 插件

这个插件会基于一个脚手架仓库来初始化一个 Gihub 或者 GitLab 仓库。

## 用例

下面的配置文件展示的是"tool file"的内容。

关于更多关于DevStream的主配置、tool file、var file的信息，请阅读[核心概念概览](../core-concepts/overview.zh.md)和[DevStream配置](../core-concepts/config.zh.md).

```yaml
--8<-- "repo-scaffolding.yaml"
```

**注意:**

- 如果你执行 `dtm delete` 命令，这个仓库将会被删除。
- 如果你执行 `dtm update` 命令,  这个仓库将会被删除然后重新创建。
- 对于 `repo-scaffolding` 插件，目前只需要 token 有 `repo`, `delete_repo` 权限即可。
- `destinationRepo` 配置字段用于表示代码仓库的配置信息，具体配置可查看[SCM配置项](./scm-option.zh.md)。
- `sourceRepo` 配置字段用于表示代码仓库的配置信息，具体配置可查看[SCM配置项](./scm-option.zh.md)。
- `destinationRepo` 如果是 `gitlab`, 则支持配置 `destinationRepo.visibility`，此配置用于设置新建仓库的权限，支持的选项有 `public`, `private` 和 `internal`。

### 变量

这个配置用于设置渲染源脚手架仓库时的变量，以下变量会默认设置：

```json
{
    "AppName": destinationRepo.repo,
    "Repo": {
        "Name": destinationRepo.repo,
        "Owner": destinationRepo.owner
    }
}
```

## Outputs

这个插件有以下三个输出：

- `owner`
- `repo`
- `repoURL`

## 示例

### 官方支持脚手架项目

以下仓库是用于在 `sourceRepo` 设置的官方脚手架仓库，你可以使用这些仓库或者创建自己的脚手架仓库。

| language    | org           | repo                       |
|-------------|---------------|----------------------------|
| Golang      | devstream-io  | dtm-scaffolding-golang     |
| Golang      | devstream-io  | dtm-scaffolding-golang-cli |
| Java Spring | spring-guides | gs-spring-boot             |
| Python      | devstream-io  | dtm-repo-scaffolding-python-flask   |

### Golang

```yaml
tools:
  - name: repo-scaffolding
    instanceID: golang-scaffolding
    options:
      destinationRepo:
        owner: test_owner
        org: ""
        name: dtm-test-golang
        branch: main
        scmType: github
      sourceRepo:
        org: devstream-io
        name: dtm-scaffolding-golang
        scmType: github
      vars:
        ImageRepo: dtm-test/golang-repo
```

这个配置在 GitHub 为用于 test_owner 创建 `dtm-test-golang` 仓库，它的生成是基于 `devstream-io/dtm-scaffolding-golang` 官方 Golang 脚手架仓库和传入的变量 `ImageRepo`。

### Golang CLI

```yaml
tools:
  - name: repo-scaffolding
    instanceID: golang-cli-scaffolding
    options:
      destinationRepo:
        owner: test_owner
        org: ""
        name: dtm-test-golang-cli
        branch: main
        scmType: github
      sourceRepo:
        org: devstream-io
        name: dtm-scaffolding-golang-cli
        scmType: github
```

这个配置在 GitHub 为用于 test_owner 创建 `dtm-test-golang-cli` 仓库，它的生成是基于 `devstream-io/dtm-scaffolding-golang-cli` 官方 Golang CLI 脚手架仓库。

### Java Spring

```yaml
tools:
  - name: repo-scaffolding
    instanceID: java-scaffolding
    options:
      destinationRepo:
        owner: test_owner
        org: ""
        name: dtm-test-java
        branch: main
        baseUrl: 127.0.0.1:30001
        visibility: public
        scmType: gitlab
      sourceRepo:
        org: spring-guides
        name: gs-spring-boot
        scmType: github
```

这个配置会在 GitLab 为用户 test_owner 创建 `dtm-test-java` 仓库，使用的是 Spring 官方的 `spring-guides/gs-spring-boot` 仓库。
