# 概览

您需要熟知以下概念：Git 核心概念，Docker, Kubernetes, 持续集成，持续交付以及 GitOps概念. 这些都是 DevStream 的核心概念。

## DevStream的构架

以下构架图展示了 DevStream 大致的流程：
![](../images/architecture-overview.png)

## 工作流程

![config state resource-status workflow](../images/config_state_resource.png)

## 配置，工具，状态和资源状态

构架文档阐述了 DevStream 的基本工作原理。请确保你在阅读文档其他部分之前完成此部分的阅读。

### 1. 配置(Config)

DevStream 在配置文件中定义了 DevOps 工具链。

有三种配置文件：

- 主配置文件(core config)
- 变量配置文件(variable config)
- 工具配置文件(tool config)

主配置文件包含了以下内容：

- `varFile`: 变量配置文件的路径
- `toolFile`: 工具文件的路径
- `pluginDir`: 插件目录的路径，默认为 `~/.devstream/plugins`, 或使用 `-d` 参数指定一个目录
- `state`: 与 State 关联的设置。更多信息，请看[这里](./state.zh.md)

变量配置文件是一个包含了键值对的YAML文件，可以用于工具配置文件。

工具配置文件是一个名为 _Tools_ 的列表，其中每一个 _Tool_ 都包含了名称，实例ID（唯一标识）以及工具的选项(options)。

_注意: 你可以将多个YAML文件合并为同一个并用三个短横线(`---`)区分。更多信息见 [这里](https://stackoverflow.com/questions/50788277/why-3-dashes-hyphen-in-yaml-file) 和 [这里](https://www.javatpoint.com/yaml-structure)_。

### 2. 工具(Tool)

- 每个 _Tool_ 对应一个插件, 即可用于安装和配置，也可用于整合 DevOps 的工具。
- 每个 _Tool_ 有名称, (实例ID)InstanceID 和选项(Options), 定义在[这里](https://github.com/devstream-io/devstream/blob/main/internal/pkg/configmanager/toolconfig.go#L13)。
- 每个 _Tool_ 可以使用`dependsOn` 字段指定其依赖项。

 `dependsOn` 是一个字符串数组, 其中每一个元素都是一个依赖。 每个依赖项都以 "TOOL_NAME.INSTANCE_ID" 为格式命名。
例子见[这里](https://github.com/devstream-io/devstream/blob/main/examples/quickstart.yaml#L22)。

### 3. 状态(State)

_State_ 记录了当前 DevOps 工具链的状态，包括了每个工具的配置以及当下状态。

- _State_ 实际上是一个记录了状态的 Map, 定义在[这里](https://github.com/devstream-io/devstream/blob/main/internal/pkg/statemanager/state.go#L24)。
- Map 中的每一个状态都是一个包含了名称、插件、选项和资源的结构体，定义在[这里](https://github.com/devstream-io/devstream/blob/main/internal/pkg/statemanager/state.go#L16)。

### 4. 资源状态(ResourceStatus)

- 我们称插件创建了 _资源(Resource)_，而插件的 `Read()` 接口返回了此资源的描述，该描述也作为资源的状态（State）的一部分保存。

配置-状态-资源状态 工作流：

![config state resource-status workflow](../images/config_state_resource.png)
