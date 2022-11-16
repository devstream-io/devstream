package golang

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/plugininstaller/github"
	githubCommon "github.com/devstream-io/devstream/pkg/util/scm/github"
)

var workflows = []*githubCommon.Workflow{
	{CommitMessage: github.CommitMessage, WorkflowFileName: github.PRBuilderFileName, WorkflowContent: prPipeline},
	{CommitMessage: github.CommitMessage, WorkflowFileName: github.MainBuilderFileName, WorkflowContent: mainPipeline},
}
