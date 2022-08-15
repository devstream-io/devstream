package template

var configGoNameTpl = "[[ .Name ]].yaml"
var configGoDirTpl = "internal/pkg/show/config/plugins/"
var configGoContentTpl = `tools:
# name of the tool
- name: [[ .Name ]]
  # id of the tool instance
  instanceID: default
  # format: name.instanceID; If specified, dtm will make sure the dependency is applied first before handling this tool.
  dependsOn: [ ]
  # options for the plugin
  options:
  # TODO(dtm): Add your default config here.
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    configGoNameTpl,
		DirTpl:     configGoDirTpl,
		ContentTpl: configGoContentTpl,
	})
}
