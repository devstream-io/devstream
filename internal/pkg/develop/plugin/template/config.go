package template

var config_go_nameTpl = "{{ .Name }}.go"
var config_go_dirTpl = "internal/pkg/show/config/plugin/"
var config_go_contentTpl = `package plugin

// TODO(dtm): Add your default config here.
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    config_go_nameTpl,
		DirTpl:     config_go_dirTpl,
		ContentTpl: config_go_contentTpl,
	})
}
