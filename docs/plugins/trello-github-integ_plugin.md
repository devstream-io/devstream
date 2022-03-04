## 1 `trello-github-integ` Plugin

This plugin creates a new Trello board and integrate it with your GitHub repo.

## 2 Usage:

_This plugin depends on the following two environment variables:_

- TRELLO_API_KEY
- TRELLO_TOKEN

Set the values accordingly before using this plugin.

## 3 Tips:
_Trello board description is managed by DevStream, please don't modify it._

To create a Trello API key and token, see [here](https://docs.servicenow.com/bundle/quebec-it-asset-management/page/product/software-asset-management2/task/generate-trello-apikey-token.html).

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
  dependsOn: TOOL1_NAME.TOOL1_KIND,TOOL2_NAME.TOOL2_KIND,...
  # options for the plugin
  options:
    # the repo's owner. It should be case-sensitive here; strictly use your GitHub user name; please change the value below.
    owner: YOUR_GITHUB_USERNAME
    # the repo where you'd like to setup GitHub Actions; please change the value below.
    repo: YOUR_REPO_NAME
    # integration tool name
    api:
      name: trello
      # name of the Trello kanban board
      kanbanBoardName: kanban-name
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: main
```

Currently, all the parameters in the example above are mandatory.
