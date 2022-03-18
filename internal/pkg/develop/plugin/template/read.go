package template

var read_go_nameTpl = "read.go"
var read_go_dirTpl = "internal/pkg/plugin/{{ .Name }}/"
var read_go_contentTpl = `package {{ .Name }}

func Read(options map[string]interface{}) (map[string]interface{}, error) {
    // TODO(dtm): Add your logic here.

    return nil, nil
}
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    read_go_nameTpl,
		DirTpl:     read_go_dirTpl,
		ContentTpl: read_go_contentTpl,
	})
}
