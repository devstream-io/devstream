package jiragithub

import (
	github "github.com/devstream-io/devstream/pkg/util/github"
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
