package template

var optionsGoNameTpl = "options.go"
var optionsGoDirTpl = "internal/pkg/plugin/[[ .Name | dirFormat ]]/"
var optionsGoContentTpl = `package [[ .Name | format ]]

// Options is the struct for configurations of the [[ .Name ]] plugin.
type Options struct {
    // TODO(dtm): Add your params here.
	Foo string
}
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    optionsGoNameTpl,
		DirTpl:     optionsGoDirTpl,
		ContentTpl: optionsGoContentTpl,
	})
}
