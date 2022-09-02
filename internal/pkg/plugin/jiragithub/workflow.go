package jiragithub

import (
	"github.com/devstream-io/devstream/pkg/util/scm/github"
)

const (
	CommitMessage   = "jira-github-integ github actions workflow, created by DevStream"
	BuilderFileName = "jira-github-integ.yml"
)

var workflow = &github.Workflow{
	CommitMessage:    CommitMessage,
	WorkflowFileName: BuilderFileName,
	WorkflowContent:  jiraIssuesBuilder,
}
