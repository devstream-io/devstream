# 插件列表


| Type                   | Plugin              | Note                            | Usage/Doc                               |
| ---------------------- | ------------------- | ------------------------------- | --------------------------------------- |
| Issue Tracking         | trello-github-integ | Trello/GitHub 整合              | [doc](trello-github-integ.md)           |
| Issue Tracking         | trello              | Trello 配置                     | [doc](trello.md)                        |
| Issue Tracking         | jira-github-integ   | Jira/GitHub 整合                | [doc](jira-github-integ.md)             |
| Issue Tracking         | zentao              | Zentao 安装                     | [doc](zentao.md)                        |
| Source Code Management | repo-scaffolding    | 应用仓库脚手架                  | [doc](repo-scaffolding.md)              |
| Source Code Management | gitlab-ce-docker    | 使用 docker 安装 GitLab CE 版本 | [doc](gitlab-ce-docker.md)              |
| CI                     | jenkins-pipeline    | 创建 Jenkins pipeline           | [doc](jenkins-pipeline.md)              |
| CI                     | github-actions      | 创建 GitHub Actions             | [doc](github-actions.md)                |
| CI                     | gitlab-ci           | 创建 GitLab CI                  | [doc](gitlab-ci.md)                     |
| CI                     | ci-generic          | 通用 CI 插件                    | [doc](ci-generic.md)                    |
| CD/GitOps              | argocdapp           | 创建 Argo CD 应用               | [doc](argocdapp.md)                     |
| Image Repository       | harbor-docker       | 使用 Docker Compose 安装 Harbor | [doc](harbor-docker.md)                 |
| Deployment             | helm-installer      | 使用 Helm 安装工具              | [doc](helm-installer/helm-installer.md) |

你也可以通过执行以下命令来获取当前的插件列表：

```shell
$ dtm list plugins
argocdapp
ci-generic
devlake-config
github-actions
gitlab-ce-docker
gitlab-ci
harbor-docker
helm-installer
jenkins-pipeline
jira-github-integ
repo-scaffolding
trello
trello-github-integ
zentao
```
