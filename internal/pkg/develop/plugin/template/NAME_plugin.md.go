package template

var NAME_plugin_md_nameTpl = "{{ .Name }}.md"
var NAME_plugin_md_dirTpl = "docs/plugins/"

// TODO(daniel-hutao): * -> `
var NAME_plugin_md_contentTpl = "# {{ .Name }} plugin\n\nTODO(dtm): Add your document here.\n## Usage\n\n" + "```" + "yaml\n\n--8<-- \"{{ .Name }}.yaml\"\n\n" + "```"

func init() {

	TplFiles = append(TplFiles, TplFile{
		NameTpl:    NAME_plugin_md_nameTpl,
		DirTpl:     NAME_plugin_md_dirTpl,
		ContentTpl: NAME_plugin_md_contentTpl,
	})
}
