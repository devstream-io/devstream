package jenkinspipeline

import (
	_ "embed"

	"github.com/devstream-io/devstream/pkg/util/jenkins"
)

//go:embed tpl/gitlab-casc.tpl.yaml
var gitlabConnectionCascConfig string

var jenkinsPlugins = []*jenkins.JenkinsPlugin{
	// gitlab-plugin is used for jenkins gitlab integration
	{
		Name:    "gitlab-plugin",
		Version: "1.5.35",
	},
}
