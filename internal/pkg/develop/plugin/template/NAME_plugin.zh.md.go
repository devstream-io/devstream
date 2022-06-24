package template

var NAME_plugin_zh_md_nameTpl = "{{ .Name }}.zh.md"
var NAME_plugin_zh_md_dirTpl = "docs/plugins/"

// TODO(daniel-hutao): * -> `
var NAME_plugin_zh_md_contentTpl = "# {{ .Name }} 插件\n\nTODO(dtm): 在这里添加文档.\n## 用例\n\n" + "```" + "yaml\n\n--8<-- \"{{ .Name }}.yaml\"\n\n" + "```"

func init() {

	TplFiles = append(TplFiles, TplFile{
		NameTpl:    NAME_plugin_zh_md_nameTpl,
		DirTpl:     NAME_plugin_zh_md_dirTpl,
		ContentTpl: NAME_plugin_zh_md_contentTpl,
	})
}
