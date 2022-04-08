package template

var NAME_plugin_md_nameTpl = "{{ .Name }}.md"
var NAME_plugin_md_dirTpl = "docs/plugins/"

// TODO(daniel-hutao): * -> `
var NAME_plugin_md_contentTpl = `# *{{ .Name }}* Plugin

// TODO(dtm): Add your document here.
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    NAME_plugin_md_nameTpl,
		DirTpl:     NAME_plugin_md_dirTpl,
		ContentTpl: NAME_plugin_md_contentTpl,
	})
}
