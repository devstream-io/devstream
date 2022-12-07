# 架构

本文介绍了DevStream的架构，总结了DevStream的主要组件，以及数据、命令是如何在各个组件之间流转的。

## 0 工作流程

下图展示了DevStream是如何执行一个用户命令的。

![DevStream架构图](../../images/architecture-overview.png)

DevStream主要由三大块组成：

- CLI：处理用户输入的命令和参数
- Plugin engine: 插件引擎，通过调用其他组件（config manager, plugin manager, state manager, backend manager）来实现DevStream的核心功能。
- 插件：实现某个DevOps工具的CRUD接口。

## 1 CLI

_注意：为了简单起见，CLI被命名为`dtm`（DevOps Toolchain Manager)。_

用户运行`dtm`时，会调用[`devstream`](https://github.com/devstream-io/devstream/tree/main/cmd/devstream)包中的一个命令。所有命令的源文件定义都在这个文件夹中。

然后，每个命令调用[`internal/pkg`](https://github.com/devstream-io/devstream/tree/main/internal/pkg/pluginengine)下的`pluginengine`包。

`pluginengine`首先调用`configmanager`，将本地YAML配置文件读取到一个结构体中，然后调用`pluginmanager`来下载所需的插件。

之后，`pluginengine`调用`statemanager`来计算config、状态和实际DevOps工具的状态之间的"差异"。最后，`pluginengine`根据这变更执行对应的操作，并更新状态。在执行过程中，`pluginengine`加载每个插件（`*.so`文件）并根据每个变更调用相应的接口。

## 2 插件引擎

`pluginengine`有以下各种职责：

- 根据配置文件来决定需要哪些插件，以及确保这些插件的存在
- 根据配置、状态和工具的实际状态来生成变更
- 通过加载每个插件的so文件并调用所需的接口来执行变更

Plugin engine通过调用以下模块来实现目标：

### 2.1 配置管理器

[`configmanager`包](https://github.com/devstream-io/devstream/blob/main/internal/pkg/configmanager/configmanager.go#L23)中的模型结构体用来表示顶级配置的结构。

### 2.2 插件管理器

[`pluginmanager`包](https://github.com/devstream-io/devstream/blob/main/internal/pkg/pluginmanager/manager.go)负责根据配置下载必要的插件。

如果本地已存在所需版本的插件，则不会再次下载。

### 2.3 状态管理器

[`statemanager`包](https://github.com/devstream-io/devstream/blob/main/internal/pkg/statemanager/manager.go)是用来管理“状态”的，即哪些变更已经成功完成，哪些失败。

`statemanager`将状态存储在[`backend`](https://github.com/devstream-io/devstream/blob/main/internal/pkg/backend/backend.go)中。

### 2.4 后端管理器

[`backend`包](https://github.com/devstream-io/devstream/tree/main/internal/pkg/backend)是“后端”管理器，在这里，“后端”指的是状态的实际存储。目前，DevStream已经支持本地和远程（与AWS S3兼容）状态。

## 3 插件

一个 _plugin_ 实现了上述的预定义接口。

它执行的包括"创建"、"读取"、"更新"和"删除"等操作。

要开发一个新的插件，请参阅[创建一个插件](dev/creating-a-plugin.zh.md)。
