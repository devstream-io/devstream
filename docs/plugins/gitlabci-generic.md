# gitlabci-generic Plugin

This plugin creates Golang GitLab CI workflow.

It downloads a template of your choice, render it with provided parameters, and creates a GitLab CI file to your repo.

## Usage

_This plugin depends on an environment variable "GITLAB_TOKEN", which is your GitLab personal access token._

TL;DR: if you are using gitlab.com (instead of a self-hosted GitLab), [click here](https://gitlab.com/-/profile/personal_access_tokens?name=DevStream+Access+token&scopes=api) to create a token for DevStream (the scope contains API only.)

If you are using self-hosted GitLab, refer to the [official doc here](https://docs.gitlab.com/ee/user/profile/personal_access_tokens.html) for more info.

_Note: when creating the token, make sure you select "API" in the "scopes" section, as DevStream uses GitLab API to add CI workflow files._

Plugin config example:

```yaml
tools:
# name of the tool
- name: gitlabci-generic
  # id of the tool instance
  instanceID: default
  # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: []
  # options for the plugin
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
