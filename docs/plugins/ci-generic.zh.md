# ci-generic 插件

这个插件可以基于本地或者远程的文件在 GitLab/GitHub 安装 CI 配置

## 用例

下面的配置文件展示的是"tool file"的内容。

关于更多关于DevStream的主配置、tool file、var file的信息，请阅读[核心概念概览](../core-concepts/overview.zh.md)和[DevStream配置](../core-concepts/config.zh.md).

``` yaml
--8<-- "ci-generic.yaml"
```

**注意事项：**

- `projectRepo` 配置字段用于表示代码仓库的配置信息，具体配置可查看[SCM配置项](./scm-option.zh.md)
- `ci.configContents` 和 `ci.configLocation` 不能同时为空。
- 如果你配置了 `projectRepo.scmType` 为 `github`，那 `ci.type` 就不能是 `gitlab-ci`。
- 如果你配置了 `projectRepo.scmType` 为 `gitlab`，那 `ci.type` 就不能是 `github-actions`。

## 示例

### 使用本地的 Workflows 目录

```yaml
tools:
- name: ci-generic
  instanceID: test-github
  options:
    ci:
      configLocation: workflows
      type: github
    projectRepo:
      owner: devstream
      org: ""
      name: test-repo
      branch: main
      scmType: github
```

这个配置将会把本地当前运行环境下的 workflows 目录放置于 GitHub 的 `.github/workflows` 目录。

### 使用 HTTP 获取远程的CI文件

```yaml
tools:
- name: ci-generic
  instanceID: test-gitlab
  options:
    ci:
      configLocation : https://raw.githubusercontent.com/DeekshithSN/Jenkinsfile/inputTest/Jenkinsfile
      type: jenkins
    projectRepo:
      owner: root
      org: ""
      name: test-repo
      branch: main
      scmType: gitlab
      baseURL: http://127.0.0.1:30000
```

这个配置将会把[URL](https://raw.githubusercontent.com/DeekshithSN/Jenkinsfile/inputTest/Jenkinsfile) 中的 Jenkinsfile 文件置于 GitLab 的仓库。

### 使用Github仓库中的CI文件
```yaml
tools:
- name: ci-generic
  instanceID: test-gitlab
  options:
    ci:
      configLocation : git@github.com:devstream-io/devstream.git//staging/dtm-jenkins-pipeline-example/general
      type: jenkins
    projectRepo:
      owner: root
      org: ""
      name: test-repo
      branch: main
      scmType: gitlab
      baseURL: http://127.0.0.1:30000
```

这个配置将会搜索[devstream 仓库](https://github.com/devstream-io/devstream)下的staging/dtm-jenkins-pipeline-example/general 目录，获取到目录下的 Jenkinsfile，置于 gitlab 仓库内。

### 在Devstream中直接配置CI文件

```yaml
tools:
- name: ci-generic
  instanceID: test-gitlab
  options:
    ci:
      configContents:
        pr.yaml: |
name: GitHub Actions Demo
run-name: ${{ github.actor }} is testing out GitHub Actions 🚀
on: [push]
jobs:
  Explore-GitHub-Actions:
    runs-on: ubuntu-latest
    steps:
      - run: echo "🎉 The job was automatically triggered by a ${{ github.event_name }} event."
      projectRepo:
        owner: test-user
        org: ""
        name: test-repo
        branch: main
        scmType: github
```

这个配置将会在用户的Github仓库`test-user/test-repo`下创建`.github/workflows/pr.yaml`文件。
