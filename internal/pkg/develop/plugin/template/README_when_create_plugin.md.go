package template

var ReadmeWhenCreatePluginMdNameTpl = "README_when_create_plugin.md"
var ReadmeWhenCreatePluginMdDirTpl = "./"
var ReadmeWhenCreatePluginMdContentTpl = `# Note with **dtm develop create-plugin**

## Done-Done Check

- [ ] I've finished all the TODO items and deleted all the *TODO(dtm)* comments.
`

func init() {
	TplFiles = append(TplFiles, TplFile{
		NameTpl:    ReadmeWhenCreatePluginMdNameTpl,
		DirTpl:     ReadmeWhenCreatePluginMdDirTpl,
		ContentTpl: ReadmeWhenCreatePluginMdContentTpl,
	})
}
