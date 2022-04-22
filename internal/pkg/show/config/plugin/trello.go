package plugin

var TrelloDefaultConfig = `tools:
# name of the tool
- name: trello
  # id of the tool instance
  instanceID: my-trello-board
  # optional; if specified, dtm will make sure the dependency is applied first before handling this tool
  dependsOn: [ "TOOL1_NAME.TOOL1_PLUGIN", "TOOL2_NAME.TOOL2_PLUGIN" ]
  # options for the plugin
  options:
    # the repo's owner
    owner: YOUR_GITHUB_USERNAME
    # the repo's org. If you set this property, then the new repo will be created under the org you're given, and the "owner" setting above will be ignored.
    org: YOUR_ORGANIZATION_NAME
    # for which repo this board will be used
    repo: YOUR_REPO_NAME
    # the Tello board name. If empty, use owner/repo as the board's name.
    kanbanBoardName: KANBAN_BOARD_NAME`
