package template

var NAME_go_nameTpl = "{{ .Name }}.go"
var NAME_go_dirTpl = "internal/pkg/plugin/{{ .Name }}/"
var NAME_go_contentTpl = `package {{ .Name }}

// TODO(dtm): Add your logic here.
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    NAME_go_nameTpl,
		DirTpl:     NAME_go_dirTpl,
		ContentTpl: NAME_go_contentTpl,
	})
}
