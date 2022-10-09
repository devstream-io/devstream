package template

import _ "embed"

var NamePluginZhMdNameTpl = "[[ .Name ]].zh.md"
var NamePluginZhMdDirTpl = "docs/plugins/"

//go:embed NAME.zh.md.tpl
var NamePluginZhMdContentTpl string

func init() {

	TplFiles = append(TplFiles, TplFile{
		NameTpl:    NamePluginZhMdNameTpl,
		DirTpl:     NamePluginZhMdDirTpl,
		ContentTpl: NamePluginZhMdContentTpl,
	})
}
