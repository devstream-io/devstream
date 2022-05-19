package template

var validate_go_nameTpl = "validate.go"
var validate_go_dirTpl = "internal/pkg/plugin/{{ .Name | dirFormat }}/"
var validate_go_contentTpl = `package {{ .Name | format }}

// validate validates the options provided by the core.
func validate(options *Options) []error {
    // TODO(dtm): Add your logic here.
	return make([]error, 0)
}
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    validate_go_nameTpl,
		DirTpl:     validate_go_dirTpl,
		ContentTpl: validate_go_contentTpl,
	})
}
