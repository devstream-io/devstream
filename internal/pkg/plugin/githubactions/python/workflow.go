package python

import (
	"github.com/devstream-io/devstream/internal/pkg/plugininstaller/github"
	githubCommon "github.com/devstream-io/devstream/pkg/util/github"
)

var workflows = []*githubCommon.Workflow{
	{CommitMessage: github.CommitMessage, WorkflowFileName: github.MainBuilderFileName, WorkflowContent: mainPipeline},
}
