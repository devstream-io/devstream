# 架构

本文件总结了DevStream的主要组成部分以及数据在这些组成部分之间的流动。

## 0 数据流

下图显示了DevStream如何执行一个用户命令的近似情况。

![DevStream架构图](../images/architecture-overview.png)

有三个主要部分:

- CLI：处理用户输入、命令和参数
- `pluginengine`：插件引擎，通过调用其他模块（`configloader`、`pluginmanager`、`statemanager`等）实现核心功能。
- 插件：为某个DevOps工具实现实际的CRUD接口。

## 1 CLI ("devstream "软件包)

注意：为了简单起见，CLI被命名为`dtm`（DevOps工具管理器），代替其全名DevStream。

每次用户运行`dtm`程序时，执行立即转移到[`devstream`](https://github.com/devstream-io/devstream/tree/main/cmd/devstream)包中的一个 "命令 "实现，所有命令的定义都在这个文件夹中。

然后，每个命令调用[`internal/pkg`](https://github.com/devstream-io/devstream/tree/main/internal/pkg/pluginengine)下的`pluginengine`包。

`pluginengine`首先调用`configloader`本地YAML配置文件读取到一个结构中。

然后调用`pluginmanager`来下载所需的插件。

之后，`pluginengine`调用状态管理器来计算congfig、状态和实际DevOps工具的状态之间的 "差异"。最后，`pluginengine`根据这变更执行对应的操作，并更新状态。在执行过程中，`pluginengine`加载每个插件（`*.so`文件）并根据每个变更调用预定义接口。

## 2 插件引擎

`pluginengine`有几个职责：

- 确保所需的插件（根据配置文件的设置）是存在的
- 根据配置、状态和工具的实际状态生成变更
- 通过加载每个插件和调用所需的动作来执行这些。

它通过调用以下模块来实现目标。

### 2.1 配置加载器

包[`configloader`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/configloader/config.go#L19)中的模型类型代表了顶层的配置结构。

### 2.2 插件管理器

[`pluginmanager`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginmanager/manager.go)负责根据配置下载必要的插件。

如果本地已经存在所需版本的插件，将不再下载。

### 2.3 状态管理器

[`statemanager`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/statemanager/manager.go)负责管理 "状态"，即哪些事情已经成功完成，哪些没有。

`statemanager`将状态存储在一个[`backend`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/backend/backend.go)中。

## 3 插件

一个_plugin_实现了上述的预定义接口。

它执行的包括 "创建"、"读取"、"更新 "和 "删除 "等操作。

要开发一个新的插件，请参阅[创建一个插件](./development/creating-a-plugin.md)。