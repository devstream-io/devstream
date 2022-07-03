# 创建一个插件

## 0 感谢你的贡献!

首先，请阅读我们的 [CONTRIBUTING](https://github.com/devstream-io/devstream/blob/main/CONTRIBUTING.md) 文档。

## 1 自动建立新插件的代码框架

运行`dtm develop create-plugin --name=YOUR-PLUGIN-NAME` ，dtm将自动生成以下文件。

> ### /cmd/plugin/YOUR-PLUGIN-NAME/main.go

这是该插件代码的唯一主要入口。

你不需要修改这个文件。如果你觉得自动生成的这个文件有问题，你可以创建一个PR来直接修改[模板](https://github.com/devstream-io/devstream/blob/main/internal/pkg/develop/plugin/template/main.go)。

> ### /docs/plugins/YOUR-PLUGIN-NAME.md

这是自动生成的插件的文档。

虽然`dtm`的目的就是自动化，但它并不能魔法般的生成文档。你需要自己编写自己想要创建的这个插件的文档。

> ###/internal/pkg/plugin/YOUR-PLUGIN-NAME/

请在这里编写插件的主要逻辑。

可以查看我们的[Standard Go Project Layout](project-layout.md)文件，了解关于项目布局的详细说明。

## 2 接口

### 2.1 定义

每个插件都需要实现[pluginengine](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L10)中定义的所有接口。

目前，有4个接口，可能会有变化。目前，这4个接口是。

- [`create`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L12)
- [`read`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L13)
- [`update`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L14)
- [`delete`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L16)

### 2.2 返回值

`create`、`read`和`update`方法返回两个值`(map[string]interface{}, error)`；第一个是 "状态"。

`delete'接口返回两个值`(bool, error)`。如果没有错误，它返回`(true, nil)`；否则将返回`(false, error)`。

如果没有发生错误，返回值将是`(true, nil)`。否则，结果将是`(false, error)`。

## 3 插件是如何工作的？

DevStream是使用[go plugin](https://pkg.go.dev/plugin)来实现自定义插件的。。

当你执行一个调用任何接口(`Create`, `Read`, `Update`, `Delete`)的命令时，DevStream的`pluginengine`会调用[`plugin.Lookup("DevStreamPlugin")`函数](https://github.com/devstream-io/devstream/blob/38307894bbc08f691b2c5015366d9e45cc87970c/internal/pkg/pluginengine/plugin_helper.go#L28)来加载插件，获得实现`DevStreamPlugin`接口的变量`DevStreamPlugin`，然后你就可以调用相应的插件接口。所以我们不建议你直接修改`/cmd/plugin/YOUR-PLUGIN-NAME/main.go`文件，因为该文件是根据接口定义自动生成好的。

注意：`/cmd/plugin/YOUR-PLUGIN-NAME/main.go`文件中的`main()`不会被执行，它只是用来避免golangci-lint错误。
