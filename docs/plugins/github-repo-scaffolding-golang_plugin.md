## 1 `github-repo-scaffolding-golang` Plugin

This plugin installs github-repo-scaffolding-golang.

_This plugin depends on the following environment variable:_

- GITHUB_TOKEN

Set it before using this plugin.

***Tips:***

*1. if uninstall, the repo on github will be completely removed*

*2. if reinstall，the repo on github will be completely removed and recreated*

## 2 Usage:

**Please note that the `owner` parameter is case-sensitive.**

```yaml
tools:
# name of the instance with github-repo-scaffolding-golang
- name: go-webapp-repo
  plugin:
    # kind of the plugin
    kind: github-repo-scaffolding-golang
    # version of the plugin
    version: 0.2.0
  # options for the plugin
  options:
    # the repo's owner. It should be case-sensitive here; strictly use your GitHub user name; please change the value below.
    owner: YOUR_GITHUB_USERNAME
    # the repo which you'd like to create; please change the value below.
    repo: YOUR_REPO_NAME
    # the branch of the repo you'd like to hold the code
    branch: main
    # the image repo you'd like to push the container image; please change the value below.
    image_repo: YOUR_DOCKERHUB_USERNAME/YOUR_DOCKERHUB_IMAGE_REPO_NAME
```

Currently, all the parameters in the example above are mandatory.
