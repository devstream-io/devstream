# 配置(Config)

DevStream使用 YAML 文件来声明 DevOps 工具链的配置。

## 配置内容

正如概述中所提到的，配置包含了以下几个部分：

- `config`
- `vars`
- `tools`
- `apps`
- `pipelineTemplates`

其中，`config` 必选，`tools` 和 `apps` 至少有一个不为空，其余部分可选。

## 组织方式

DevStream 支持两种组织配置的方式：

- 单文件：你可以把这些部分写到一个 YAML 文件中
- 目录：也可以把它们分开放到同一个文件夹下的多个 YAML 文件中，只要这些文件的名字以 `.yaml` 或 `.yml` 结尾，且内容拼接后包含了 DevStream 所需的各个部分即可。

然后在执行 `init` 等命令时，加上 `-f` 或 `--config-file` 参数指定配置文件/目录的路径。

如：

- 单文件：`dtm init -f config.yaml`
- 目录：`dtm init -f dirname`

## 主配置

指 DevStream 本身的配置，即 `config` 部分，比如状态存储的方式等。详见 [这里](./state.zh.md)

## 变量语法

DevStream 提供了变量语法。使用 `vars` 用来定义变量的值，而后可以在 `tools`、`apps`、`pipelineTemplates` 中使用，语法是 `[[ varName ]]`。

示例：

```yaml
vars:
  githubUsername: daniel-hutao # 定义变量
  repoName: dtm-test-go
  defaultBranch: main

tools:
- name: repo-scaffolding
  instanceID: default
  options:
    destinationRepo:
      owner: [[ githubUsername ]] # 使用变量
      name: [[ repoName ]]
      branch: [[ defaultBranch ]]
      scmType: github
    # <后面略...>
```

## 工具的配置

`tools` 部分声明了工具链中的工具，详见 [这里](./tools.zh.md)

## 应用与流水线模板的配置

`apps` 部分声明了 应用 的配置，`pipelineTemplates` 声明了 流水线模板，详见 [这里](./apps.zh.md)
