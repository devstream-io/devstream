# Creating a Plugin

## 0 Thanks for Contributing!

First, please read our [CONTRIBUTING](https://github.com/devstream-io/devstream/blob/main/CONTRIBUTING.md) doc.

## 1 Scaffolding Automagically

Run `dtm develop create-plugin --name=YOUR-PLUGIN-NAME` , and dtm will automatically generate the following file.

> ### /cmd/plugin/YOUR-PLUGIN-NAME/main.go

This is the only main entrance to the plugin code.

You do not need to change this file. If you want, feel free to submit a pull request to change the [plugin template](https://github.com/devstream-io/devstream/blob/main/internal/pkg/develop/plugin/template/main.go) directly.

> ### /docs/plugins/YOUR-PLUGIN-NAME.md

This is the automatically generated plugin documentation.

Although dtm is automagic, but it can’t read your mind. I’m afraid that you will have to write your own doc.

> ### /internal/pkg/plugin/YOUR-PLUGIN-NAME/

Please fill in the main logic of the plugin here.

You can check out our [Standard Go Project Layout](project-layout) document for detailed instruction on the project layout.

## 2 Interfaces

### 2.1 Definition

Each plugin needs to satisfy all the interfaces defined in the [plugin engine](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L10).

At the moment, there are 4 interfaces, which might be subject to change. Currently, the three interfaces are:

- [`Create`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L12)
- [`Read`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L13)
- [`Update`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L14)
- [`Delete`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L16)

### 2.2 Return Value

`Create`, `Read`, and `Update` interfaces return two values `(map[string]interface{}, error)`; the first being the "state".

`Delete` interface returns two values `(bool, error)`. It returns `(true, nil)` if there is no error; otherwise `(false, error)` will be returned.

If no error occurred, the return value would be `(true, nil)`. Otherwise, the result would be `(false, error)`.

## 3 How does plugin work?

DevStream uses [go plugin](https://pkg.go.dev/plugin) to implement custom DevOps plugins.

When you execute a command which calls any of the interfaces(`Create`, `Read`, `Update`, `Delete`), devstream's pluginengine will call the [`plugin.Lookup("DevStreamPlugin")` function](https://github.com/devstream-io/devstream/blob/38307894bbc08f691b2c5015366d9e45cc87970c/internal/pkg/pluginengine/plugin_helper.go#L28) to load the plugin, get the variable `DevStreamPlugin` that implements the ` DevStreamPlugin` interface, and then you can call the corresponding plugin logic functions.  This is why it is not recommended to modify the `/cmd/plugin/YOUR-PLUGIN-NAME/main.go` file directly.

Note: The `main()` in `/cmd/plugin/YOUR-PLUGIN-NAME/main.go` file will not be executed, it is only used to avoid the goclangci-lint error.
