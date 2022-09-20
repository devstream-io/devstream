#  dtm delete

## 1 普通（非强制）删除

使用 `dtm delete` 执行 _删除_ 时，会发生：

- 读取 _Config_
- 对于 _Config_ 中定义的每一个 _Tool_，如果有对应的 _State_，就会调用 `Delete` 接口。

_Note: the `Delete` change will be executed only if the _State_ exists._

## 2 强制删除

使用 `dtm delete --force` 执行 _强制删除_ 时，会发生：

- 读取 _Config_
- 对于 _Config_ 中定义的每一个 _Tool_，调用 `Delete` 接口。

_说明："普通删除"和"强制删除"的区别在于，无论 _State_ 是否存在，`dtm` 都会尝试调用 `Delete` 接口。此举目的在于，当 _State_ 损坏甚至丢失时（我们希望这只发生在开发环境），仍能有办法清除 _Tools_。

## 3 命令行参数

| 短  | 长            | 默认值                    | 描述        |
|-----|---------------|--------------------------|------------|
|     | --force       | `false`                  | 强制删除    |
| -f  | --config-file | `"config.yaml"`          | 配置文件路径 |
| -d  | --plugin-dir  | `"~/.devstream/plugins"` | 插件目录    |
| -y  | --yes         | `false`                  | 取消二次确认 |
