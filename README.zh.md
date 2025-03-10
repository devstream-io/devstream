# Chatflow 智能工作流引擎

> 注意：DevChat IDE 插件尚未迁移至本实现，在当前插件中创建智能工作流请参见[文档说明](https://docs.devchat.ai/zh/quick-start/create_workflow)。

本实现遵循[《Chatflow 接口规范》版本v1.0.0](specs/interface-v1.0.0.zh.md)。未实现 [`shell-python` 运行时系统](specs/runtime-v0.1.0.zh.md)的 Python 环境和包自动管理能力。

## 示例

本项目 `/demo` 目录是一个包含了工作流 `/hello_world` 实现的示例，下面操作步骤可实现对该工作流的调用。以运行于 Linux 或 macOS 操作系统为例。

- 设置运行时依赖的环境变量，在终端执行如下命令：
  - `export $CHATFLOW_BASE=/path/to/demo/`，其中 `/path/to/demo/` 应替换为本地 `demo` 目录的完整路径。
