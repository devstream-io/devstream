package python

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/github"
	githubCommon "github.com/devstream-io/devstream/pkg/util/scm/github"
)

var workflows = []*githubCommon.Workflow{
	{CommitMessage: github.CommitMessage, WorkflowFileName: "lint.yml", WorkflowContent: lintPipeline},
	{CommitMessage: github.CommitMessage, WorkflowFileName: "test.yml", WorkflowContent: testPipeline},
	{CommitMessage: github.CommitMessage, WorkflowFileName: "docker.yml", WorkflowContent: dockerPipeline},
}
