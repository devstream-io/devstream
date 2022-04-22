package plugin

var GitlabCIGenericDefaultConfig = `tools:
# name of the tool
- name: gitlabci-generic
  # id of the tool instance
  instanceID: myapp-ci
  # optional; if specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ "TOOL1_NAME.TOOL1_PLUGIN", "TOOL2_NAME.TOOL2_PLUGIN" ]
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
      key2: value2`
