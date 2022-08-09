package template

import "fmt"

var NAME_plugin_md_nameTpl = "{{ .Name }}.md"
var NAME_plugin_md_dirTpl = "docs/plugins/"

var NAME_plugin_md_contentTpl = `# {{ .Name }} plugin

TODO(dtm): Add your document here.

## Usage

%s yaml
--8<-- "{{ .Name }}.yaml"
%s
`

func init() {

	TplFiles = append(TplFiles, TplFile{
		NameTpl:    NAME_plugin_md_nameTpl,
		DirTpl:     NAME_plugin_md_dirTpl,
		ContentTpl: fmt.Sprintf(NAME_plugin_md_contentTpl, "```", "```"),
	})
}
