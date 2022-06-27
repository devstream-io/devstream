# 测试

## 单元测试

运行所有的单元测试：

```shell
go test ./...
```

## 端到端测试

### 1 使用 GitHub Action 运行端到端测试

GitHub Actions 会在代码提交时自动运行端到端测试。

### 2 在本地运行端到端测试

```shell
bash hack/e2e/e2e-run.sh
```

请在测试前设置以下环境变量：

- GITHUB_USER
- GITHUB_TOKEN
- DOCKERHUB_USERNAME
- DOCKERHUB_TOKEN
