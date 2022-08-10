package template

import _ "embed"

var NAME_plugin_md_nameTpl = "{{ .Name }}.md"
var NAME_plugin_md_dirTpl = "docs/plugins/"

//go:embed NAME_plugin.md.tpl
var NAME_plugin_md_contentTpl string

func init() {

	TplFiles = append(TplFiles, TplFile{
		NameTpl:    NAME_plugin_md_nameTpl,
		DirTpl:     NAME_plugin_md_dirTpl,
		ContentTpl: NAME_plugin_md_contentTpl,
	})
}
