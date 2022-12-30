package jira

import (
	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

// addJiraSecret will add trello secret in github
func addJiraSecret(rawOptions configmanager.RawOptions) error {
	opts, err := newOptions(rawOptions)
	if err != nil {
		return err
	}

	scmClient, err := scm.NewClientWithAuth(opts.Scm)
	if err != nil {
		return err
	}
	if err := scmClient.AddRepoSecret("GH_TOKEN", opts.Scm.Token); err != nil {
		return err
	}
	if err := scmClient.AddRepoSecret("JIRA_API_TOKEN", opts.Jira.Token); err != nil {
		return err
	}
	return nil
}
