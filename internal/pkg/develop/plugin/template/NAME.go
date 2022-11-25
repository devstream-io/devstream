package template

var NameGoNameTpl = "[[ .Name | format  ]].go"
var NameGoDirTpl = "internal/pkg/plugin/[[ .Name | dirFormat ]]/"
var NameGoMustExistFlag = true
var NameGoContentTpl = `package [[ .Name | format ]]

const Name = "[[ .Name ]]"

// TODO(dtm): Add your logic here.
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:       NameGoNameTpl,
		DirTpl:        NameGoDirTpl,
		ContentTpl:    NameGoContentTpl,
		MustExistFlag: NameGoMustExistFlag,
	})
}
