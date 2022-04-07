# `gitlabci-generic` Plugin

This plugin creates Golang GitLab CI workflow.

It downloads a template of your choice, render it with provided parameters, and creates a GitLab CI file to your repo.

## Usage:

_This plugin depends on an environment variable "GITLAB_TOKEN", which is your GitLab personal access token. Set it before using this plugin._

See the [official document](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) for more information.

```yaml
tools:
- name: myapp-ci
  # name of the plugin
  plugin: gitlabci-generic
  options:
    # owner/repo; "path with namespace" is only GitLab API's way of saying the same thing; please change the values below.
    pathWithNamespace: YOUR_GITLAB_USERNAME/YOUR_GITLAB_REPO_NAME
    # main branch of the repo (to which branch the plugin will submit the workflows)
    branch: main
    # url of the GitLab CI template
    templateURL: https://someplace.com/to/download/your/template
    # custom variables keys and values
    templateVariables:
      key1: value1
      key2: value2
```

Or, run `dtm show config --plugin=gitlabci-generic` to get the default config.

All parameters are mandatory.
