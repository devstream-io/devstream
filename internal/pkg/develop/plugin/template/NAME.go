package template

var NAME_go_nameTpl = "{{ .Name | format  }}.go"
var NAME_go_dirTpl = "internal/pkg/plugin/{{ .Name | dirFormat }}/"
var NAME_go_mustExistFlag = true
var NAME_go_contentTpl = `package {{ .Name | format }}

// TODO(dtm): Add your logic here.
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:       NAME_go_nameTpl,
		DirTpl:        NAME_go_dirTpl,
		ContentTpl:    NAME_go_contentTpl,
		MustExistFlag: NAME_go_mustExistFlag,
	})
}
