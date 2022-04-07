package template

var README_when_create_plugin_md_nameTpl = "README_when_create_plugin.md"
var README_when_create_plugin_md_dirTpl = "./"
var README_when_create_plugin_md_contentTpl = `# Note with **dtm develop create-plugin**

## Done-Done Check

- [ ] I've finished all the TODO items and deleted all the *TODO(dtm)* comments.
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    README_when_create_plugin_md_nameTpl,
		DirTpl:     README_when_create_plugin_md_dirTpl,
		ContentTpl: README_when_create_plugin_md_contentTpl,
	})
}
