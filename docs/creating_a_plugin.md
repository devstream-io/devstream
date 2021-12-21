## 1 Interfaces

### 1.1 Definition

Each plugin needs to satisfy all the interfaces defined in the [plugin engine](https://github.com/merico-dev/stream/blob/main/internal/pkg/pluginengine/plugin.go#L12).

At the moment, there are only three interfaces, which might be subject to change. Currently, the three interfaces are:

- `Install`
- `Reinstall`
- `Uninstall`

### 1.2 Return Value

Each interface returns two values `(bool, error)`.

If no error occurred, the return value would be `(true, nil)`. Otherwise, the result would be `(false, error)`.

## 2 Package

TL;DR:
Each plugin should have a separate folder under the `cmd/` directory. Refer to [this example](https://github.com/merico-dev/stream/blob/main/cmd/githubactions/main.go).
Put most of a plugin's code under the `internal/pkg/` directory and only use `cmd/` as the main entrance to the plugin code.

Check out our [Standard Go Project Layout](https://github.com/merico-dev/stream/blob/main/docs/project_layout.md) document for detailed instruction on the project layout.
