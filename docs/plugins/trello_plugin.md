## 1 `trello` Plugin

This plugin creates a new Trello board and lists.

## 2 Usage:

_This plugin depends on the following two environment variables:_

- TRELLO_API_KEY
- TRELLO_TOKEN

Set the values accordingly before using this plugin.

## 3 Tips:
_Trello board description is managed by DevStream, please don't modify it._

To create a Trello API key and token, see [here](https://trello.com/app-key).

```yaml
tools:
- name: my-trello-board
  # plugin profile
  plugin:
    # kind of this plugin
    kind: trello
    # version of the plugin
    version: 0.2.0
  # options for the plugin, checkout the version from the GitHub releases
  options:
    # the repo's owner
    owner: YOUR_GITHUB_USERNAME
    # for which repo this board will be used
    repo: YOUR_REPO_NAME
    # the Tello board name. If empty, use owner/repo as the board's name.
    kanbanBoardName: KANBAN_BOARD_NAME
```

Replace the following from the config above:

- `YOUR_GITHUB_USERNAME`
- `YOUR_REPO_NAME`
- `KANBAN_BOARD_NAME`

## 3. Outputs

This plugin has four outputs:

- `boardId`
- `todoListId`
- `doingListId`
- `doneListId`

which can be used by the `trello-github-integ` plugin. Refer to the [`trello-github-integ` plugin doc](./trello-github-integ_plugin.md) for more details.
