package jenkinspipeline

import (
	_ "embed"
)

//go:embed tpl/gitlab-casc.tpl.yaml
var gitlabConnectionCascConfig string

var jenkinsPlugins = []string{
	"gitlab-plugin",
	"kubernetes",
	"git",
	"configuration-as-code",
}
