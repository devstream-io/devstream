package java

import (
	_ "embed"

	"github.com/devstream-io/devstream/pkg/util/template"
)

//go:embed gitlabci-template.yml
var gitlabCITemplate string

// Render gitlab-ci.yml template with Options
func renderTmpl(Opts *Options) (string, error) {
	return template.New().FromContent(gitlabCITemplate).SetDefaultRender("gitlabci-java", Opts).Render()
}
