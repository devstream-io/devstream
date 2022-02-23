## 1 `gitlabci-golang` Plugin

This plugin creates Golang GitLab CI workflow.

## 2 Usage:

_This plugin depends on an environment variable "GITLAB_TOKEN", which is your GitLab personal access token. Set it before using this plugin._

See the [official document](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) for more information.

```yaml
tools:
- name: go-hello-world
  plugin:
    # name of the plugin
    kind: gitlabci-golang
    # version of the plugin
    # checkout the version from the GitHub releases
    version: 0.0.2
  # options for the plugin
  options:
    # owner/repo; "path with namespace" is only GitLab API's way of saying the same thing.
    pathWithNamespace: ironcore864/golang-demo
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: main
```

All parameters are mandatory.
