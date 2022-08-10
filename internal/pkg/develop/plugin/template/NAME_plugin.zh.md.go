package template

import _ "embed"

var NAME_plugin_zh_md_nameTpl = "{{ .Name }}.zh.md"
var NAME_plugin_zh_md_dirTpl = "docs/plugins/"

//go:embed NAME_plugin.zh.md.tpl
var NAME_plugin_zh_md_contentTpl string

func init() {

	TplFiles = append(TplFiles, TplFile{
		NameTpl:    NAME_plugin_zh_md_nameTpl,
		DirTpl:     NAME_plugin_zh_md_dirTpl,
		ContentTpl: NAME_plugin_zh_md_contentTpl,
	})
}
