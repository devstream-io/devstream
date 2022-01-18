## 1 Trello-Github Integration Plugin

This plugin creates some GitHub Actions workflows.

Currently, only Golang is supported.

We now support trello, according to our roadmap.

## 2 Usage:

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
