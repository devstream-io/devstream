package plugin

var DevlakeDefaultConfig = `tools:
# name of the tool
- name: devlake
  # id of the tool instance
  instanceID: default
  # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: []`
