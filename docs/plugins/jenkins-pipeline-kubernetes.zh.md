# jenkins-pipeline-kubernetes 插件

这个插件在已有的 Jenkins 上建立 Jenkins job, 将 GitHub 作为 SCM。

步骤：

1. 按需修改配置项，其中 `githubRepoUrl` 为 GitHub 仓库地址，应预先建立一个 GitHub 仓库，并创建一个名为 "Jenkinsfile" 的文件放至仓库根目录。
2. 设置环境变量
    - `GITHUB_TOKEN`
    - `JENKINS_PASSWORD`

## 用例

```yaml

--8<-- "jenkins-pipeline-kubernetes.yaml"

```
