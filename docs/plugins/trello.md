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
tools:
# name of the tool
- name: trello
  # id of the tool instance
  instanceID: default
  # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool
  dependsOn: []
  # options for the plugin
  options:
    # the repo's owner
    owner: YOUR_GITHUB_USERNAME
    # the repo's org. If you set this property, then the new repo will be created under the org you're given, and the "owner" setting above will be ignored.
    org: YOUR_ORGANIZATION_NAME
    # for which repo this board will be used
    repo: YOUR_REPO_NAME
    # the Tello board name. If empty, use owner/repo as the board's name.
    kanbanBoardName: KANBAN_BOARD_NAME
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
