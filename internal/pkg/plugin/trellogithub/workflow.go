package trellogithub

import (
	"github.com/devstream-io/devstream/internal/pkg/plugin/trellogithub/trello"
	"github.com/devstream-io/devstream/pkg/util/github"
)

const (
	defaultCommitMessage = "builder by DevStream"
	BuilderYmlTrello     = "trello-github-integ.yml"
)

var trelloWorkflow = &github.Workflow{
	CommitMessage:    defaultCommitMessage,
	WorkflowFileName: BuilderYmlTrello,
	WorkflowContent:  trello.IssuesBuilder,
}

// Options is the struct for configurations of the trellogithub plugin.
type Options struct {
	Owner       string
	Org         string
	Repo        string
	Branch      string
	BoardId     string
	TodoListId  string
	DoingListId string
	DoneListId  string
}
