## 1 `trello-github-integ` Plugin

This plugin creates a new GitHub Actions file(trello-github-integration) and upload to your GitHub repo.

## 2 Usage:

_This plugin depends on the plugin `trello`:_

```yaml
tools:
- name: trello-github-integ-default
  # plugin profile
  plugin:
    # kind of this plugin
    kind: trello-github-integ
    # version of the plugin
    # checkout the version from the GitHub releases
    version: 0.2.0
  # optional; if specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ "TOOL1_NAME.TOOL1_KIND", "TOOL2_NAME.TOOL2_KIND" ]
  # options for the plugin
  options:
    # the repo's owner. It should be case-sensitive here; strictly use your GitHub user name; please change the value below.
    owner: YOUR_GITHUB_USERNAME
    # the repo where you'd like to setup GitHub Actions; please change the value below.
    repo: YOUR_REPO_NAME
    # reference parameters come from dependency, their usage will be explained later
    boardId: ${{ default.trello.outputs.boardId }}
    todoListId: ${{ default.trello.outputs.todoListId }}
    doingListId: ${{ default.trello.outputs.doingListId }}
    doneListId: ${{ default.trello.outputs.doneListId }}
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: main
```

## 3. Use Together with the `trello` Plugin

This plugin can be used together with the `trello` plugin (see document [here](./trello_plugin.md).)

See the example below:

```yaml
---
tools:
  - name: trello_init_demo
    plugin:
      kind: trello
      version: 0.2.0
    options:
      owner: lfbdev
      repo: golang-demo
      kanbanBoardName: kanban-name
  - name: trello_github_integ_demo
    plugin:
      kind: trello-github-integ
      version: 0.2.0
    dependsOn: ["trello_init_demo.trello"]
    options:
      owner: lfbdev
      repo: golang-demo
      boardId: ${{ trello_init_demo.trello.outputs.boardId }}
      todoListId: ${{ trello_init_demo.trello.outputs.todoListId }}
      doingListId: ${{ trello_init_demo.trello.outputs.doingListId }}
      doneListId: ${{ trello_init_demo.trello.outputs.doneListId }}
      branch: main
```

In the example above:

- We put `default.trello` as dependency by using the `dependsOn` keyword.
- We use `default.trello`'s output as input for the `default_trello_github` plugin.

Pay attention to the `${{ xxx }}` part in the example. `${{ TOOL_NAME.TOOL_KIND.outputs.var}}` is the syntax for using an output.

