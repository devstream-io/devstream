package plugin

var TrelloGithubDefaultConfig = `tools:
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
    branch: main`
