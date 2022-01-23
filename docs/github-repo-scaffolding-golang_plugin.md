## 1 github-repo-scaffolding-golang Plugin

This plugin installs github-repo-scaffolding-golang.

## 2 Usage:

```yaml
tools:
# name of the instance with github-repo-scaffolding-golang
- name: github-repo-scaffolding-golang-demo
  plugin:
    # kind of the plugin
    kind: github-repo-scaffolding-golang
    # version of the plugin
    version: 0.0.1
  # options for the plugin
  options:
    # the repo's owner
    owner: daniel-hutao
    # the repo which you'd like to create
    repo: golang-demo
    # the branch of the repo you'd like to hold the code
    branch: main
    # the image repo you'd like to push the container image
    image_repo: daniel-hutao/golang-demo
```

Currently, all the parameters in the example above are mandatory.
