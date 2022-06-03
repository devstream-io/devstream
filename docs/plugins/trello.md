# trello Plugin

This plugin creates a new Trello board and lists.

## Usage

_This plugin depends on the following two environment variables:_

- TRELLO_API_KEY
- TRELLO_TOKEN

Set the values accordingly before using this plugin.

## Tips:
_Trello board description is managed by DevStream, please don't modify it._

To create a Trello API key and token, see [here](https://trello.com/app-key).

```yaml
--8<-- "trello.yaml"
```

Replace the following from the config above:

- `YOUR_GITHUB_USERNAME`
- `YOUR_REPO_NAME`
- `KANBAN_BOARD_NAME`

## Outputs

This plugin has four outputs:

- `boardId`
- `todoListId`
- `doingListId`
- `doneListId`

which can be used by the `trello-github-integ` plugin. Refer to the [`trello-github-integ` plugin doc](./trello-github-integ.md) for more details.
