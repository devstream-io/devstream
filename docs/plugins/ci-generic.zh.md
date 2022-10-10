# ci-generic 插件

这个插件可以基于本地或者远程的文件在 GitLab/GitHub 安装 CI 配置

## 使用

下面的配置文件展示的是"tool file"的内容。

关于更多关于DevStream的主配置、tool file、var file的信息，请阅读[核心概念概览](../core-concepts/core-concepts.zh.md)和[DevStream配置](../core-concepts/config.zh.md).

``` yaml
--8<-- "ci-generic.yaml"
```

### 字段配置

| key                    | description                                                                                                      |
| ----                   | ----                                                                                                             |
| ci.localPath           | If your ci file is local, you can set the this field to the ci file location, which can be a directory or a file |
| ci.remoteURL           | If your ci file is remote, you can set this field to url address                                                 |
| ci.type                | ci type, support gitlab, github, jenkins for now                                                                 |
| projectRepo.owner      | destination repo owner                                                                                           |
| projectRepo.org        | destination repo org                                                                                             |
| projectRepo.repo       | destination repo name                                                                                            |
| projectRepo.branch     | destination repo branch                                                                                          |
| projectRepo.repoType  | destination repo type, support github/gitlab for now                                                             |
| projectRepo.baseURL   | if you use self-build gitlab, you can set this field to gitlab address                                           |
| projectRepo.visibility | if you use gitlab, you can set this field for repo permission                                                    |

**注意事项：**

- `ci.localPath` 和 `ci.remoteURL` 不能同时为空。
- 如果你同时设置了 `ci.localPath` 和 `ci.remoteURL`，`ci.localPath` 将会被优先使用。
- 如果你的 `projectRepo.repoType` 配置是 `gitlab`，`ci.type` 就不能是 `github`。
- 如果你的 `projectRepo.repoType` 配置是 `github`, `ci.type` 就不能是 `gitlab`。

### 示例

#### 本地的 Workflows 目录部署到 GitHub 仓库

```yaml
tools:
  - name: ci-generic
    instanceID: test-github
    options:
      ci:
        localPath: workflows
        type: github
      projectRepo:
        owner: devstream
        org: ""
        repo: test-repo
        branch: main
        repoType: github
```

这个配置将会把本地的 workflows 目录放置于 GitHub 的 `.github/workflows` 目录。

#### Remote Jenkinsfile With Gitlab

```yaml
tools:
  - name: ci-generic
    instanceID: test-gitlab
    options:
      ci:
        remoteURL : https://raw.githubusercontent.com/DeekshithSN/Jenkinsfile/inputTest/Jenkinsfile
        type: jenkins
      projectRepo:
        owner: root
        org: ""
        repo: test-repo
        branch: main
        repoType: gitlab
        baseURL: http://127.0.0.1:30000
```

这个配置将会把[URL](https://raw.githubusercontent.com/DeekshithSN/Jenkinsfile/inputTest/Jenkinsfile) 中的 Jenkinsfile 文件置于 GitLab 的仓库。
