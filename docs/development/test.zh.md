# 测试

## 单元测试

运行所有的单元测试：

```shell
go test ./...
```

_注意：目前，不是所有的测试都是真正的“单元"测试，因为有些测试依赖于网络。热切期望你的帮助 : )_

## 端到端(End-to-End)测试

### GitHub Actions

当代码库的主分支代码更新时，GitHub Actions 会自动运行端到端测试。

GitHub Action 工作流程的定义在[这里](https://github.com/devstream-io/devstream/blob/main/.github/workflows/e2e-test.yml)，运行端到端测试时使用的 `dtm` 配置文件在[这里](https://github.com/devstream-io/devstream/tree/main/test/e2e/yaml)。

### 本地运行端到端测试

目前，我们编写了针对以下插件的简单端到端测试：

- `github-repo-scaffolding-golang`
- `githubactions-golang`
- `argocd`
- `argocdapp`

本地运行端到端测试的配置模板在[这里](https://github.com/devstream-io/devstream/blob/main/test/e2e/yaml/e2e-test-local.yaml)。

在测试前，请先确保 Docker 已经启动，并设置以下环境变量：

- GITHUB_USER
- GITHUB_TOKEN
- DOCKERHUB_USERNAME
- DOCKERHUB_TOKEN

然后执行：

```shell
bash hack/e2e/e2e-run.sh
```

测试脚本将会下载 kind/kubectl，启动一个 K8s 容器，并执行 `dtm` 命令，检查运行结果，最后清理环境。
