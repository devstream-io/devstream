package python

import (
	ga "github.com/merico-dev/stream/internal/pkg/plugin/githubactions"
	github "github.com/merico-dev/stream/pkg/util/github"
)

var workflows = []*github.Workflow{
	{CommitMessage: ga.CommitMessage, WorkflowFileName: ga.MainBuilderFileName, WorkflowContent: mainPipeline},
}
