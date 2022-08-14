# Debug

Take `localstack-plugin` as an example.

## Prepare

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

## Debugging in VSCode

See [Go for Visual Studio Code](https://github.com/golang/vscode-go).

### Build Plugin

```bash
VERSION=0.8.0 DEV=true make build-plugin.localstack
```

### launch.json File

For `dtm apply -f config.yaml`.

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

### Debug Plugin

Through the above `laucnh.json` file, we can trace the code until `internal/pkg/pluginengine/plugin.go`.

- `p.Create(tool.Options)`
- `p.Read(tool.Options)`
- `p.Update(tool.Options)`
- `p.Delete(tool.Options)`

We can watch the value of `tool.Options`, but we can not trace the implementation of `p.CRUD` now.

Currently, the plugin can enter debug mode through a [Ginkgo](https://github.com/onsi/ginkgo) test suite for CRUD.

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
