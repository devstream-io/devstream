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
  - name: default
    # plugin profile
    plugin:
      # kind of this plugin
      kind: trello
      # version of the plugin
      version: 0.2.0
    # options for the plugin, checkout the version from the GitHub releases
    options:
      # the repo's owner (if kanbanBoardName is empty, use owner/repo as the boardname)
      owner: lfbdev
      # the repo name  (if kanbanBoardName is empty, use owner/repo as the boardname)
      repo: golang-demo
      # the Tello board name 
      kanbanBoardName: kanban-name
```

## 3. Use Together with the `trello-github-integ` Plugin

This plugin can be used together with the `trello-github-integ` plugin (see document [here](./trello-github-integ_plugin.md).)

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

- We put `default.trello` as a dependency by using the `dependsOn` keyword.
- We use `default.trello`'s output as input for the `default_trello_github` plugin.

Pay attention to the `${{ xxx }}` part in the example. `${{ TOOL_NAME.TOOL_KIND.outputs.var}}` is the syntax for using an output.
