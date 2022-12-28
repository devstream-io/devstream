# dtm verify

`dtm verify` 命令将检查以下内容：

## 1 配置文件

`dtm verify` 首先验证配置文件是否能成功加载。

若失败，可能会打印出以下信息：

- 若配置文件不存在，它会提醒你是否忘记使用 "-f" 参数来指定配置文件。
- 若配置文件格式不正确，它将打印一些错误。

## 2 插件列表

接着，`dtm verify` 会检查配置文件中引用的所有插件是否存在。

如果不存在，它会提示你忘了先运行 `dtm init`（也有可能是 "plugin-dir" 参数不对）。

## 3 State

`dtm verify` 也会尝试创建一个操作 _backend_ 的 _State_ 管理器。若 _State_ 出现问题，将提示错误所在。

## 4 Config / State / ResourceStatus

关于 _Config_ 、 _State_ 和 _ResourceStatus_ 的定义，见[核心概念](../core-concepts/overview.zh.md)。

`dtm verify` 尝试判断 _Config_ 是否与 _State_ 和 _ResourceStatus_ 匹配。如果不匹配，它会告诉你到底是什么不一样，以及如果你运行 `dtm apply` 会发生什么。

如果所有上述检查都成功，`dtm verify` 最后将输出 "Verify succeeded." 日志提示。

## 5 命令行参数

| 短  | 长            | 默认值                    | 描述        |
|-----|---------------|--------------------------|------------|
| -f  | --config-file | `"config.yaml"`          | 配置文件路径 |
| -d  | --plugin-dir  | `"~/.devstream/plugins"` | 插件目录    |
