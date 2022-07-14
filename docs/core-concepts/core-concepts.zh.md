# 工具、配置和资源

# 主要概念概览

您需要熟知以下概念：核心Git，Docker, Kubernetes, 持续集成，持续交付以及GitOps. 这些都是DevStream的核心概念。

## DevStream的构架

以下构架图展示了DevStream大致的流程：
![](./images/architecture-overview.png)

## Workflow 流程

![config state resource workflow](../images/config_state_resource.png)

## 配置，工具，状态和资源

构架文档阐述了DevStream基本工作原理。请确保您在阅读文档其他部分之前首先阅读这一部分。

### 1. 配置

DevStream在配置文件中定义了开发运维工具链。

有三种配置文件：

- 主配置文件
- var配置文件
- 工具配置文件

主配置文件包含了以下内容：

- `varFile`: var文件的文件路径
- `toolFile`: 工具文件的文件路径
- `state`: 与状态关联的设置。更多信息，请看[这里](./stateconfig.md)

变量配置文件是一个包含了键值对的YAML文件，可以用于工具配置文件。

工具配置文件是一个 _Tools_，每一个都包含了名称，实例ID（唯一标识）以及工具的选项。

_注意: 你可以将多个YAML文件合并为同一个并用三个短破折号(`---`)区分。更多信息 [这里](https://stackoverflow.com/questions/50788277/why-3-dashes-hyphen-in-yaml-file) 和 [这里](https://www.javatpoint.com/yaml-structure)._

### 2. 工具

- 每个 _Tool_ 对应一个插件, 即可用于安装和配置，也可用于整合开发运维的工具。
- 每个 _Tool_ 有名称, InstanceID和选项, 定义在此 [这里](https://github.com/devstream-io/devstream/blob/main/internal/pkg/configloader/toolconfig.go#L13).
- 每个 _Tool_ 可以有自己的依赖项，具体根据`dependsOn` 字段指定.

 `dependsOn` 是一个字符串数组, 其中每一个元素都是一个依赖. 每个依赖项都以 "TOOL_NAME.INSTANCE_ID"为格式命名.
例子见 [这里](https://github.com/devstream-io/devstream/blob/main/examples/quickstart.yaml#L22) .

### 3. 状态

_state_ 记录了当下开发运维工具链的状况，包括了每个工具的配置以及当下状态。
- 这个 _State_ 实际上是状态的Map, 定义在 [这里](https://github.com/devstream-io/devstream/blob/main/internal/pkg/statemanager/state.go#L24).
- Map中的每一个状态都是一个包含了名称，插件，选项和资源的结构，定义在 [这里](https://github.com/devstream-io/devstream/blob/main/internal/pkg/statemanager/state.go#L16).

### 4. 资源

- 我们称创建的插件为 _Resource_，而插件的`Read()`界面返回此资源的描述，该描述也作为状态的一部分保存。

整个工作流为: 配置-状态-资源 工作流：

![config state resource workflow](../images/config_state_resource.png)
