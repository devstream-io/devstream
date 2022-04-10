package plugin

var TrelloGithubDefaultConfig = `tools:
- name: my-trello-board
  # name of the plugin
  plugin: trello
  dependsOn: ["demo.github-repo-scaffolding-golang"]
  options:
    owner: YOUR_GITHUB_USERNAME
    repo: YOUR_REPO_NAME
    kanbanBoardName: KANBAN_BOARD_NAME
- name: trello-github
  # name of the plugin
  plugin: trello-github-integ
  # optional; if specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ "my-trello-board.trello" ]
  # options for the plugin
  options:
    # the repo's owner. It should be case-sensitive here; strictly use your GitHub user name; please change the value below.
    owner: YOUR_GITHUB_USERNAME
    # the repo where you'd like to setup GitHub Actions; please change the value below.
    # the repo's org. If you set this property, then the new repo will be created under the org you're given, and the "owner" setting above will be ignored.
    org: YOUR_ORGANIZATION_NAME
    repo: YOUR_REPO_NAME
    # reference parameters come from dependency, their usage will be explained later
    boardId: ${{ my-trello-board.trello.outputs.boardId }}
    todoListId: ${{ my-trello-board.trello.outputs.todoListId }}
    doingListId: ${{ my-trello-board.trello.outputs.doingListId }}
    doneListId: ${{ my-trello-board.trello.outputs.doneListId }}
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: main`
