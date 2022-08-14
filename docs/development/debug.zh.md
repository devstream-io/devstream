# 调试

以 `localstack-plugin` 为例。

## 准备

`config.yaml`:

```yaml
---
varFile: ""
toolFile: ""
state:
  backend: local
  options:
    stateFile: devstream.state
---
tools:
- name: localstack
  instanceID: default
  dependsOn: [ ]
  options:
    create_namespace: true
    repo:
      name: localstack-charts
      url: https://localstack.github.io/helm-charts
    chart:
      chart_name: localstack-charts/localstack
      release_name: localstack
      namespace: localstack
      wait: true
      timeout: 5m
      upgradeCRDs: true
      values_yaml: |
        debug: true
        updateStrategy:
          type: Recreate
```

## 使用 VSCode 调试

查阅 [Go for Visual Studio Code](https://github.com/golang/vscode-go).

### 构建插件

```bash
VERSION=0.8.0 DEV=true make build-plugin.localstack
```

### launch.json 文件

调试 `dtm apply -f config.yaml`.

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "dtm-apply",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceRoot}/cmd/devstream",
            "args": [
                "apply",
                "-f",
                "${workspaceFolder}/config.yaml",
                "-d",
                "${workspaceFolder}/.devstream",
                "-y",
                "true"
            ],
            "showLog": true,
            "logOutput": "dap",
            "stopOnEntry": true,
            "buildFlags": "-ldflags='-X github.com/devstream-io/devstream/internal/pkg/version.Version=0.8.0'"
        }
    ]
}
```

### 调试插件

通过上面的 `laucnh.json` 文件, 我们可以跟踪代码直到 `internal/pkg/pluginengine/plugin.go`.

- `p.Create(tool.Options)`
- `p.Read(tool.Options)`
- `p.Update(tool.Options)`
- `p.Delete(tool.Options)`

我们能够查看 `tool.Options` 的值, 现在还不能够直接跟踪 `p.CRUD` 的实现代码.

目前, 一个插件可以通过 [Ginkgo](https://github.com/onsi/ginkgo) CRUD 测试套件进入调试模式。

```go
package localstack

import (
    "testing"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

func TestLocalStackPlugin(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "LocalStack Plugin Suite")
}
```
