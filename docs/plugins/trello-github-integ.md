# trello-github-integ Plugin

This plugin creates a new GitHub Actions workflow(trello-github-integration) and uploads it to your GitHub repo.

## Usage

This plugin depends on and can be used together with the `trello` plugin (see document [here](./trello.md)).

`trello-github-integ` plugin can also use `trello` plugin's outputs as input. See the example below:

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
# name of the tool
- name: trello-github-integ
  # id of the tool instance
  instanceID: default
  # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ "trello.default" ]
  # options for the plugin
  options:
    # the repo's owner. It should be case-sensitive here; strictly use your GitHub user name; please change the value below.
    owner: YOUR_GITHUB_USERNAME
    # the repo where you'd like to setup GitHub Actions; please change the value below.
    # the repo's org. If you set this property, then the new repo will be created under the org you're given, and the "owner" setting above will be ignored.
    org: YOUR_ORGANIZATION_NAME
    repo: YOUR_REPO_NAME
    # reference parameters come from dependency, their usage will be explained later
    boardId: ${{ trello.default.outputs.boardId }}
    todoListId: ${{ trello.default.outputs.todoListId }}
    doingListId: ${{ trello.default.outputs.doingListId }}
    doneListId: ${{ trello.default.outputs.doneListId }}
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: main
```

Replace the following from the config above:

- `YOUR_GITHUB_USERNAME`
- `YOUR_REPO_NAME`
- `KANBAN_BOARD_NAME`

In the example above:

- We create a Trello board using `trello` plugin, and the board is marked to be used for repo YOUR_GITHUB_USERNAME/YOUR_REPO_NAME.
- `trello-github-integ` plugin depends on `trello` plugin, because we use `trello` plugin's outputs as the input for the `trello-github-integ` plugin.

Pay attention to the `${{ xxx }}` part in the example. `${{ TOOL_NAME.TOOL_INSTANCE_ID.outputs.var}}` is the syntax for using an output.
