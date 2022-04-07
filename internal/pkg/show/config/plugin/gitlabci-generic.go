package plugin

var GitlabCIGenericDefaultConfig = `tools:
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
      key2: value2`
