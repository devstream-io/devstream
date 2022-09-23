#  dtm destroy

`dtm destroy` 的效果类似于 `dtm apply -f an_empty_config.yaml`.

`destroy` 命令的目的在于，一旦你在测试的时候不小心删除了配置文件，仍能通过此方式删除 _State_ 中定义的所有内容，以拥有干净的环境。

## 1 命令行参数

| 短  | 长            | 默认值                    | 描述        |
|-----|---------------|--------------------------|------------|
|     | --force       | `false`                  | 强制删除    |
| -d  | --plugin-dir  | `"~/.devstream/plugins"` | 插件目录    |
| -f  | --config-file | `"config.yaml"`          | 配置文件路径 |
| -y  | --yes         | `false`                  | 取消二次确认 |
