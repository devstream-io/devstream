package python

import (
	ga "github.com/devstream-io/devstream/internal/pkg/plugin/githubactions"
	github "github.com/devstream-io/devstream/pkg/util/github"
)

var workflows = []*github.Workflow{
	{CommitMessage: ga.CommitMessage, WorkflowFileName: ga.MainBuilderFileName, WorkflowContent: mainPipeline},
}
