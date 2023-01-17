# trello Plugin

This plugin integrates Trello with your GitHub repo.

## Usage

The following content is an example of the "tool file".

For more information on the main config, the tool file, and the var file of DevStream, see [Core Concepts Overview](../core-concepts/overview.md) and [DevStream Configuration](../core-concepts/config.md).

```yaml
--8<-- "trello.yaml"
```

**Notes:**

- Trello board description is managed by DevStream, please don't modify it.
- To config a `board.token`, see [here](https://trello.com/app-key).
- `scm` config option represents codebase location; for more info, you can refer to [SCM Config](./scm-option.md).

## Outputs

This plugin has four outputs:

- `boardId`
- `todoListId`
- `doingListId`
- `doneListId`
