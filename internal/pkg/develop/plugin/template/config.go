package template

var config_go_nameTpl = "{{ .Name }}.yaml"
var config_go_dirTpl = "internal/pkg/show/config/plugins/"
var config_go_contentTpl = `tools:
# name of the tool
- name: {{ .Name }}
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
		NameTpl:    config_go_nameTpl,
		DirTpl:     config_go_dirTpl,
		ContentTpl: config_go_contentTpl,
	})
}
