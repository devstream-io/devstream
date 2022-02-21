## 1 Trello-Github Integration Plugin

This plugin creates a new Trello board and integrate it with your GitHub repo.

## 2 Usage:

_This plugin depends on the following two environment variables:_

- TRELLO_API_KEY
- TRELLO_TOKEN

Set the values accordingly before using this plugin.

To create a Trello API key and token, see [here](https://docs.servicenow.com/bundle/quebec-it-asset-management/page/product/software-asset-management2/task/generate-trello-apikey-token.html).

```yaml
tools:
- name: trello-github-integ-default
  # plugin profile
  plugin:
    # kind of this plugin
    kind: trello-github-integ
    # version of the plugin
    version: 0.0.1
  # options for the plugin
  # checkout the version from the GitHub releases
  options:
    # the repo's owner
    owner: ironcore864
    # the repo where you'd like to setup GitHub Actions
    repo: go-hello-http
    # integration tool name
    api:
      name: trello
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: master
```

Currently, all the parameters in the example above are mandatory.
