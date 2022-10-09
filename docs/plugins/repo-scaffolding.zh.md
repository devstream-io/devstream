# repo-scaffolding 插件

这个插件会基于一个脚手架仓库来初始化一个 Gihub 或者 GitLab 仓库。

## 运行需求

这个插件基于你使用的代码仓库类型需要设置以下配置：

### GitHub

- GITHUB_TOKEN: 在使用插件之前请先设置这个环境变量，如果你不知道如何获取这个 token，可以查看文档 [Creating a personal access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)。

### GitLab

- GITLAB_TOKEN： 在使用插件之前请先设置这个环境变量，如果你不知道如何获取这个 token，可以查看文档 [Personal access tokens](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html)。
- `destinationRepo.baseUrl`：如果你使用的是自建的 GitLab 仓库，需要将该配置设置为自建 GItLab 的 URL 地址。
- `destinationRepo.visibility`：此配置用于设置新建仓库的权限，支持的选项有 `public`, `private` 和 `internal`。

*注意：*

- 如果你执行 `dtm delete` 命令，这个仓库将会被删除。

- 如果你执行 `dtm update` 命令,  这个仓库将会被删除然后重新创建。 

- 对于 `repo-scaffolding` 插件，目前只需要 token 有 `repo`, `delete_repo` 权限即可。 

## 使用方法

下面的配置文件展示的是"tool file"的内容。

关于更多关于DevStream的主配置、tool file、var file的信息，请阅读[核心概念概览](../core-concepts/core-concepts.zh.md)和[DevStream配置](../core-concepts/config.zh.md).

```yaml
--8<-- "repo-scaffolding.yaml"
```

**请注意这里的设置参数都是大小写敏感的**

在配置文件中替换以下配置：

### destinationRepo

这个是目标仓库的配置，包括以下几个配置项：

- `YOUR_DESTINATION_USERNAME`
- `YOUR_DESTINATION_ORGANIZATION_NAME`
- `YOUR_DESTINATION_REPO_NAME`
- `YOUR_DESTINATION_REPO_MAIN_BRANCH`
- `YOUR_DESTINATION_REPO_TYPE`

`owner`，`org` 和 `repo` 目前是必填的，`branch` 的默认值是  "main"，`repoType` 配置目前支持 `gitlab` 和 `github`。

### sourceRepo

这个是源脚手架仓库的配置（目前只支持 Github），包括以下几个配置：

- `YOUR_TEMPLATE_REPO_ORG`
- `YOUR_TEMPLATE_REPO_NAME`
- `YOUR_TEMPLATE_REPO_TYPE`

目前这两个配置项都是必填的，`repoType` 配置目前只支持 `github`。

### vars

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

## 示例 

### 官方支持脚手架项目

以下仓库是用于在 `sourceRepo` 设置的官方脚手架仓库，你可以使用这些仓库或者创建自己的脚手架仓库。

| language | org | repo |
|  ----  | ----  |----  |
| Golang | devstream-io | dtm-scaffolding-golang |
| Java Spring | spring-guides | gs-spring-boot |


### Golang

```yaml
tools:
  - name: repo-scaffolding
    instanceID: golang-scaffolding
    options:
      destinationRepo:
        owner: test_owner
        org: ""
        repo: dtm-test-golang
        branch: main
        repoType: github
      sourceRepo:
        org: devstream-io
        repo: dtm-scaffolding-golang
        repoType: github
      vars:
        ImageRepo: dtm-test/golang-repo
```

这个配置在 GitHub 为用于 test_owner 创建 `dtm-test-golang` 仓库，它的生成是基于 `devstream-io/dtm-scaffolding-golang` 官方 Golang 脚手架仓库和传入的变量 `ImageRepo`。

### Java Spring

```yaml
tools:
  - name: repo-scaffolding
    instanceID: java-scaffolding
    options:
      destinationRepo:
        owner: test_owner
        org: ""
        repo: dtm-test-java
        branch: main
        baseUrl: 127.0.0.1:30001
        visibility: public
        repoType: gitlab
      sourceRepo:
        org: spring-guides
        repo: gs-spring-boot
        repoType: github
```

这个配置会在 GitLab 为用户 test_owner 创建 `dtm-test-java` 仓库，使用的是 Spring 官方的 `spring-guides/gs-spring-boot` 仓库。

## Outputs

这个插件有以下三个输出：

- `owner`
- `repo`
- `repoURL`
