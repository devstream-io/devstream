## 1 `gitlabci-golang` Plugin

This plugin creates Golang GitLab CI workflow.

## 2 Usage:

_This plugin depends on an environment variable "GITLAB_TOKEN", which is your GitLab personal access token. Set it before using this plugin._

See the [official document](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) for more information.

```yaml
tools:
- name: go-hello-world
  # name of the plugin
  plugin: gitlabci-golang
  # options for the plugin
  options:
    # owner/repo; "path with namespace" is only GitLab API's way of saying the same thing; please change the values below.
    pathWithNamespace: YOUR_GITLAB_USERNAME/YOUR_GITLAB_REPO_NAME
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: main
```

All parameters are mandatory.
