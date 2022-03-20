## 1 `trello-github-integ` Plugin

This plugin creates a new GitHub Actions workflow(trello-github-integration) and uploads it to your GitHub repo.

## 2 Usage:

This plugin depends on and can be used together with the `trello` plugin (see document [here](./trello_plugin.md)).

`trello-github-integ` plugin can also use `trello` plugin's outputs as input. See the example below:

```yaml
tools:
- name: my-trello-board
  plugin:
    kind: trello
    version: 0.3.0
  dependsOn: ["demo.github-repo-scaffolding-golang"]
  options:
    owner: YOUR_GITHUB_USERNAME
    repo: YOUR_REPO_NAME
    kanbanBoardName: KANBAN_BOARD_NAME
- name: trello-github
  # plugin profile
  plugin:
    # kind of this plugin
    kind: trello-github-integ
    # version of the plugin
    # checkout the version from the GitHub releases
    version: 0.3.0
  # optional; if specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ "my-trello-board.trello" ]
  # options for the plugin
  options:
    # the repo's owner. It should be case-sensitive here; strictly use your GitHub user name; please change the value below.
    owner: YOUR_GITHUB_USERNAME
    # the repo where you'd like to setup GitHub Actions; please change the value below.
    repo: YOUR_REPO_NAME
    # reference parameters come from dependency, their usage will be explained later
    boardId: ${{ my-trello-board.trello.outputs.boardId }}
    todoListId: ${{ my-trello-board.trello.outputs.todoListId }}
    doingListId: ${{ my-trello-board.trello.outputs.doingListId }}
    doneListId: ${{ my-trello-board.trello.outputs.doneListId }}
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

Pay attention to the `${{ xxx }}` part in the example. `${{ TOOL_NAME.TOOL_KIND.outputs.var}}` is the syntax for using an output.
