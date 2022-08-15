package template

var validateGoNameTpl = "validate.go"
var validateGoDirTpl = "internal/pkg/plugin/[[ .Name | dirFormat ]]/"
var validateGoContentTpl = `package [[ .Name | format ]]

// validate validates the options provided by the core.
func validate(options *Options) []error {
    // TODO(dtm): Add your logic here.
	return make([]error, 0)
}
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    validateGoNameTpl,
		DirTpl:     validateGoDirTpl,
		ContentTpl: validateGoContentTpl,
	})
}
