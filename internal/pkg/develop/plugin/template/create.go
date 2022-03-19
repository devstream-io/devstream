package template

var create_go_nameTpl = "create.go"
var create_go_dirTpl = "internal/pkg/plugin/{{ .Name }}/"
var create_go_contentTpl = `package {{ .Name }}

func Create(options map[string]interface{}) (map[string]interface{}, error) {
    // TODO(dtm): Add your logic here.

    return nil, nil
}
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    create_go_nameTpl,
		DirTpl:     create_go_dirTpl,
		ContentTpl: create_go_contentTpl,
	})
}
