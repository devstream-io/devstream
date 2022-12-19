*[dtm]: The commandline tool of DevStream (short for DevsTreaM).
*[Config]: A set of YAML files serving as DevStream's configuration.
*[配置]: 由一个或一组 YAML 文件组成，定义了 DevStream 需要的所有信息。
*[Plugin]: DevStream uses core-plugin architecture, where the core is a state machine and each plugin handles the CRUD and integration of a certain DevOps tool.
*[Plugins]: DevStream uses core-plugin architecture, where the core is a state machine and each plugin handles the CRUD and integration of a certain DevOps tool.
*[插件]: DevStream 使用 core-plugin 架构，core 用作状态机，插件负责管理 DevOps 工具的 CRUD。
*[Tool]: A type of DevStream config, corresponding to a DevStream plugin, which does the CRUD and integration of a certain DevOps tool. The concept of Tool is used mainly in the context of Config.
*[Tools]: A type of DevStream config, corresponding to a DevStream plugin, which does the CRUD and integration of a certain DevOps tool. The concept of Tool is used mainly in the context of Config.
*[工具]: DevStream 配置的一种类型。每个工具(Tool)对应了一个 DevStream 插件，它可以安装、配置或集成一些 DevOps 工具。工具这一概念主要用在配置中。
*[App]: A type of DevStream config, corresponding to a real-world application, for which different tools such as CI/CD can be easily configured.
*[Apps]: A type of DevStream config, corresponding to a real-world application, for which different tools such as CI/CD can be easily configured.
*[应用]: DevStream 配置的一种类型，对应现实中的应用程序，使用应用这种类型的配置可以简化例如 CI/CD 等工具的配置。
*[PipelineTemplate]: A CI/CD pipeline definition which can be refered to by App.
*[PipelineTemplates]: A CI/CD pipeline definition which can be refered to by App.
*[流水线模板]: CI/CD 流水线的定义，可被应用所引用。
*[Output]: Each Tool can have some Output (in the format of a key/value map), so that other tools can refer to it as their input.
*[输出]: 每个工具(Tool)可能会有些输出，这是一个 map，可以在配置其他工具时引用
*[State]: Records the current status of your DevOps platform defined and created by DevStream.
*[状态]: 记录了当前 DevOps 工具链的状态，包括每个工具的配置以及当下状态。
