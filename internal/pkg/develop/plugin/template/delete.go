package template

var delete_go_nameTpl = "delete.go"
var delete_go_dirTpl = "internal/pkg/plugin/{{ .Name }}/"
var delete_go_contentTpl = `package {{ .Name }}

func Delete(options map[string]interface{}) (bool, error) {
    // TODO(dtm): Add your logic here.

    return false, nil
}
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    delete_go_nameTpl,
		DirTpl:     delete_go_dirTpl,
		ContentTpl: delete_go_contentTpl,
	})
}
