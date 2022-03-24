package template

var update_go_nameTpl = "update.go"
var update_go_dirTpl = "internal/pkg/plugin/{{ .Name }}/"
var update_go_contentTpl = `package {{ .Name }}

func Update(options map[string]interface{}) (map[string]interface{}, error) {
    // TODO(dtm): Add your logic here.

    return nil, nil
}
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    update_go_nameTpl,
		DirTpl:     update_go_dirTpl,
		ContentTpl: update_go_contentTpl,
	})
}
