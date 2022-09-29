package template

import _ "embed"

var NamePluginMdNameTpl = "[[ .Name ]].md"
var NamePluginMdDirTpl = "docs/plugins/"

//go:embed NAME.md.tpl
var NamePluginMdContentTpl string

func init() {

	TplFiles = append(TplFiles, TplFile{
		NameTpl:    NamePluginMdNameTpl,
		DirTpl:     NamePluginMdDirTpl,
		ContentTpl: NamePluginMdContentTpl,
	})
}
