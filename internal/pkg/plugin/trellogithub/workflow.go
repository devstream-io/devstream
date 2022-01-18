package trellogithub

import (
	"github.com/merico-dev/stream/internal/pkg/plugin/trellogithub/trello"
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
		{defaultCommitMessage, BuilderYmlTrello, trello.IssuesBuilder},
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
	Name string
}

// Workflow is the struct for a GitHub Actions Workflow.
type Workflow struct {
	commitMessage    string
	workflowFileName string
	workflowContent  string
}

type ToolString string

type workflows map[string][]*Workflow

func (ws *workflows) GetWorkflowByNameVersionTypeString(nvtStr string) []*Workflow {
	workflowList, exist := (*ws)[nvtStr]
	if exist {
		return workflowList
	}
	return nil
}
