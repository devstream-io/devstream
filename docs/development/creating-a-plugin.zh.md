# 创建一个插件

## 0 感谢你的贡献!

首先, 请认真阅读[如何贡献](https://github.com/devstream-io/devstream/blob/main/CONTRIBUTING.md) 文档。

## 1 自动生成插件脚手架

运行 `dtm develop create-plugin --name=YOUR-PLUGIN-NAME`, 然后dtm就会自动生成以下文件.

> ### /cmd/plugin/YOUR-PLUGIN-NAME/main.go

这个文件是生成的插件代码的唯一入口

你不需要对它做任何更改. 当然如果你想的话, 请直接提交一个pull request 来直接更改[插件模板](https://github.com/devstream-io/devstream/blob/main/internal/pkg/develop/plugin/template/main.go)。

> ### /docs/plugins/YOUR-PLUGIN-NAME.md

这个是自动生成的插件文档

虽然 dtm 很智能，但是它无法看穿你的心思。所以你还是需要自己编辑这个文档。

> ### /internal/pkg/plugin/YOUR-PLUGIN-NAME/

请在这儿补充该插件的主要逻辑

具体详细说明你可以通过阅读[标准的go项目布局](project-layout.zh.md) 来了解具体项目布局的说明

## 2 Interfaces

### 2.1 Definition

每个新创建的插件都需要实现[DevStreamPlugin](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L10)中的所有接口。

目前有四个接口，它们分别是:

- [`Create`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L12)
- [`Read`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L13)
- [`Update`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L14)
- [`Delete`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L16)

### 2.2 接口返回值

`Create`, `Read`, 和 `Update` 接口都是两个返回值 `(map[string]interface{}, error)`; 第一个返回值就是 "状态".

`Delete` 接口也有两个返回值 `(bool, error)`. 如果返回 `(true, nil)` 则说明接口执行成功成功，否则返回 `(false, error)` .

## 3 插件是怎么工作的?

DevStream 通过[go plugin](https://pkg.go.dev/plugin)来实现自定义的devops插件

如果你执行的命令中调用到了(`Create`, `Read`, `Update`, `Delete`)这几个接口中的某个,比如调用 `Create`接口， DevStream的插件引擎[`plugin.Lookup("DevStreamPlugin")` 函数](https://github.com/devstream-io/devstream/blob/38307894bbc08f691b2c5015366d9e45cc87970c/internal/pkg/pluginengine/plugin_helper.go#L28)就会去加载该插件, 并获取实现 `DevStreamPlugin` 接口的变量 ` symDevStreamPlugin` 后, 就可以用它去执行我想要执行的`create`函数了.  这也是不推荐直接修改 `/cmd/plugin/YOUR-PLUGIN-NAME/main.go` 文件的原因.

注意: 文件 `/cmd/plugin/YOUR-PLUGIN-NAME/main.go` 的 `main` 方法不会被执行, 它仅仅是用来避免golangci-lint报错。

