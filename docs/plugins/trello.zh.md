# trello 插件

该插件用于将 Github 中的 Issue 实时同步到 trello 中。

# 用例

```yaml
--8<-- "trello.yaml"
```

**Notes:**

- Devstream 会帮你管理 Trello 看板的描述，请不要修改它。
- 该插件需要配置 `board.token`, 如何获取可以查看该 [文档](https://trello.com/app-key)。
- `scm` 配置字段用于表示代码仓库的配置信息，具体配置可查看[SCM配置项](./scm-option.zh.md)。

## Outputs

该插件产生以下输出:

- `boardId`
- `todoListId`
- `doingListId`
- `doneListId`