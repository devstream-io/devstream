# dtm apply

当使用 `dtm` _apply_ 配置文件时，会发生以下事情：

## 1 对于 _Config（配置文件）_ 中定义的每个 _Tool_

我们将会对比 _Tool_、它的 _State_ 和它之前创建的 _Resoruce_（如果存在 state）。

根据对比结果，变更规则如下：

- 若该 _Tool_ 不在 _State_ 中，调用 `Create` 接口；
- 若该 _Tool_ 存在于 _State_ 中，但当前 _Config_ 中关于该 _Tool_ 的配置与 _State_ 中的定义不同（意味着用户可能在上一次 `apply` 之后更新了配置），则调用 `Update` 接口；
- 若该 _Tool_ 存在于 _State_ 中，且当前 _Config_ 中关于该 _Tool_ 的配置与 _State_ 相同。我们将会继续尝试通过 `Read` 接口读取 _Resource_ ，并与 _State_ 中记录的 _Resource_ 比对：
  - 若从 `Read` 读取到的 _Resource_ 不存在，调用 `Create` 接口。这可能表明 _Resource_ 在最后一次成功 `apply` 后被手动删除；
  - 若从 `Read` 读取到的 _Resource_ 存在，但与 _State_ 中记录的 _Resource_ 不一致（意味着有人修改了 _State_ 或插件状态发生了变化），调用 `Update` 接口；
  - 最后，若读取到的 _Resource_ 和 _State_ 中的 _Resource_ 一致，什么都不做。

## 2 _State_ 中含有某 _Tool_，但 _Config_ 中没有

我们将对其执行"删除"操作，以删除相应的 _Resource_ 。因为 _State_ 中含有此 _Tool_，但配置文件中不存在了，这意味着用户先前为该 _Tool_，创建了 _Resource_，但后面从 _Config_ 中删除了该 _Tool_，表明用户不想要该 _Resource_ 了。

## 3 命令行参数

| 短  | 长            | 默认值                    | 描述        |
|-----|---------------|--------------------------|------------|
| -f  | --config-file | `"config.yaml"`          | 配置文件路径 |
| -d  | --plugin-dir  | `"~/.devstream/plugins"` | 插件目录    |
| -y  | --yes         | `false`                  | 取消二次确认 |
