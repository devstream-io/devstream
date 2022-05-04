package template

// TplFile is a file contains some template tags like "{{ .Name }}".
// eg. internal/pkg/develop/plugin/template/create.go is a TplFile.
type TplFile struct {
	NameTpl       string
	DirTpl        string
	ContentTpl    string
	MustExistFlag bool
}

// File is a rendered TplFile that doesn't contain any template tags like "{{ .Name }}".
type File struct {
	Name          string
	Dir           string
	Content       string
	MustExistFlag bool
}

// TplFiles filled by functions at other go files.
// eg. internal/pkg/develop/plugin/template/create.go init()
var TplFiles = make([]TplFile, 0)
