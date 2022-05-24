package template

var options_go_nameTpl = "options.go"
var options_go_dirTpl = "internal/pkg/plugin/{{ .Name | dirFormat }}/"
var options_go_contentTpl = `package {{ .Name | format }}

// Options is the struct for configurations of the {{ .Name }} plugin.
type Options struct {
    // TODO(dtm): Add your params here.
	Foo string
}
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    options_go_nameTpl,
		DirTpl:     options_go_dirTpl,
		ContentTpl: options_go_contentTpl,
	})
}
