package jiragithub

import (
	"github.com/spf13/viper"

	"github.com/devstream-io/devstream/pkg/util/github"
)

func setRepoSecrets(gitHubClient *github.Client) error {

	// JIRA_API_TOKEN, how to get it: "https://help.siteimprove.com/support/solutions/articles/80000448174-how-to-create-an-api-token-from-your-atlassian-account"
	if err := gitHubClient.AddRepoSecret("JIRA_API_TOKEN", viper.GetString("jira_api_token")); err != nil {
		return err
	}

	if err := gitHubClient.AddRepoSecret("GH_TOKEN", viper.GetString("github_token")); err != nil {
		return err
	}

	return nil
}
