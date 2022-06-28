# 创建一个插件

## 0 谢谢你的贡献!

首先，请阅读我们的 [CONTRIBUTING](https://github.com/devstream-io/devstream/blob/main/CONTRIBUTING.md) doc。

## 1 自动建立脚手架

运行`dtm develop create-plugin --name=YOUR-PLUGIN-NAME` ，dtm将自动生成以下文件。

> ### /cmd/plugin/YOUR-PLUGIN-NAME/main.go

这是该插件代码的唯一主要入口。

你不需要改变这个文件。如果你愿意，可以随时提交一个拉动请求，直接修改[插件模板](https://github.com/devstream-io/devstream/blob/main/internal/pkg/develop/plugin/template/main.go)。

> ### /docs/plugins/YOUR-PLUGIN-NAME.md

这是自动生成的插件文档。

虽然dtm是自动的，但它不能完全如你所想。所以恐怕还得你自己写文档了。

> ###/internal/pkg/plugin/YOUR-PLUGIN-NAME/

请在这里填写该插件的主要逻辑。

你可以查看我们的[Standard Go Project Layout](project-layout.md)文件，了解关于项目布局的详细说明。

## 2 接口

### 2.1 定义

每个插件都需要实现[插件引擎](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L10)中接口定义的所有方法。

目前，有4个接口，可能会有变化。目前，这4个接口是。

- [`创建`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L12)
- [`读取`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L13)
- [`更新`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L14)
- [`删除`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginengine/plugin.go#L16)

### 2.2 返回值

`创建`、`读取`和`更新`方法返回两个值`(map[string]interface{}, error)`；第一个是 "状态"。

`删除'接口返回两个值`(bool, error)`。如果没有错误，它返回`(true, nil)`；否则将返回`(false, error)`。

如果没有发生错误，返回值将是`(true, nil)`。否则，结果将是`(false, error)`。

## 3 插件是如何工作的？

DevStream使用[go plugin](https://pkg.go.dev/plugin)来实现自定义DevOps插件。

当你执行一个调用任何接口(`Create`, `Read`, `Update`, `Delete`)的命令时，devstream的pluginengine会调用[`plugin.Lookup("DevStreamPlugin")`函数](https://github.com/devstream-io/devstream/blob/38307894bbc08f691b2c5015366d9e45cc87970c/internal/pkg/pluginengine/plugin_helper.go#L28)来加载插件，获得实现`DevStreamPlugin`接口的变量`DevStreamPlugin`，然后你就可以调用相应的插件逻辑函数。 这就是为什么不建议直接修改`/cmd/plugin/YOUR-PLUGIN-NAME/main.go`文件。

注意：`/cmd/plugin/YOUR-PLUGIN-NAME/main.go`文件中的`main()`不会被执行，它只是用来避免golangci-lint错误。
