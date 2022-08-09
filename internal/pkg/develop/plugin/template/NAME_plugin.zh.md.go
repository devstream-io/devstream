package template

import "fmt"

var NAME_plugin_zh_md_nameTpl = "{{ .Name }}.zh.md"
var NAME_plugin_zh_md_dirTpl = "docs/plugins/"

var NAME_plugin_zh_md_contentTpl = `# {{ .Name }} 插件

TODO(dtm): 在这里添加文档.

## 用例

%s yaml
--8<-- "{{ .Name }}.yaml"
%s
`

func init() {

	TplFiles = append(TplFiles, TplFile{
		NameTpl:    NAME_plugin_zh_md_nameTpl,
		DirTpl:     NAME_plugin_zh_md_dirTpl,
		ContentTpl: fmt.Sprintf(NAME_plugin_zh_md_contentTpl, "```", "```"),
	})
}
