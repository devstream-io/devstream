package trellogithub

import (
	"github.com/merico-dev/stream/internal/pkg/plugin/trellogithub/trello"
	"github.com/merico-dev/stream/pkg/util/github"
)

const (
	defaultCommitMessage = "builder by DevStream"
	BuilderYmlTrello     = "trello-github-integ.yml"
)

var extTrello = &Api{
	Name: "trello",
}

var defaultWorkflows = workflows{
	extTrello.Name: {
		{
			CommitMessage:    defaultCommitMessage,
			WorkflowFileName: BuilderYmlTrello,
			WorkflowContent:  trello.IssuesBuilder,
		},
	},
}

// Options is the struct for configurations of the trellogithub plugin.
type Options struct {
	Owner  string
	Repo   string
	Branch string
	Api    *Api
	Jobs   *Jobs
}

type Api struct {
	Name        string
	BoardId     string
	todoListId  string
	doingListId string
	doneListId  string
}

type workflows map[string][]*github.Workflow

func (ws *workflows) GetWorkflowByNameVersionTypeString(nvtStr string) []*github.Workflow {
	workflowList, exist := (*ws)[nvtStr]
	if exist {
		return workflowList
	}
	return nil
}
