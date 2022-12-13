package jira

import (
	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/internal/pkg/configmanager"
	"github.com/devstream-io/devstream/pkg/util/scm"
)

// addScmSecret will add jira secret in scm
func addScmSecret(options configmanager.RawOptions) error {
	opts, err := newOptions(options)
	if err != nil {
		return err
	}
	scmClient, err := scm.NewClientWithAuth(opts.Scm)
	if err != nil {
		return err
	}

	// JIRA_API_TOKEN, how to get it: "https://help.siteimprove.com/support/solutions/articles/80000448174-how-to-create-an-api-token-from-your-atlassian-account"
	if err := scmClient.AddRepoSecret("JIRA_API_TOKEN", viper.GetString("jira_api_token")); err != nil {
		return err
	}

	if err := scmClient.AddRepoSecret("GH_TOKEN", viper.GetString("github_token")); err != nil {
		return err
	}
	return nil
}
