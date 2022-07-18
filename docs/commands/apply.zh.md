# dtm apply

当 _applying_ 配置文件 `dtm` 时, 会发生下面这些情况:

## 1 每个 _Tool_ 在 _Config_ 中定义

我们比较了 _Tool_、它的状态和它之前创建的资源（如果状态存在）。

我们根据比较结果生成变更计划：

- 如果 _Tool_ 不存在 _State_， 将调用 `Create` 接口.
- 如果 _Tool_ 存在 _State_， 但是 _Config_ 与 _State_ 不同 (这意味着用户可能在上次 `apply` 应用之后更新了配置，) 将调用 `Update` 接口。
- 如果 _Tool_ 存在 _State_， 并且 _Config_ 与 _State_ 相同， 我们尝试读取 _Resource_。
    - 如果 _Resource_ 不存在， 将调用 `Create` 接口。这可能表明在最后一次成功的 `apply` 之后手动删除了 _Resource_。
    - 如果 _Resource_ 存在，但 _State_ 漂移 (意思是有人修改了它)， 将调用 `Update` 接口。
    - 最后但同样重要的是， 如果 _Resource_ 与 _State_ 完全相同，则不会发生任何事情。

## 2 每个  _State_ _Config_ 中没有 _Tool_

我们生成一个 `Delete` 更改来删除 _Resource_。因为配置中没有 _Tool_，但存在 _State_， 这意味着可能之前已经创建了 _Resource_，然后用户从 _Config_ 中删除了 _Tool_，这意味着用户不再需要 _Resource_。
