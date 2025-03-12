# Devstream 智能工作流引擎

> 注意：DevChat IDE 插件尚未迁移至本实现，在当前插件中创建智能工作流请参见[文档说明](https://docs.devchat.ai/zh/quick-start/create_workflow)。

本实现遵循[Devstream 接口规范》版本v1.0.0](specs/interface-v1.0.0.zh.md)。未实现 [`shell-python` 运行时系统](specs/runtime-v0.1.0.zh.md)的 Python 环境和包自动管理能力。

本工作流引擎采用 Python 编程语言实现，可作为 Python 包被用户脚本或其他工作流引入和使用。引擎默认使用调用时所处 Python 环境的解释器。

## 示例

本项目 `/demo` 目录是一个包含了工作流 `/hello_world` 实现的示例，下面操作步骤可实现对该工作流的调用。命令行以 Linux 或 macOS 操作系统终端为例。

1. 设置运行时依赖的环境变量，在终端执行如下命令：
   - `export DEVSTREAM_BASE=/path/to/demo/`，其中 `/path/to/demo/` 替换为本地 `demo` 目录的完整路径。
2. 本地应安装好 Python 环境，并在本地 `demo` 目录下的 `settings.yml` 文件中注册，该文件形如：
   
   ```yaml
   environments:
    myenv: "/path/to/python/bin"
   ```
   
   其中 `/path/to/python/bin` 替换为本地 Python 环境解释器的完整路径。
3. 我们假设一个用户脚本或另外的工作流希望调用 `/hello_world`，那么 [`/demo/caller.py`](demo/caller.py) 是一个实现这种调用的实例：
   ```python
   import devstream
   
   devstream.call("/hello_world", "This is some user input.")
   ```
4. 在终端中执行如下命令行，可实测调用 `/hello_world` 的结果：
   - `python /path/to/caller.py`，其中 `python` 指向终端当前所处的 Python 环境，也是运行 Devstream 引擎的环境；`/path/to/caller.py` 替换为脚本实际所处的路径。
