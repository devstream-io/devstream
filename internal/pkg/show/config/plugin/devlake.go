package plugin

var DevlakeDefaultConfig = `tools:
- name: devlake
  # name of the plugin
  plugin: devlake
  # if specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ "TOOL1_NAME.TOOL1_PLUGIN", "TOOL2_NAME.TOOL2_PLUGIN" ]`
